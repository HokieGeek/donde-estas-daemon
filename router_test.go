package dondeestas

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	// "math/rand"
	// "log"
)

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

	fmt.Println(response.Body)
	var person Person
	if err := getJson(ioutil.NopCloser(response.Body), person); err != nil {
		t.Fatalf("Encountered error when retrieving json from string: %s", err)
		// TODO: } else if !arePersonEqual(expectedPerson, &person) {
		// t.Fatal("Did not receive expected person")
	}
}
/*
func TestRouting_UpdatePersonHandler(t *testing.T) {
	log := log.New(os.Stdout, "", 0)
	response := httptest.NewRecorder()

	expectedPerson, _ := createRandomPerson()
    req := httptest.NewRequest("GET", createRandomString(), bytes.NewBufferString(expectedPersonStr))

	db, server, _ := createRandomDbClient()
	defer server.Close()

	UpdatePersonHandler(log, db, response, req)
}

func TestRouting_PersonRequestHandler(t *testing.T) {
	log := log.New(os.Stdout, "", 0)
	response := httptest.NewRecorder()

    req := httptest.NewRequest("GET", createRandomString(), bytes.NewBufferString(expectedPersonStr))

	db, server, _ := createRandomDbClient()
	defer server.Close()

	PersonRequestHandler(log, db, response, req)
}

func TestRouting_New(t *testing.T) { // TODO TEST FOR BAD INPUT VALUES
	port := rand.Int()
	log := log.New(os.Stdout, "", 0)
	db, server, _ := createRandomDbClient()
	defer server.Close()

	New(log *log.Logger, port int, db *dbclient)
}
*/
