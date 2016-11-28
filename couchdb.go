package dondeestas

import (
	"fmt"
	// "net/http"
	"errors"
)

type couchdb struct {
}

func (db *couchdb) Init() error {
	fmt.Println("Init()")
	return nil
}

func (db *couchdb) Create(p Person) error {
	fmt.Println("Create(p)")
	return errors.New("Not implemented")
}

func (db *couchdb) Get(id int) (*Person, error) {
	fmt.Printf("Get(%d)\n", id)
	return nil, errors.New("Not implemented")
}

func (db *couchdb) Update(p Person) error {
	fmt.Println("Update(p)")
	return errors.New("Not implemented")
}

func (db *couchdb) Remove(id int) error {
	fmt.Printf("Remove(%d)\n", id)
	return errors.New("Not implemented")
}
