package dondeestas

import (
	"net/http/httptest"
	"testing"
)

func createRandomDbClientParams() (DbClientParams, *httptest.Server) {
	db, server, _ := createRandomDbCouchUninitialized()
	// params := DbClientParams{CouchDB, createRandomString(), db.hostname, db.port}
	params := DbClientParams{CouchDB, createRandomString(), db.hostname, db.port} // TODO: random type
	return params, server
}

func createRandomDbClient() (*DbClient, *httptest.Server, error) {
	params, server := createRandomDbClientParams()

	client, err := NewDbClient(params)
	if err != nil {
		return nil, nil, err
	}

	return client, server, nil
}

func TestNewDbClient(t *testing.T) {
	params, server := createRandomDbClientParams()

	if _, err := NewDbClient(params); err != nil {
		t.Fatalf("Error when creating new DbClient: %s", err)
	}

	params.DbName = ""
	if _, err := NewDbClient(params); err == nil {
		t.Error("Unexpectedly created DbClient with empty DB name")
	}
	params.DbName = createRandomString()

	params.DbType = DbClientTypes(42) // TODO
	if _, err := NewDbClient(params); err == nil {
		t.Error("Unexpectedly created DbClient with empty DB name")
	}

	// Simulate no network connectivity
	server.Close()
	if _, err := NewDbClient(params); err == nil {
		t.Error("Unexpetedly created DbClient with no connectivity")
	}
}
