package mysql

import (
	"encoding/json"
	"flag"
	"testing"
	"time"
)

func init() {
	dsn := flag.String("mysql", "www:www@tcp(127.0.0.1:3306)/test?charset=utf8", "mysql dsn")
	flag.Parse()
	AddDB(
		"test",
		*dsn,
		1, 1,
	)
}

type TestModel struct {
	Id    int64
	Value string
}

func TestMysql(t *testing.T) {
	conn := GetConnector("test")
	if _, err := conn.Execute("DROP TABLE IF EXISTS test_mysql_kelp"); err != nil {
		t.Fatal("drop table failed", err)
	}
	if _, err := conn.Execute("CREATE TABLE test_mysql_kelp (" +
		"id INT NOT NULL AUTO_INCREMENT," +
		"value VARCHAR(10) DEFAULT NULL," +
		"PRIMARY KEY (id))"); err != nil {
		t.Fatal("create table failed", err)
	}
	lastId, _ := conn.Insert("INSERT INTO test_mysql_kelp (value) VALUES (?)", "test_data")
	if lastId != 1 {
		t.Fatal("last insert id wrong", lastId)
	}
	ignoreId, _ := conn.Insert("INSERT IGNORE INTO test_mysql_kelp VALUES (1,?)", "test_data")
	if ignoreId != 0 {
		t.Fatal("ignore insert id wrong", ignoreId)
	}
	affectRow, _ := conn.Execute("UPDATE test_mysql_kelp SET value = ? WHERE id = 1 LIMIT 1", "test_other")
	if affectRow != 1 {
		t.Fatal("affect row wrong", affectRow)
	}
	testModel := &TestModel{}
	if err := conn.QueryOne(testModel, "SELECT * FROM test_mysql_kelp WHERE id = 1 LIMIT 1"); err != nil {
		t.Fatal("query error", err)
	}
	if testModel.Value != "test_other" {
		t.Fatal("update failed", testModel)
	}

	trans, _ := conn.Begin()
	trans.Insert("INSERT INTO test_mysql_kelp (value) VALUES (?)", "test_trans")
	trans.Execute("UPDATE test_mysql_kelp SET value = ? WHERE id = 1 LIMIT 1", "test_trans")
	trans.Rollback()
	testModel = &TestModel{}
	if err := conn.QueryOne(testModel, "SELECT * FROM test_mysql_kelp WHERE id = 1 LIMIT 1"); err != nil {
		t.Fatal("query error", err)
	}
	if testModel.Value == "test_trans" {
		t.Fatal("rollback failed", testModel)
	}

	trans, _ = conn.Begin()
	trans.Insert("INSERT INTO test_mysql_kelp (value) VALUES (?)", "test_trans")
	trans.Execute("UPDATE test_mysql_kelp SET value = ? WHERE id = 1 LIMIT 1", "test_trans")
	trans.Commit()
	testModel = &TestModel{}
	if err := conn.QueryOne(testModel, "SELECT * FROM test_mysql_kelp WHERE id = 1 LIMIT 1"); err != nil {
		t.Fatal("query error", err)
	}
	if testModel.Value != "test_trans" {
		t.Fatal("commit failed", testModel)
	}
}

type QueryModel struct {
	Id    int
	Str   string
	Dbl   float64
	Dt    time.Time
	Tsr   string `column:"dt"`
	Tx    string
	Ch    string
	Bl    bool
	Extra string `column:"str"`
}

func TestQuery(t *testing.T) {
	conn := GetConnector("test")
	conn.Execute("DROP TABLE IF EXISTS test_query")
	conn.Execute("CREATE TABLE test_query (" +
		"id INT NOT NULL AUTO_INCREMENT," +
		"str VARCHAR(10) DEFAULT NULL," +
		"it INT DEFAULT NULL," +
		"uit INT UNSIGNED DEFAULT NULL," +
		"dbl DOUBLE(9,2) DEFAULT NULL," +
		"dt DATETIME DEFAULT NULL," +
		"tx TEXT DEFAULT NULL," +
		"ch CHAR(1) DEFAULT NULL," +
		"bl BOOLEAN DEFAULT NULL," +
		"PRIMARY KEY (id))")
	now := time.Now()
	conn.Insert(
		"INSERT INTO test_query (str,it,uit,dbl,dt,tx,ch,bl) VALUES(?,?,?,?,?,?,?,?)",
		"str",
		1,
		2,
		123.12,
		now.Format("2006-01-02 15:04:05"),
		"text",
		"c",
		true,
	)
	var list []*QueryModel
	err := conn.Query(&list, "SELECT * FROM test_query WHERE id = ?", 1)
	if err != nil {
		t.Fatal(err)
	}
	// list should be:
	// [{
	//   "Id":1,
	//   "Str":"str",
	//   "Dbl":123.12,
	//   "Tsr":"2017-10-26 11:08:00",
	//   "Tx":"text",
	//   "Ch":"c",
	//   "Bl":true
	// }]
	if len(list) != 1 ||
		list[0].Id != 1 ||
		list[0].Str != "str" ||
		list[0].Dbl != 123.12 ||
		//!now.Round(time.Second).Equal(list[0].Dt) ||
		list[0].Tsr != now.Format("2006-01-02 15:04:05") ||
		list[0].Tx != "text" ||
		list[0].Ch != "c" ||
		list[0].Bl != true {
		str, _ := json.Marshal(list)
		t.Fatal("query bind wrong", string(str))
	}
	dest := &QueryModel{}
	err = conn.QueryOne(dest, "SELECT * FROM test_query")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 ||
		dest.Id != 1 ||
		dest.Str != "str" ||
		dest.Dbl != 123.12 ||
		//!now.Round(time.Second).Equal(list[0].Dt) ||
		dest.Tsr != now.Format("2006-01-02 15:04:05") ||
		dest.Tx != "text" ||
		dest.Ch != "c" ||
		dest.Bl != true {
		str, _ := json.Marshal(list)
		t.Fatal("query one bind wrong", string(str))
	}
}
