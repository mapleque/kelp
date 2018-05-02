package mysql

// pool is mysql connection pool
// which to hold all mysql connection in a service
type pool struct {
	store map[string]Connector
}

// Connector is mysql database connector inerface
// implement by query and transaction
type Connector interface {
	Begin() (Connector, error)
	Commit() error
	Rollback() error
	Query(destList interface{}, sql string, params ...interface{}) error
	QueryOne(destObject interface{}, sql string, params ...interface{}) error
	Insert(sql string, params ...interface{}) (lastInsertId int64, err error)
	Execute(sql string, params ...interface{}) (affectRows int64, err error)
}

// p used as a connection pool storage
var p *pool

func init() {
	if p != nil {
		return
	}
	p = &pool{make(map[string]Connector)}
}

// Get return a Connector which have been added into pool
func Get(name string) Connector {
	return p.store[name]
}

// Add store a connector into pool
func Add(name string, conn Connector) {
	p.store[name] = conn
}
