package mysql

type TestDB struct{}

func NewTestDB() *TestDB {
	return &TestDB{}
}

func (this *TestDB) Begin() (Connector, error) {
	log.Debug("begin a transaction")
	return this, nil
}
func (this *TestDB) Commit() error {
	log.Debug("commit a transaction")
	return nil
}
func (this *TestDB) Rollback() error {
	log.Debug("rollback a transation")
	return nil
}
func (this *TestDB) Query(destList interface{}, sql string, params ...interface{}) error {
	log.Debug("do query", sql, params)
	return nil
}
func (this *TestDB) QueryOne(destObject interface{}, sql string, params ...interface{}) error {
	log.Debug("do query one", sql, params)
	return nil
}
func (this *TestDB) Insert(sql string, params ...interface{}) (int64, error) {
	log.Debug("do insert", sql, params)
	return 1, nil
}
func (this *TestDB) Execute(sql string, params ...interface{}) (int64, error) {
	log.Debug("do execute", sql, params)
	return 1, nil
}
