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

func TestRouting_GetJson(t *testing.T) {
	expectedPerson, _ := createRandomPerson()
	expectedPersonJson, _ := json.Marshal(expectedPerson)
	expectedPersonStr := string(expectedPersonJson)
	req := httptest.NewRequest("GET", "http://blah.com/foo", bytes.NewBufferString(expectedPersonStr))

	var person Person
	if err := getJson(req.Body, person); err != nil {
		t.Fatalf("Encountered error when retrieving json from string: %s", err)
	}
}

func TestRouting_PostJson(t *testing.T) {
	expectedPerson, _ := createRandomPerson()
	expectedStatus := http.StatusOK
	response := httptest.NewRecorder()
	postJson(response, expectedStatus, expectedPerson)

	if response.Code != expectedStatus {
		t.Fatalf("Expected http status %d but found %d", expectedStatus, response.Code)
	}

	var person Person
	if err := json.Unmarshal(response.Body.Bytes(), &person); err != nil {
		t.Fatalf("Encountered error when retrieving json from string: %s", err)
	} else if !arePersonEqual(expectedPerson, &person) {
		t.Fatal("Did not receive expected person")
	}
	// TODO: more?
}

func TestRouting_UpdatePersonHandler(t *testing.T) {
	log := log.New(os.Stdout, "", 0)
	response := httptest.NewRecorder()

	expectedPerson, _ := createRandomPerson()
	expectedPersonJson, _ := json.Marshal(expectedPerson)
	expectedPersonStr := string(expectedPersonJson)
	req := httptest.NewRequest("GET", "http://"+createRandomString(), bytes.NewBufferString(expectedPersonStr))

	db, server, _ := createRandomDbClient()
	defer server.Close()

	UpdatePersonHandler(log, db, response, req)

	if response.Code != http.StatusCreated {
		t.Errorf("Did not get expected HTTP status code. Instead got: %d", response.Code)
	}

	if person, err := (*db).Get(expectedPerson.Id); err != nil {
		t.Fatalf("Encountered error when retrieving person: %s", err)
	} else if !arePersonEqual(expectedPerson, person) {
		t.Fatal("Retrieved Person is not equivalent to the expected Person")
	}

	// TODO: more?
}

func TestRouting_PersonRequestHandler(t *testing.T) {
	log := log.New(os.Stdout, "", 0)
	response := httptest.NewRecorder()

	expectedPerson, _ := createRandomPerson()

	var dataReq PersonDataRequest
	dataReq.Ids = make([]string, 1)
	dataReq.Ids[0] = expectedPerson.Id
	personDataRequestJson, _ := json.Marshal(dataReq)
	personDataRequestStr := string(personDataRequestJson)
	req := httptest.NewRequest("GET", "http://"+createRandomString(), bytes.NewBufferString(personDataRequestStr))

	db, server, _ := createRandomDbClient()
	defer server.Close()

	(*db).Create(*expectedPerson)

	PersonRequestHandler(log, db, response, req)

	// TODO: check status code
	// TODO: check number of people returned

	var person PersonDataResponse
	if err := json.Unmarshal(response.Body.Bytes(), &person); err != nil {
		t.Fatalf("Encountered error when retrieving json from string: %s", err)
	} else if !arePersonEqual(expectedPerson, &person.People[0]) {
		t.Fatal("Did not receive expected person")
	}

	// Test that ... the ... dummy db works...
	/* TODO: this mostly just tests the mock db anyway
	req = httptest.NewRequest("GET", "http://"+createRandomString(), bytes.NewBufferString(createRandomString()))
	PersonRequestHandler(log, db, response, req)
	if response.Code != http.StatusUnprocessableEntity {
		t.Errorf("Did not receive expected failure HTTP status on bad request body: %d", response.Code)
	}
	*/
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
