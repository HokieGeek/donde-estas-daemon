package dondeestas

// DbClient is the interface used by all structs which provide access to a database
type DbClient interface {
	Init(dbname, hostname string, port int) error
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
	Port             int
}

// NewDbClient creates a database client of specified type at a specified URL
func NewDbClient(params DbClientParams) (*DbClient, error) {
	var db DbClient

	switch params.DbType {
	case CouchDB:
		couch := new(couchdb)
		db = DbClient(couch)
	}

	if err := db.Init(params.DbName, params.Hostname, params.Port); err != nil {
		return nil, err
	}

	return &db, nil
}
