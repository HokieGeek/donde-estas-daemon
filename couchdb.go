package dondeestas

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

type couchdb struct {
	dbname, hostname string
	port             uint16
	url              string
	personPaths      map[string]string
}

type request struct {
	command, path string
	person        *Person
	headers       map[string]string
}

func (db couchdb) req(params *request) (resp *http.Response, err error) {
	var data bytes.Buffer
	if params.person != nil {
		p, _ := json.Marshal(params.person)
		data = *bytes.NewBuffer(p)
	}
	req, err := http.NewRequest(params.command, db.url+"/"+params.path, &data)
	if err != nil {
		return
	}
	req.Header.Set("Content-Length", fmt.Sprintf("%d", data.Len()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if params.headers != nil {
		for k, v := range params.headers {
			req.Header.Set(k, v)
		}
	}

	if bytes, err := httputil.DumpRequest(req, true); err == nil {
		log.Println("::Request Begin::")
		log.Println(string(bytes))
		log.Println("::Request End::")
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if bytes, err := httputil.DumpResponse(resp, true); err == nil {
		log.Println("::Response Begin::")
		log.Println(string(bytes))
		log.Println("::Response End::")
	}

	return
}

func (db couchdb) dbExists() bool {
	resp, err := db.req(&request{command: "HEAD", path: db.dbname})
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}

func (db couchdb) dbCreate() error {
	if db.dbname == "" {
		return errors.New("Database name is blank")
	}

	if _, err := db.req(&request{command: "PUT", path: db.dbname}); err != nil {
		return err
	}

	return nil
}

func (db couchdb) personPath(id string) string {
	if path, ok := db.personPaths[id]; ok {
		return path
	}

	var buf bytes.Buffer
	buf.WriteString(db.dbname)
	buf.WriteString("/")
	buf.WriteString(id)
	db.personPaths[id] = buf.String()

	return db.personPaths[id]
}

func (db *couchdb) Init(dbname, hostname string, port uint16) error {
	log.Printf("Init(%s, %s, %d)", dbname, hostname, port)

	if len(strings.TrimSpace(dbname)) == 0 {
		return errors.New("No database name specified")
	}

	if len(strings.TrimSpace(hostname)) == 0 {
		return errors.New("Hostname not specified")
	}

	db.personPaths = make(map[string]string)
	db.dbname = dbname
	db.hostname = hostname
	db.port = port

	var url bytes.Buffer
	if len(db.hostname) < 4 || db.hostname[:4] != "http" {
		url.WriteString("http://")
	}
	url.WriteString(fmt.Sprintf("%s:%d", db.hostname, db.port))
	db.url = url.String()

	if db.dbExists() {
		log.Printf("Found database: %s\n", db.dbname)
		return nil
	}

	log.Printf("Creating database: %s\n", db.dbname)
	return db.dbCreate()
}

func (db couchdb) Create(p Person) error {
	log.Printf("Create(%s)\n", p.ID)
	return db.Update(p)
}

func (db couchdb) Exists(id string) bool {
	resp, err := db.req(&request{command: "HEAD", path: db.personPath(id)})
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}

func (db couchdb) Get(id string) (*Person, error) {
	resp, err := db.req(&request{command: "GET", path: db.personPath(id)})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Encountered an unknown error")
	}

	person := new(Person)
	return person, readCloserJSONToStruct(resp.Body, person)
}

func (db couchdb) getRevisionID(p Person) (string, error) {
	resp, err := db.req(&request{command: "HEAD", path: db.personPath(p.ID)})
	if err != nil {
		return "", err
	}

	var revID string
	switch resp.StatusCode {
	case http.StatusOK, http.StatusNotModified:
		revID = resp.Header.Get("Etag")
	}

	return revID, nil
}

type docResp struct {
	ID  string `json:"id"`  // Document ID
	Ok  bool   `json:"ok"`  // Operaion Status
	Rev string `json:"rev"` // Revision MVCC token
}

func (db couchdb) updateWithRevision(p Person, revID string) error {
	req := &request{command: "PUT", path: db.personPath(p.ID), person: &p}
	if revID != "" {
		req.headers = make(map[string]string)
		req.headers["If-Match"] = revID
	}

	// Perform the update
	resp, err := db.req(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Encountered an unexpected database error: %d", resp.StatusCode)
	}

	/*
		// TODO: use this
		test := new(docResp)
		if err := readCloserJSONToStruct(resp.Body, test); err != nil {
			return err
		}

		log.Println("Update response:")
		log.Printf("%+v\n", test)
	*/

	return nil
}

func (db couchdb) Update(p Person) error {
	log.Printf("Update(%s)\n", p.ID)

	revID, err := db.getRevisionID(p)
	if err != nil {
		return err
	}

	return db.updateWithRevision(p, revID)
}

func (db couchdb) Remove(id string) error {
	resp, err := db.req(&request{command: "DELETE", path: db.personPath(id)})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return errors.New("Encountered an unknown error")
	}

	return nil
}
