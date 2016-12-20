package dondeestas

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func getExpectedPersonRequest() (*http.Request, *Person, *httptest.ResponseRecorder) {
	expectedPerson, _ := createRandomPerson()
	expectedPersonJSON, _ := json.Marshal(expectedPerson)
	expectedPersonStr := string(expectedPersonJSON)
	req := httptest.NewRequest("GET", "http://"+createRandomString(), bytes.NewBufferString(expectedPersonStr))

	return req, expectedPerson, httptest.NewRecorder()
}

func TestRouting_PostJson(t *testing.T) {
	expectedPerson, _ := createRandomPerson()
	expectedStatus := http.StatusOK
	response := httptest.NewRecorder()
	if err := postJSON(response, expectedStatus, expectedPerson); err != nil {
		t.Fatalf("Unexpected error posting a JSON string: %s", err)
	}

	if response.Code != expectedStatus {
		t.Fatalf("Expected http status %d but found %d", expectedStatus, response.Code)
	}

	var person Person
	if err := json.Unmarshal(response.Body.Bytes(), &person); err != nil {
		t.Fatalf("Encountered error when retrieving json from string: %s", err)
	} else if !arePersonEqual(expectedPerson, &person) {
		t.Fatal("Did not receive expected person")
	}

	// Attempt to post a struct with a struct that cannot be converted to JSON
	badJSON := struct{ IntChan chan int }{IntChan: make(chan int)}
	if err := postJSON(response, expectedStatus, badJSON); err == nil {
		t.Error("Unexpectedly did not encounter an error when posting bad JSON struct")
	}
}

func getPersonBufferString(p *Person) *bytes.Buffer {
	json, _ := json.Marshal(p)
	return bytes.NewBufferString(string(json))
}

func TestRouting_updatePersonHandler(t *testing.T) {
	log := log.New(os.Stdout, "", 0)
	db, server, _ := createRandomDbClient()

	// Build a person
	expectedPerson, _ := createRandomPerson()

	req := httptest.NewRequest("GET", "http://"+createRandomString(), getPersonBufferString(expectedPerson))
	response := httptest.NewRecorder()
	updatePersonHandler(log, db, response, req)
	if response.Code != http.StatusCreated {
		t.Errorf("Did not get expected HTTP status code. Instead got: %d", response.Code)
	}

	if person, err := (*db).Get(expectedPerson.ID); err != nil {
		t.Fatalf("Encountered error when retrieving person: %s", err)
	} else if !arePersonEqual(expectedPerson, person) {
		t.Fatal("Retrieved Person is not equivalent to the expected Person")
	}

	// Test updating the same person
	expectedPerson.Name = createRandomString()
	req = httptest.NewRequest("GET", "http://"+createRandomString(), getPersonBufferString(expectedPerson))
	response = httptest.NewRecorder()
	updatePersonHandler(log, db, response, req)
	if response.Code != http.StatusCreated {
		t.Errorf("Did not get expected HTTP status code. Instead got: %d", response.Code)
	}

	if person, err := (*db).Get(expectedPerson.ID); err != nil {
		t.Fatalf("Encountered error when retrieving person: %s", err)
	} else if !arePersonEqual(expectedPerson, person) {
		t.Fatal("Retrieved Person is not equivalent to the expected Person")
	}

	// Test bad request
	req = httptest.NewRequest("GET", "http://"+createRandomString(), bytes.NewBufferString(createRandomString()))
	response = httptest.NewRecorder()
	updatePersonHandler(log, db, response, req)
	if response.Code != http.StatusUnprocessableEntity {
		t.Errorf("Did not get expected HTTP error status code of %d. Instead got: %d", http.StatusUnprocessableEntity, response.Code)
	}

	// Test unable to process the body
	expectedPerson.ID = ""
	req = httptest.NewRequest("GET", "http://"+createRandomString(), getPersonBufferString(expectedPerson))
	response = httptest.NewRecorder()
	updatePersonHandler(log, db, response, req)
	if response.Code != http.StatusInternalServerError {
		t.Errorf("Did not get expected HTTP error status code of %d. Instead got: %d", http.StatusInternalServerError, response.Code)
	}

	// Test when the database is unable to comply with the request
	req = httptest.NewRequest("GET", "http://"+createRandomString(), getPersonBufferString(expectedPerson))
	response = httptest.NewRecorder()
	server.Close()
	updatePersonHandler(log, db, response, req)
	if response.Code != http.StatusInternalServerError {
		t.Errorf("Did not get expected HTTP error status code of %d. Instead got: %d", http.StatusInternalServerError, response.Code)
	}
}

func createPersonDataRequest(ids []string) (*httptest.ResponseRecorder, *http.Request) {
	var dataReq personDataRequest
	dataReq.Ids = make([]string, len(ids))

	for i, v := range ids {
		dataReq.Ids[i] = v
	}
	personDataRequestJSON, _ := json.Marshal(dataReq)
	personDataRequestStr := string(personDataRequestJSON)
	req := httptest.NewRequest("GET", "http://"+createRandomString(), bytes.NewBufferString(personDataRequestStr))

	return httptest.NewRecorder(), req
}

func TestRouting_personRequestHandler(t *testing.T) {
	log := log.New(os.Stdout, "", 0)

	db, server, _ := createRandomDbClient()
	defer server.Close()

	expectedPerson, _ := createRandomPerson()

	(*db).Create(*expectedPerson)

	// Test usual case
	response, req := createPersonDataRequest([]string{expectedPerson.ID})
	personRequestHandler(log, db, response, req)
	if response.Code != http.StatusOK {
		t.Fatalf("Encountered unexpected HTTP code: %d\n", response.Code)
	}

	var person personDataResponse
	if err := json.Unmarshal(response.Body.Bytes(), &person); err != nil {
		t.Fatalf("Encountered error when retrieving json from string: %s", err)
	} else if len(person.People) > 1 {
		t.Fatalf("Expected 1 Person object returned but found %d\n", len(person.People))
	} else if !arePersonEqual(expectedPerson, &person.People[0]) {
		t.Fatal("Did not receive expected person")
	}

	// Test getting partial list
	response, req = createPersonDataRequest([]string{expectedPerson.ID, createRandomString()})
	personRequestHandler(log, db, response, req)
	if response.Code != http.StatusPartialContent {
		t.Fatalf("Expected HTTP Partial Content code (%d) but found %d instead\n", http.StatusPartialContent, response.Code)
	}
	// TODO: verify the number of people returned and that the correct number was returned

	// Test bad request
	req = httptest.NewRequest("GET", "http://"+createRandomString(), bytes.NewBufferString(createRandomString()))
	response = httptest.NewRecorder()
	personRequestHandler(log, db, response, req)
	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("Expected HTTP error code %d but got %d instead\n", http.StatusUnprocessableEntity, response.Code)
	}
}

func ExampleListenAndServe(t *testing.T) {
	db, err := NewDbClient(DbClientParams{CouchDB, "example_db", "localhost", 5934})
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", 0)

	logger.Fatal(ListenAndServe(logger, 8080, db))
}
