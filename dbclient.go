package dondeestas

// DbClientTypes is an enumeration of database types supported by this library
type DbClientTypes int

const (
	couchDB DbClientTypes = 0 + iota
)

type dbclient interface {
	Init(dbname, hostname string, port int) error
	Create(p Person) error
	Exists(id string) bool
	Get(id string) (*Person, error)
	Update(p Person) error
	Remove(id string) error
}

// DbClientParams encapsulates the options available for newDbClient
type DbClientParams struct {
	DbType           DbClientTypes
	DbName, Hostname string
	Port             int
}

func newDbClient(params DbClientParams) (*dbclient, error) {
	var db dbclient

	switch params.DbType {
	case couchDB:
		couch := new(couchdb)
		db = dbclient(couch)
	}

	if err := db.Init(params.DbName, params.Hostname, params.Port); err != nil {
		return nil, err
	}

	return &db, nil
}
