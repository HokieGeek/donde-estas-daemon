package dondeestas

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

func (db couchdb) req(command, path string, person *Person) (*http.Response, error) {
	var req *http.Request
	var err error
	if person != nil {
		p, err := json.Marshal(*person)
		if err != nil {
			return nil, err
		}
		data := bytes.NewBuffer(p)
		req, err = http.NewRequest(command, db.url+"/"+path, data)
		req.Header.Set("Content-Length", fmt.Sprintf("%d", data.Len()))
	} else {
		req, err = http.NewRequest(command, db.url+"/"+path, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	bytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(bytes))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (db couchdb) createDbIfNotExist() error {
	resp, err := db.req("HEAD", db.dbname, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		log.Printf("Creating database: %s\n", db.dbname)

		_, err := db.req("PUT", db.dbname, nil)
		if err != nil {
			return err
		}
	} else {
		log.Printf("Found database: %s\n", db.dbname)
	}

	return nil
}

func (db couchdb) personPath(id int) string {
	var buf bytes.Buffer
	buf.WriteString(db.dbname)
	buf.WriteString("/")
	buf.WriteString(fmt.Sprintf("%d", id))
	return buf.String()
}

func (db *couchdb) Init(dbname, hostname string, port int) error {
	db.dbname = dbname
	db.hostname = hostname
	db.port = port

	db.url = "http://" + db.hostname + ":" + fmt.Sprintf("%d", db.port)

	err := db.createDbIfNotExist()
	if err != nil {
		return err
	}

	return nil
}

func (db couchdb) Create(p Person) error {
	return db.Update(p)
}

func (db couchdb) Exists(id int) bool {
	resp, err := db.req("HEAD", db.personPath(id), nil)
	if err != nil {
		return false
	}

	return resp.StatusCode == 200
}

func (db couchdb) Get(id int) (*Person, error) {
	resp, err := db.req("GET", db.personPath(id), nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Encountered an unknown error")
	}

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		return nil, err
	}

	if err = resp.Body.Close(); err != nil {
		return nil, err
	}

	p := new(Person)
	if err := json.Unmarshal(body, p); err != nil {
		return nil, err
	}

	return p, nil
}

type DocResp struct {
	Id  string `json:"id"`
	Ok  bool   `json:"ok"`
	Rev string `json:"rev"`
}

func (db couchdb) Update(p Person) error {
	resp, err := db.req("PUT", db.personPath(p.Id), &p)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 && resp.StatusCode != 202 {
		return errors.New(fmt.Sprintf("Encountered an unexpected error: %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		return err
	}

	if err = resp.Body.Close(); err != nil {
		return err
	}

	test := new(DocResp)
	if err := json.Unmarshal(body, test); err != nil {
		return err
	}

	log.Println("Update response:")
	log.Printf("%+v\n", test)

	return nil
}

func (db couchdb) Remove(id int) error {
	resp, err := db.req("DELETE", db.personPath(id), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Encountered an unknown error")
	}

	return nil
}
