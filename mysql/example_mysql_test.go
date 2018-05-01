package mysql_test

import (
	"flag"
	"fmt"

	"github.com/mapleque/kelp/mysql"
)

var dsn *string

func init() {
	dsn := flag.String("mysql", "www:www@tcp(127.0.0.1:3306)/test?charset=utf8", "mysql dsn")
	flag.Parse()
	// add db connector into pool
	mysql.AddDB(
		"test",
		*dsn,
		1, 1,
	)
}

type TestModel struct {
	Id    int64
	Value string
}

func Example_mysql() {
	// get db connector from pool
	conn := mysql.Get("test")

	// execute sql
	conn.Execute("DROP TABLE IF EXISTS test_mysql_kelp")
	conn.Execute("CREATE TABLE test_mysql_kelp (" +
		"id INT NOT NULL AUTO_INCREMENT," +
		"value VARCHAR(10) DEFAULT NULL," +
		"PRIMARY KEY (id))")

	// insert returns last id
	lastId, _ := conn.Insert("INSERT INTO test_mysql_kelp (value) VALUES (?)", "test_data")
	fmt.Println("last id is", lastId)

	// execute returns affect rows
	affectRow, _ := conn.Execute("UPDATE test_mysql_kelp SET value = ? WHERE id = 1 LIMIT 1", "test_other")
	fmt.Println("affect rows is", affectRow)

	// query one
	testModel := &TestModel{}
	conn.QueryOne(testModel, "SELECT * FROM test_mysql_kelp WHERE id = 1 LIMIT 1")
	fmt.Println("test model is", testModel)

	// begin a transaction
	trans, _ := conn.Begin()
	// insert in transaction
	trans.Insert("INSERT INTO test_mysql_kelp (value) VALUES (?)", "test_trans")
	trans.Execute("UPDATE test_mysql_kelp SET value = ? WHERE id = 1 LIMIT 1", "test_trans")
	// rollback in transaction
	trans.Rollback()

	testModel = &TestModel{}
	conn.QueryOne(testModel, "SELECT * FROM test_mysql_kelp WHERE id = 1 LIMIT 1")
	fmt.Println("test model is still", testModel)

	// begin a transaction
	trans, _ = conn.Begin()
	// insert in transaction
	trans.Insert("INSERT INTO test_mysql_kelp (value) VALUES (?)", "test_trans")
	// execute in transation
	trans.Execute("UPDATE test_mysql_kelp SET value = ? WHERE id = 1 LIMIT 1", "test_trans")
	// commit transaction
	trans.Commit()

	testModel = &TestModel{}
	conn.QueryOne(testModel, "SELECT * FROM test_mysql_kelp WHERE id = 1 LIMIT 1")
	fmt.Println("test model is change to", testModel)
	// Output:
	// last id is 1
	// affect rows is 1
	// test model is &{1 test_other}
	// test model is still &{1 test_other}
	// test model is change to &{1 test_trans}
}
