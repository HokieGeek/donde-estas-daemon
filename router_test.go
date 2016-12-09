package dondeestas

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func getExpectedPersonRequest() (*http.Request, *Person, *httptest.ResponseRecorder) {
	expectedPerson, _ := createRandomPerson()
	expectedPersonJson, _ := json.Marshal(expectedPerson)
	expectedPersonStr := string(expectedPersonJson)
	req := httptest.NewRequest("GET", "http://"+createRandomString(), bytes.NewBufferString(expectedPersonStr))

	return req, expectedPerson, httptest.NewRecorder()
}

func TestRouting_PostJson(t *testing.T) {
	expectedPerson, _ := createRandomPerson()
	expectedStatus := http.StatusOK
	response := httptest.NewRecorder()
	if err := postJson(response, expectedStatus, expectedPerson); err != nil {
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
}

func getPersonBufferString(p *Person) *bytes.Buffer {
	json, _ := json.Marshal(p)
	return bytes.NewBufferString(string(json))
}

func TestRouting_UpdatePersonHandler(t *testing.T) {
	log := log.New(os.Stdout, "", 0)
	db, server, _ := createRandomDbClient()

	// Build a person
	expectedPerson, _ := createRandomPerson()

	req := httptest.NewRequest("GET", "http://"+createRandomString(), getPersonBufferString(expectedPerson))
	response := httptest.NewRecorder()
	UpdatePersonHandler(log, db, response, req)
	if response.Code != http.StatusCreated {
		t.Errorf("Did not get expected HTTP status code. Instead got: %d", response.Code)
	}

	if person, err := (*db).Get(expectedPerson.Id); err != nil {
		t.Fatalf("Encountered error when retrieving person: %s", err)
	} else if !arePersonEqual(expectedPerson, person) {
		t.Fatal("Retrieved Person is not equivalent to the expected Person")
	}

	// Test updating the same person
	expectedPerson.Name = createRandomString()
	req = httptest.NewRequest("GET", "http://"+createRandomString(), getPersonBufferString(expectedPerson))
	response = httptest.NewRecorder()
	UpdatePersonHandler(log, db, response, req)
	if response.Code != http.StatusCreated {
		t.Errorf("Did not get expected HTTP status code. Instead got: %d", response.Code)
	}

	if person, err := (*db).Get(expectedPerson.Id); err != nil {
		t.Fatalf("Encountered error when retrieving person: %s", err)
	} else if !arePersonEqual(expectedPerson, person) {
		t.Fatal("Retrieved Person is not equivalent to the expected Person")
	}

	// Test bad request
	req = httptest.NewRequest("GET", "http://"+createRandomString(), bytes.NewBufferString(createRandomString()))
	response = httptest.NewRecorder()
	UpdatePersonHandler(log, db, response, req)
	if response.Code != http.StatusUnprocessableEntity {
		t.Errorf("Did not get expected HTTP error status code of %d. Instead got: %d", http.StatusUnprocessableEntity, response.Code)
	}

	// Test unable to process the body
	expectedPerson.Id = ""
	req = httptest.NewRequest("GET", "http://"+createRandomString(), getPersonBufferString(expectedPerson))
	response = httptest.NewRecorder()
	UpdatePersonHandler(log, db, response, req)
	if response.Code != http.StatusInternalServerError {
		t.Errorf("Did not get expected HTTP error status code of %d. Instead got: %d", http.StatusInternalServerError, response.Code)
	}

	// Test when the database is unable to comply with the request
	req = httptest.NewRequest("GET", "http://"+createRandomString(), getPersonBufferString(expectedPerson))
	response = httptest.NewRecorder()
	server.Close()
	UpdatePersonHandler(log, db, response, req)
	if response.Code != http.StatusInternalServerError {
		t.Errorf("Did not get expected HTTP error status code of %d. Instead got: %d", http.StatusInternalServerError, response.Code)
	}
}

func createPersonDataRequest(ids []string) (*httptest.ResponseRecorder, *http.Request) {
	var dataReq PersonDataRequest
	dataReq.Ids = make([]string, len(ids))

	for i, v := range ids {
		dataReq.Ids[i] = v
	}
	personDataRequestJson, _ := json.Marshal(dataReq)
	personDataRequestStr := string(personDataRequestJson)
	req := httptest.NewRequest("GET", "http://"+createRandomString(), bytes.NewBufferString(personDataRequestStr))

	return httptest.NewRecorder(), req
}

func TestRouting_PersonRequestHandler(t *testing.T) {
	log := log.New(os.Stdout, "", 0)

	db, server, _ := createRandomDbClient()
	defer server.Close()

	expectedPerson, _ := createRandomPerson()

	(*db).Create(*expectedPerson)

	// Test usual case
	response, req := createPersonDataRequest([]string{expectedPerson.Id})
	PersonRequestHandler(log, db, response, req)
	if response.Code != http.StatusOK {
		t.Fatalf("Encountered unexpected HTTP code: %d\n", response.Code)
	}

	var person PersonDataResponse
	if err := json.Unmarshal(response.Body.Bytes(), &person); err != nil {
		t.Fatalf("Encountered error when retrieving json from string: %s", err)
	} else if len(person.People) > 1 {
		t.Fatalf("Expected 1 Person object returned but found %d\n", len(person.People))
	} else if !arePersonEqual(expectedPerson, &person.People[0]) {
		t.Fatal("Did not receive expected person")
	}

	// Test getting partial list
	response, req = createPersonDataRequest([]string{expectedPerson.Id, createRandomString()})
	PersonRequestHandler(log, db, response, req)
	if response.Code != http.StatusPartialContent {
		t.Fatalf("Expected HTTP Partial Content code (%d) but found %d instead\n", http.StatusPartialContent, response.Code)
	}
	// TODO: verify the number of people returned and that the correct number was returned

	// Test bad request
	req = httptest.NewRequest("GET", "http://"+createRandomString(), bytes.NewBufferString(createRandomString()))
	response = httptest.NewRecorder()
	PersonRequestHandler(log, db, response, req)
	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("Expected HTTP error code %d but got %d instead\n", http.StatusUnprocessableEntity, response.Code)
	}
}

func TestRouting_ListenAndServe(t *testing.T) {
	t.Skip("TODO")
	port := rand.Int()
	log := log.New(os.Stdout, "", 0)
	db, server, _ := createRandomDbClient()
	defer server.Close()

	// TODO: use a goroutine to do some tests
	go ListenAndServe(log, port, db)
}
