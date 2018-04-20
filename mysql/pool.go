package mysql

type Connector interface {
	Begin() (Connector, error)
	Commit() error
	Rollback() error
	Query(destList interface{}, sql string, params ...interface{}) error
	QueryOne(destObject interface{}, sql string, params ...interface{}) error
	Insert(sql string, params ...interface{}) (lastInsertId int64, err error)
	Execute(sql string, params ...interface{}) (affectRows int64, err error)
}

type _Pool struct {
	pool map[string]Connector
}

var pool *_Pool

func init() {
	if pool != nil {
		return
	}
	pool = &_Pool{make(map[string]Connector)}
}

func GetConnector(name string) Connector {
	return pool.pool[name]
}
