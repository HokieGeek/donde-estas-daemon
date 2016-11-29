package dondeestas

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getServer(data string) *httptest.Server {
	// TODO
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, data)
	}))

	return ts
}

func TestInit(t *testing.T) {
	// dummyServer := getServer(expectedData)
	// defer dummyServer.Close()
	// dummyServer.URL
}

/*
func (db couchdb) req(command, path string, person *Person) (*http.Response, error) {
func (db couchdb) createDbIfNotExist() error {
func (db couchdb) personPath(id int) string {
func (db *couchdb) Init(dbname, hostname string, port int) error {
func (db couchdb) Create(p Person) error {
func (db couchdb) Exists(id int) bool {
func (db couchdb) Get(id int) (*Person, error) {
func (db couchdb) Update(p Person) error {
func (db couchdb) Remove(id int) error {
*/
