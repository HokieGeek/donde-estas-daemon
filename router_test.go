package dondeestas

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter(t *testing.T) {
	t.Skip("TODO?")
	// req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	// w := httptest.NewRecorder()

	// dummyServer := getServer(expectedData)
	// defer dummyServer.Close()
	// dummyServer.URL
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
	if err := getJson(ioutil.NopCloser(response.Body), person); err != nil {
		t.Fatalf("Encountered error when retrieving json from string: %s", err)
		// TODO: } else if !arePersonEqual(expectedPerson, &person) {
		// t.Fatal("Did not receive expected person")
	}
	// TODO: test more
}

/*
// func TestRouting_New(log *log.Logger, port int, db *dbclient) {

func TestRouting_PersonRequestHandler(t *testing.T) {
	PersonRequestHandler(log *log.Logger, db *dbclient, w http.ResponseWriter, r *http.Request)
}

func TestRouting_UpdatePersonHandler(t *testing.T) {
	UpdatePersonHandler(log *log.Logger, db *dbclient, w http.ResponseWriter, r *http.Request)
}
*/
