package mysql

import (
	"fmt"
	"os"
)

type Schema struct {
	tables []*Table
}

func NewSchema() *Schema {
	return &Schema{
		tables: []*Table{},
	}
}

func (this *Schema) Add(model interface{}) *Table {
	table := NewTable(model)
	this.tables = append(this.tables, table)
	return table
}

func (this *Schema) ToFiles(filepath string) {
	if len(filepath) > 0 {
		fmt.Println("start generation sql files")
		defer fmt.Println("done")
	}
	for _, table := range this.tables {
		if err := writeFile(
			filepath,
			table.getSqlFilename(),
			table.getCreateTableSql(),
		); err != nil {
			panic(err)
		}
	}
}

func writeFile(filepath, filename, content string) error {
	var file *os.File
	var err error
	if filepath == "" {
		file = os.Stdout
	} else {
		abfile := filepath + filename
		file, err = os.Create(abfile)
		if err != nil {
			return err
		}
		defer file.Close()
		fmt.Printf("create file: %s\n", abfile)
	}
	fmt.Fprint(file, content)
	return nil
}
