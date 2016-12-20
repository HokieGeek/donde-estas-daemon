package dondeestas

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

type couchdb struct {
	dbname   string
	hostname string
	port     int
	url      string
}

type request struct {
	command string
	path    string
	person  *Person
	headers map[string]string
}

func (db couchdb) req(params *request) (*http.Response, error) {
	var req *http.Request
	var err error
	if params.person != nil {
		p, _ := json.Marshal(params.person)
		data := bytes.NewBuffer(p)
		req, err = http.NewRequest(params.command, db.url+"/"+params.path, data)
		if err == nil {
			req.Header.Set("Content-Length", fmt.Sprintf("%d", data.Len()))
		}
	} else {
		req, err = http.NewRequest(params.command, db.url+"/"+params.path, nil)
	}
	if err != nil {
		return nil, err
	}
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if bytes, err := httputil.DumpResponse(resp, true); err == nil {
		log.Println("::Response Begin::")
		log.Println(string(bytes))
		log.Println("::Response End::")
	}

	return resp, nil
}

func (db couchdb) dbExists() bool {
	resp, err := db.req(&request{"HEAD", db.dbname, nil, nil})
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}

func (db couchdb) dbCreate() (bool, error) {
	if db.dbname == "" {
		return false, errors.New("Database name is blank")
	}

	if _, err := db.req(&request{"PUT", db.dbname, nil, nil}); err != nil {
		return false, err
	}

	return true, nil
}

func (db couchdb) personPath(id string) string {
	var buf bytes.Buffer
	buf.WriteString(db.dbname)
	buf.WriteString("/")
	buf.WriteString(id)
	return buf.String()
}

func (db *couchdb) Init(dbname, hostname string, port int) error {
	log.Printf("Init(%s, %s, %d)", dbname, hostname, port)
	// TODO: dbname and hostname cannot b e whitespace
	if len(dbname) == 0 {
		return errors.New("No database name specified")
	}

	if len(hostname) == 0 {
		return errors.New("Hostname not specified")
	}

	if port < 0 {
		return errors.New("Invalid port number")
	}

	db.dbname = dbname
	db.hostname = hostname
	db.port = port

	var url bytes.Buffer
	if len(db.hostname) < 4 || db.hostname[:4] != "http" {
		url.WriteString("http://")
	}
	url.WriteString(db.hostname)
	if db.port > -1 {
		url.WriteString(":")
		url.WriteString(fmt.Sprintf("%d", db.port))
	}
	db.url = url.String()

	if db.dbExists() {
		log.Printf("Found database: %s\n", db.dbname)
	} else {
		if ok, err := db.dbCreate(); !ok {
			if err != nil {
				return err
			}
		}
		log.Printf("Created database: %s\n", db.dbname)
	}

	return nil
}

func (db couchdb) Create(p Person) error {
	log.Printf("Create(%s)\n", p.ID)
	return db.Update(p)
}

func (db couchdb) Exists(id string) bool {
	resp, err := db.req(&request{"HEAD", db.personPath(id), nil, nil})
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}

func (db couchdb) Get(id string) (*Person, error) {
	resp, err := db.req(&request{"GET", db.personPath(id), nil, nil})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Encountered an unknown error")
	}

	p := new(Person)
	if err := readCloserJSONToStruct(resp.Body, p); err != nil {
		return nil, err
	}

	return p, nil
}

type docResp struct {
	ID  string `json:"id"`  // Document ID
	Ok  bool   `json:"ok"`  // Operaion Status
	Rev string `json:"rev"` // Revision MVCC token
}

func (db couchdb) getRevisionId(p Person) (*http.Response, error) {
	var req request
	req.command = "HEAD"
	req.path = db.personPath(p.ID)

	resp, err := db.req(&req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (db couchdb) Update(p Person) error {
	log.Printf("Update(%s)\n", p.ID)

	var req request
	if resp, err := db.getRevisionId(p); err != nil {
		return err
	} else if resp.StatusCode == http.StatusOK {
		req.headers = make(map[string]string)
		req.headers["If-Match"] = resp.Header.Get("Etag")
	}

	// Perform the update
	req.command = "PUT"
	req.path = db.personPath(p.ID)
	req.person = &p
	resp, err := db.req(&req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Encountered an unexpected database error: %d", resp.StatusCode)
	}

	/*
		// TODO: use this
		test := new(DocResp)
		if err := readCloserJSONToStruct(resp.Body, test); err != nil {
			return err
		}

		log.Println("Update response:")
		log.Printf("%+v\n", test)
	*/

	return nil
}

func (db couchdb) Remove(id string) error {
	resp, err := db.req(&request{"DELETE", db.personPath(id), nil, nil})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return errors.New("Encountered an unknown error")
	}

	return nil
}
