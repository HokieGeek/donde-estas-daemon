package dondeestas

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

type DummyCouchDb struct {
	Name   string
	People map[int]string
}

func getTestCouchDbServer(db *DummyCouchDb) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path[1:], "/")
		if len(path) == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			// fmt.Println(r.Method)
			// fmt.Println(path)

			switch r.Method {
			case "GET":
				id, _ := strconv.Atoi(path[1])
				if _, ok := db.People[id]; ok {
					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, db.People[id])
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			case "PUT":
				if len(path) == 1 {
					db.Name = path[0]
					w.WriteHeader(http.StatusCreated)
				} else {
					id, _ := strconv.Atoi(path[1])
					if _, ok := db.People[id]; ok {
						defer r.Body.Close()
						body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							fmt.Fprint(w, err)
						} else {
							db.People[id] = string(body)
							fmt.Println(db.People[id])
							w.WriteHeader(http.StatusCreated)
						}
					} else {
						w.WriteHeader(http.StatusNotFound)
					}
				}
			case "HEAD":
				if len(path) == 1 {
					if path[0] == db.Name {
						w.WriteHeader(http.StatusOK)
					} else {
						w.WriteHeader(http.StatusNotFound)
					}
				} else {
					id, _ := strconv.Atoi(path[1])
					if _, ok := db.People[id]; ok {
						w.WriteHeader(http.StatusOK)
					} else {
						w.WriteHeader(http.StatusNotFound)
					}
				}
			case "DELETE":
				if len(path) >= 1 {
					id, _ := strconv.Atoi(path[1])
					if _, ok := db.People[id]; ok {
						delete(db.People, id)
						w.WriteHeader(http.StatusOK)
					} else {
						w.WriteHeader(http.StatusNotFound)
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
		}
	}))

	return ts
}

func splitUrl(url string) (string, int) {
	sepPos := strings.LastIndex(url, ":")
	p, err := strconv.Atoi(url[sepPos+1:])
	if err != nil {
		// TODO
		return "", sepPos
	}
	return url[:sepPos], p
}

func TestCouchDb_Init(t *testing.T) {
	dummyServer := getTestCouchDbServer(new(DummyCouchDb))
	defer dummyServer.Close()

	host, port := splitUrl(dummyServer.URL)
	dbname := "foobar"

	db := new(couchdb)

	// Straight up init
	if err := db.Init(dbname, host, port); err != nil {
		t.Fatal(err)
	}

	// Remove the scheme
	if err := db.Init(dbname, host[7:], port); err != nil {
		t.Fatal(err)
	}

	// Blank out the fields
	if err := db.Init("", host, port); err == nil {
		t.Error("Database unexpectedly initialized with empty name")
	}

	if err := db.Init(dbname, "", port); err == nil {
		t.Error("Database unexpectedly initialized with empty hostname")
	}

	if err := db.Init(dbname, host, -1); err == nil {
		t.Error("Database unexpectedly initialized with invalid port number")
	}

	// TODO: test for whitespace
}

func TestCouchDb_Req(t *testing.T) {
	t.Skip("TODO")
}
func TestCouchDb_CreateDbIfNotExist(t *testing.T) {
	t.Skip("TODO")
}
func TestCouchDb_PersonPath(t *testing.T) {
	t.Skip("TODO")
}
func TestCouchDb_Exists(t *testing.T) {
	t.Skip("TODO")
}
func TestCouchDb_Get(t *testing.T) {
	t.Skip("TODO")
}
func TestCouchDb_Create(t *testing.T) {
	t.Skip("TODO")
}
func TestCouchDb_Update(t *testing.T) {
	t.Skip("TODO")
}
func TestCouchDb_Remove(t *testing.T) {
	t.Skip("TODO")
}
