package dondeestas

import "errors"

// DbClient is the interface used by all structs which provide access to a database
type DbClient interface {
	Init(dbname, hostname string, port uint16) error
	Create(p Person) error
	Exists(id string) bool
	Get(id string) (*Person, error)
	Update(p Person) error
	Remove(id string) error
}

// DbClientTypes is an enumeration of database types supported by this library
type DbClientTypes int

// Enumeration which specifies which type of client to create when calling NewDbClient
const (
	CouchDB DbClientTypes = 0 + iota // CouchDB client type
)

// DbClientParams encapsulates the options available for NewDbClient
type DbClientParams struct {
	DbType           DbClientTypes
	DbName, Hostname string
	Port             uint16
}

// NewDbClient creates a database client of specified type at a specified URL
func NewDbClient(params DbClientParams) (*DbClient, error) {
	var db DbClient

	switch params.DbType {
	case CouchDB:
		db = DbClient(new(couchdb))
	default:
		return nil, errors.New("Did not recognize the client type")
	}

	return &db, db.Init(params.DbName, params.Hostname, params.Port)
}
