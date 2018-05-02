package mysql_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mapleque/kelp/mysql"
)

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
	conn := mysql.Get("test")
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
