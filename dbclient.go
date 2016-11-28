package dondeestas

type DbClientTypes int

const (
	CouchDB DbClientTypes = 0 + iota
)

type dbclient interface {
	Init(dbname, hostname string, port int) error
	Create(p Person) error
	Exists(id int) bool
	Get(id int) (*Person, error)
	Update(p Person) error
	Remove(id int) error
}

type DbClientParams struct {
	DbType           DbClientTypes
	DbName, Hostname string
	Port             int
}

func NewDbClient(params DbClientParams) (*dbclient, error) {
	var db dbclient

	switch params.DbType {
	case CouchDB:
		couch := new(couchdb)
		db = dbclient(couch)
	}

	if err := db.Init(params.DbName, params.Hostname, params.Port); err != nil {
		return nil, err
	}

	return &db, nil
}
