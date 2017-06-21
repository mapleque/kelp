package db

import (
	"testing"
)

func TestMysql(t *testing.T) {
	DB.Pool = make(map[string]*DBQuery)
	AddDB("test",
		"www:www@tcp(127.0.0.1:3306)/test?charset=utf8",
		1, 1,
	)
	Execute("test", "DROP TABLE IF EXISTS test_mysql")
	Execute("test", "CREATE TABLE test_mysql ("+
		"id INT NOT NULL AUTO_INCREMENT,"+
		"value VARCHAR(10) DEFAULT NULL,"+
		"PRIMARY KEY (id))")
	lastId := Insert("test", "INSERT INTO test_mysql (value) VALUES (?)", "test_data")
	if lastId != 1 {
		t.Fatal("last insert id wrong", lastId)
	}
	list := Select("test", "SELECT id FROM test_mysql WHERE id = ?", 1)
	if len(list) != 1 {
		t.Fatal("list length wrong", list)
	}
	ignoreId := Insert("test", "INSERT IGNORE INTO test_mysql VALUES (1,?)", "test_data")
	if ignoreId != 0 {
		t.Fatal("ignore insert id wrong", ignoreId)
	}
	affectRow := Update("test", "UPDATE test_mysql SET value = ? WHERE id = 1 LIMIT 1", "test_other")
	if affectRow != 1 {
		t.Fatal("affect row wrong", affectRow)
	}
	res := Select("test", "SELECT * FROM test_mysql WHERE id = 1")
	if len(res) != 1 || res[0]["value"] != "test_other" {
		t.Fatal("update result wrong", res)
	}
}
