Mysql Package
====

Connection
----

```
mysql.AddDB(
    "db1",
    "user:password@tcp(localhost:3306)/database?charset=utf8",
    1, // max connection
    1, // max idle connection
)

conn := mysql.GetConnector("db1")

```

Api
----

```
type Connector interface {
	Begin() (Connector, error)  // start an transaction
	Commit() error              // commit transaction
	Rollback() error            // rollback transaction
	Query(destList interface{}, sql string, params ...interface{}) error
	QueryOne(destObject interface{}, sql string, params ...interface{}) error
	Insert(sql string, params ...interface{}) (lastInsertId int64, err error)
	Execute(sql string, params ...interface{}) (affectRows int64, err error)
}
```

Mock Connector
----

```
type MockConnector struct {}
func (this *MockConnector) Begin() (Connection, error) {
	return this, nil
}
func (this *MockConnector) Commit() error {
	return nil
}
func (this *MockConnector) Rollback() error {
	return nil
}
func (this *MockConnector) Query(destList interface{}, sql string, params ...interface{}) error {
	return nil
}
func (this *MockConnector) QueryOne(destObject interface{}, sql string, params ...interface{}) error {
	return nil
}
func (this *MockConnector) Insert(sql string, params ...interface{}) (int64, error) {
	return 1, nil
}
func (this *MockConnector) Execute(sql string, params ...interface{}) (int64, error) {
	return 1, nil
}

func mock() {
    conn := &MockConnector{}
}
```

Example
----

```
```
