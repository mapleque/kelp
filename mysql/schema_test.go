package mysql_test

import (
	"../mysql"
)

type MyTable struct {
	Id    int64  `json:"id" column_schema:"INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT"`
	Key   string `json:"json_key" column:"key" column_schema:"VARCHAR(128) NOT NULL" comment:"key"`
	Value string `json:"value" column_schema:"TEXT" comment:"value"`
}

// This example showing how to generation
// a table creating sql from model define
func Example_schema() {
	schema := mysql.NewSchema()

	schema.Add(
		MyTable{},
	).Append(
		"ALTER TABLE `my_table` ADD UNIQUE(`key`);",
	).SetCharset(
		"utf8mb4",
	)

	// Create files into filepath
	// if filepath is empty, output to stdout
	schema.ToFiles("")

	// Output:
	// DROP TABLE IF EXISTS `my_table`;
	// CREATE TABLE `my_table` (
	//	`id` INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
	//	`key` VARCHAR(128) NOT NULL,
	//	`value` TEXT
	// ) DEFAULT CHARSET=utf8mb4;
	// ALTER TABLE `my_table` ADD UNIQUE(`key`);
}
