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
)

type couchdb struct {
	dbname   string
	hostname string
	port     int
	url      string
}

func (db *couchdb) req(command, path string, person *Person) (*http.Response, error) {
	var data *bytes.Buffer

	if person != nil {
		p, err := json.Marshal(*person)
		if err != nil {
			return nil, err
		}
		data = bytes.NewBuffer(p)
	}

	req, err := http.NewRequest(command, db.url+"/"+path, data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (db *couchdb) createDbIfNotExist() error {
	resp, err := db.req("HEAD", db.dbname, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		log.Println("Creating database")

		_, err := db.req("PUT", db.dbname, nil)
		if err != nil {
			return err
		}

	}

	return nil
}

func (db *couchdb) personPath(id int) string {
	var buf bytes.Buffer
	buf.WriteString(db.dbname)
	buf.WriteString("/")
	buf.WriteString(fmt.Sprintf("%d", id))
	return buf.String()
}

func (db *couchdb) Init(dbname, hostname string, port int) error {
	fmt.Println("Init()")

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

func (db *couchdb) Create(p Person) error {
	fmt.Println("Create(p)")
	return db.Update(p)
}

func (db *couchdb) Exists(id int) bool {
	fmt.Printf("Exists(%d)\n", id)

	resp, err := db.req("HEAD", db.personPath(id), nil)
	if err != nil {
		return false
	}

	return resp.StatusCode == 200
}

func (db *couchdb) Get(id int) (*Person, error) {
	fmt.Printf("Get(%d)\n", id)

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

func (db *couchdb) Update(p Person) error {
	fmt.Println("Update(p)")

	resp, err := db.req("PUT", db.personPath(p.Id), &p)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Encountered an unknown error")
	}

	return nil
}

func (db *couchdb) Remove(id int) error {
	fmt.Printf("Remove(%d)\n", id)

	resp, err := db.req("DELETE", db.personPath(id), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Encountered an unknown error")
	}

	return nil
}
