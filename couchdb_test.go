package dondeestas

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

type MockCouchDb struct {
	Name      string
	People    map[string]string
	Revisions map[string]string
}

func splitURL(url string) (string, uint16) {
	sepPos := strings.LastIndex(url, ":")
	port, err := strconv.ParseUint(url[sepPos+1:], 10, 16)
	if err != nil {
		return "", 0
	}
	return url[:sepPos], uint16(port)
}

func getMockCouchDbServer(db *MockCouchDb) *httptest.Server {
	db.People = make(map[string]string)
	db.Revisions = make(map[string]string)

	db.People["BADPERSON"] = createRandomString()
	db.Revisions["BADPERSON"] = createRandomString()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path[1:], "/")
		if len(path) > 0 {
			switch r.Method {
			case "HEAD":
				if len(path) == 1 {
					if path[0] == db.Name {
						w.WriteHeader(http.StatusOK)
					}
				} else if _, ok := db.People[path[1]]; ok {
					w.Header().Set("Etag", db.Revisions[path[1]])
					w.WriteHeader(http.StatusOK)
				}
			case "GET":
				if len(path) > 1 {
					if _, ok := db.People[path[1]]; ok {
						fmt.Fprint(w, db.People[path[1]])
						w.WriteHeader(http.StatusOK)
					}
				}
			case "PUT":
				if len(path) == 1 {
					db.Name = path[0]
					w.WriteHeader(http.StatusCreated)
				} else if path[1] != "" {
					if _, ok := db.People[path[1]]; ok && r.Header.Get("If-Match") != db.Revisions[path[1]] {
						w.WriteHeader(http.StatusConflict)
					} else {
						defer r.Body.Close()
						if body, err := ioutil.ReadAll(r.Body); err != nil {
							w.WriteHeader(http.StatusBadRequest)
							fmt.Fprint(w, err)
						} else {
							db.People[path[1]] = string(body)
							db.Revisions[path[1]] = createRandomString()
							w.WriteHeader(http.StatusCreated)
							docResp := &docResp{ID: path[1],
								Ok:  true,
								Rev: db.Revisions[path[1]]}
							docRespStr, _ := json.Marshal(docResp)
							fmt.Fprint(w, string(docRespStr))
						}
					}
				}
			case "DELETE":
				if len(path) > 1 {
					if _, ok := db.People[path[1]]; ok {
						delete(db.People, path[1])
						w.WriteHeader(http.StatusOK)
					}
				}
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func createRandomDbCouchUninitialized() (*couchdb, *httptest.Server, error) {
	server := getMockCouchDbServer(new(MockCouchDb))

	host, port := splitURL(server.URL)

	db := new(couchdb)
	db.dbname = createRandomString()
	db.hostname = host
	db.port = port
	db.url = server.URL

	return db, server, nil
}

func createRandomDbCouch() (*couchdb, *httptest.Server, error) {
	server := getMockCouchDbServer(new(MockCouchDb))

	host, port := splitURL(server.URL)

	db := new(couchdb)
	db.Init(createRandomString(), host, port)

	return db, server, nil
}

func TestCouchDb_Req(t *testing.T) {
	db, server, _ := createRandomDbCouchUninitialized()
	defer server.Close()
	person, _ := createRandomPerson()

	var req request
	req.command = "HEAD"
	req.path = db.dbname

	// Good values
	/// Without a person
	if _, err := db.req(&req); err != nil {
		t.Fatalf("Unexpectedly encountered error: %s", err)
	}

	/// With person
	req.person = person
	if _, err := db.req(&req); err != nil {
		t.Fatalf("Unexpectedly encountered error: %s", err)
	}

	// Bad values
	/// Bad command
	req.command = "❤"
	if _, err := db.req(&req); err == nil {
		t.Error("Did not encounter expected error on bad HTTP method")
	}

	/// Bad hostname
	req.command = "HEAD"
	db.url = ""
	if r, err := db.req(&req); err == nil {
		t.Error("Did not encounter expected error with blank hostname")
		t.Logf("Code: %d\n", r.StatusCode)
	}

	// Simulate not having a network connection
	db, server2, _ := createRandomDbCouchUninitialized()
	server2.Close()
	if _, err := db.req(&req); err == nil {
		t.Fatal("Did not receive expected connection error")
	}
}

func TestCouchDb_DbCreate(t *testing.T) {
	db, server, _ := createRandomDbCouchUninitialized()

	// Create new
	if err := db.dbCreate(); err != nil {
		t.Fatalf("Did not create database: %s", err)
	}

	// Attempt to create from blank name
	db.dbname = ""
	if err := db.dbCreate(); err == nil {
		t.Fatal("Unexpectedly created database with a blank name")
	}
	db.dbname = createRandomString()

	// Let's fail on network connectivity
	server.Close()
	if err := db.dbCreate(); err == nil {
		t.Fatal("Unexpectedly created database with a no connection to the server")
	}
}

func TestCouchDb_DbExists(t *testing.T) {
	db, server, _ := createRandomDbCouchUninitialized()

	// Check that Db doesn't already exist
	if db.dbExists() {
		t.Fatal("Non-existent database comes back as existent")
	}

	// Test if we can find created database
	if err := db.dbCreate(); err != nil {
		t.Fatal("Unexpectedly failed at creating a database")
	}
	if !db.dbExists() {
		t.Fatal("Did not find database which was created")
	}

	// Let's fail on network connectivity
	server.Close()
	if db.dbExists() {
		t.Fatal("Unexpectedly created database with a no connection to the server")
	}
}

func TestCouchDb_PersonPath(t *testing.T) {
	db, server, _ := createRandomDbCouch()
	defer server.Close()

	id := createRandomString()
	expectedPath := db.dbname + "/" + id
	if path := db.personPath(id); path != expectedPath {
		t.Fatalf("Expected path '%s' but found '%s'", expectedPath, path)
	}

	expectedPath = db.dbname + "/"
	if path := db.personPath(""); path != expectedPath {
		t.Fatalf("Expected path '%s' but found '%s'", expectedPath, path)
	}
}

func TestCouchDb_Init(t *testing.T) {
	server := getMockCouchDbServer(new(MockCouchDb))

	host, port := splitURL(server.URL)
	dbname := createRandomString()

	db := new(couchdb)

	// Straight up init
	if err := db.Init(dbname, host, port); err != nil {
		t.Fatalf("Error when initializing the database: %s", err)
	}

	// Remove the scheme
	if err := db.Init(dbname, host[7:], port); err != nil {
		t.Fatalf("Error when initializing the database with no scheme in the URL: %s", err)
	}

	// Blank out the fields
	if err := db.Init("", host, port); err == nil {
		t.Error("Database unexpectedly initialized with empty name")
	}

	if err := db.Init(dbname, "", port); err == nil {
		t.Error("Database unexpectedly initialized with empty hostname")
	}

	// Test for whitespace
	if err := db.Init(" ", host, port); err == nil {
		t.Error("Database unexpectedly initialized with name as a whitespace character")
	}

	if err := db.Init(dbname, " ", port); err == nil {
		t.Error("Database unexpectedly initialized with hostname as a whitespace character")
	}

	// Simulate connectivity error
	server.Close()
	db = new(couchdb)
	if err := db.Init(dbname, host, port); err == nil {
		t.Error("Unexpectedly initialized the database without error when there was no connectivity")
	}
}

func TestCouchDb_getRevisionID(t *testing.T) {
	db, server, _ := createRandomDbCouch()

	// Update a non-existent person
	expectedPerson, _ := createRandomPerson()
	revID, err := db.getRevisionID(*expectedPerson)
	if err != nil {
		t.Fatalf("Encountered error when retrieving revision id: %s", err)
	} else if revID != "" {
		t.Error("Unexpectedly received a revision id for a person which does not exist in the database")
	}

	// TODO: what about when it does find it?!

	// Simulate loosing network connectivity
	server.Close()
	if _, err := db.getRevisionID(*expectedPerson); err == nil {
		t.Error("Unexpectedly updated a person without network connectivity")
	}
}

func TestCouchDb_Create(t *testing.T) {
	db, server, _ := createRandomDbCouch()

	// Create a person
	person, _ := createRandomPerson()
	if err := db.Create(*person); err != nil {
		t.Fatalf("Encountered error when creating a new person: %s", err)
	}

	// Create the same person again
	if err := db.Create(*person); err != nil {
		t.Fatalf("Encountered error when creating a person a second time: %s", err)
	}

	person.ID = ""
	if err := db.Create(*person); err == nil {
		t.Fatal("Unexpectedly created a person with a blank id")
	}

	// Simulate loosing network connectivity
	server.Close()
	person, _ = createRandomPerson()
	if err := db.Create(*person); err == nil {
		t.Error("Unexpectedly created a new person without network connectivity")
	}
}

func TestCouchDb_Exists(t *testing.T) {
	db, server, _ := createRandomDbCouch()

	if db.Exists(createRandomString()) {
		t.Fatal("Unexpectedly found person with random id which should not be in the database")
	}

	person, _ := createRandomPerson()
	if err := db.Create(*person); err != nil {
		t.Fatalf("Encountered error when creating a new person: %s", err)
	}

	if !db.Exists(person.ID) {
		t.Fatal("Did not find person which exists in the database")
	}

	// Simulate connectivity error
	person, _ = createRandomPerson()
	if err := db.Create(*person); err != nil {
		t.Fatalf("Encountered error when creating a new person: %s", err)
	}
	server.Close()

	if db.Exists(person.ID) {
		t.Fatal("Found person in the database even though there is no connectivity")
	}
}

func TestCouchDb_Get(t *testing.T) {
	db, server, _ := createRandomDbCouch()

	// Retrieve a non-existent person
	if _, err := db.Get(createRandomString()); err == nil {
		t.Error("Retrieved Person object from empty database")
	}

	// Create a person and retrieve it
	expectedPerson, _ := createRandomPerson()
	if err := db.Create(*expectedPerson); err != nil {
		t.Fatalf("Encountered error when creating a new person: %s", err)
	}

	if person, err := db.Get(expectedPerson.ID); err != nil {
		t.Fatalf("Encountered error when retrieving person: %s", err)
	} else if !arePersonEqual(expectedPerson, person) {
		t.Fatal("Retrieved Person is not equivalent to the expected Person")
	}

	// Get a bad person
	if _, err := db.Get("BADPERSON"); err == nil {
		t.Fatal("Did not receive expected error when retrieving a bad Person object")
	}

	// Simulate connectivity error
	server.Close()
	if _, err := db.Get(expectedPerson.ID); err == nil {
		t.Fatal("Unexpectedly retrieved person with connectivity error")
	}
}

func TestCouchDb_updateWithRevision(t *testing.T) {
	db, server, _ := createRandomDbCouch()

	// Update a non-existent person
	expectedPerson, _ := createRandomPerson()
	if err := db.updateWithRevision(*expectedPerson, ""); err != nil {
		t.Fatalf("Encountered error when 'updating' a new person: %s", err)
	}

	expectedRevID, err := db.getRevisionID(*expectedPerson)
	if err != nil {
		t.Fatalf("Encountered error when retrieving revision id: %s", err)
	}

	// Update the same person again
	expectedName := createRandomString()
	expectedPerson.Name = expectedName
	if err := db.updateWithRevision(*expectedPerson, expectedRevID); err != nil {
		t.Fatalf("Encountered error when updating an existent person: %s", err)
	}

	if person, err := db.Get(expectedPerson.ID); err != nil {
		t.Fatalf("Encountered error when retrieving person: %s", err)
	} else if person.Name != expectedName {
		t.Fatalf("Expected name to have changed to '%s' but found '%s'", expectedName, person.Name)
	}

	expectedRevID, err = db.getRevisionID(*expectedPerson)
	if err != nil {
		t.Fatalf("Encountered error when retrieving revision id: %s", err)
	}

	expectedPerson.ID = ""
	if err := db.updateWithRevision(*expectedPerson, expectedRevID); err == nil {
		t.Fatal("Unexpectedly updated a person with a blank id")
	}

	// t.Skip("TODO")

	// Simulate loosing network connectivity
	server.Close()
	if err := db.updateWithRevision(*expectedPerson, expectedRevID); err == nil {
		t.Error("Unexpectedly updated a person without network connectivity")
	}
}

func TestCouchDb_Update(t *testing.T) {
	t.Skip("TODO")
	db, server, _ := createRandomDbCouch()

	// Update a non-existent person
	expectedPerson, _ := createRandomPerson()
	if err := db.Update(*expectedPerson); err != nil {
		t.Fatalf("Encountered error when 'updating' a new person: %s", err)
	}

	// Update the same person again
	expectedName := createRandomString()
	expectedPerson.Name = expectedName
	if err := db.Update(*expectedPerson); err != nil {
		t.Fatalf("Encountered error when updating an existent person: %s", err)
	}

	if person, err := db.Get(expectedPerson.ID); err != nil {
		t.Fatalf("Encountered error when retrieving person: %s", err)
	} else if person.Name != expectedName {
		t.Fatalf("Expected name to have changed to '%s' but found '%s'", expectedName, person.Name)
	}

	expectedPerson.ID = ""
	if err := db.Update(*expectedPerson); err == nil {
		t.Fatal("Unexpectedly updated a person with a blank id")
	}

	// Simulate loosing network connectivity
	server.Close()
	if err := db.Update(*expectedPerson); err == nil {
		t.Error("Unexpectedly updated a person without network connectivity")
	}
}

func TestCouchDb_Remove(t *testing.T) {
	db, server, _ := createRandomDbCouch()

	// Create a person
	expectedPerson, _ := createRandomPerson()
	if err := db.Create(*expectedPerson); err != nil {
		t.Fatalf("Encountered error when creating a new person: %s", err)
	}

	// Verify that they exist in the database
	if person, err := db.Get(expectedPerson.ID); err != nil {
		t.Fatalf("Encountered error when retrieving person: %s", err)
	} else if !arePersonEqual(expectedPerson, person) {
		t.Fatal("Retrieved Person is not equivalent to the expected Person")
	}

	// Remove that person
	if err := db.Remove(expectedPerson.ID); err != nil {
		t.Fatalf("Encountered error when removing a person from the database: %s", err)
	}

	// Verify that they no longer exist in the database
	if person, err := db.Get(expectedPerson.ID); err == nil {
		t.Error("Unexpectedly did not receive an error when retrieving a removed person")
		if arePersonEqual(expectedPerson, person) {
			t.Fatal("Person was not removed")
		}
	}

	// Remove nonexistant person
	if err := db.Remove(createRandomString()); err == nil {
		t.Error("Unexpectedly did not receive an error when retrieving a person not in the database")
	}

	// Simulate connectivity error
	// Create a person
	expectedPerson, _ = createRandomPerson()
	if err := db.Create(*expectedPerson); err != nil {
		t.Fatalf("Encountered error when creating a new person: %s", err)
	}

	server.Close()

	// Remove that person
	if err := db.Remove(expectedPerson.ID); err == nil {
		t.Error("Did not receive error when attempting to remove Person with connectivity problems")
	}
}
