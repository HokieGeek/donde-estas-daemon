package dondeestas

type DbClientTypes int

const (
	CouchDB DbClientTypes = 0 + iota
)

type dbclient interface {
	Init() error
	Create(p Person) error
	Get(id int) (*Person, error)
	Update(p Person) error
	Remove(id int) error
}

func NewDbClient(dbtype DbClientTypes) (*dbclient, error) {
	var db dbclient

	switch dbtype {
	case CouchDB:
		couch := new(couchdb)
		db = dbclient(couch)
	}

	if err := db.Init(); err != nil {
		return nil, err
	}

	return &db, nil
}
