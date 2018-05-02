package mysql

import (
	"io/ioutil"
	"os"
	"strings"
)

// TestDB is a mock db connector
// which only log out the sql.
type TestDB struct{}

func NewTestDB() *TestDB {
	return &TestDB{}
}

func (this *TestDB) Begin() (Connector, error) {
	log.Info("begin a transaction")
	return this, nil
}
func (this *TestDB) Commit() error {
	log.Info("commit a transaction")
	return nil
}
func (this *TestDB) Rollback() error {
	log.Info("rollback a transation")
	return nil
}
func (this *TestDB) Query(destList interface{}, sql string, params ...interface{}) error {
	log.Info("do query", sql, params)
	return nil
}
func (this *TestDB) QueryOne(destObject interface{}, sql string, params ...interface{}) error {
	log.Info("do query one", sql, params)
	return nil
}
func (this *TestDB) Insert(sql string, params ...interface{}) (int64, error) {
	log.Info("do insert", sql, params)
	return 1, nil
}
func (this *TestDB) Execute(sql string, params ...interface{}) (int64, error) {
	log.Info("do execute", sql, params)
	return 1, nil
}

// InitTestDB is a mysql helper used for build a test database with sql schema.
// This usually used in unittest to initial a test database with service's table created.
func InitTestDB(name, dsn, schemaDir string) {
	if err := AddDB(
		name,
		dsn,
		10,
		10,
	); err != nil {
		panic(err)
	}
	sqlDir, err := ioutil.ReadDir(schemaDir)
	if err != nil {
		panic(err)
	}
	for _, fi := range sqlDir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), "SQL") {
			sqlfile := schemaDir + string(os.PathSeparator) + fi.Name()
			schemaSql, err := ioutil.ReadFile(sqlfile)
			if err != nil {
				panic(err)
			}
			for _, sql := range strings.Split(string(schemaSql), ";") {
				if len(strings.TrimSpace(sql)) > 0 {
					if _, err := Get(name).Execute(sql); err != nil {
						panic(err)
					}
				}
			}
		}
	}
}
