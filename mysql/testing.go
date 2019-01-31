package mysql

import (
	"io/ioutil"
	"os"
	"strings"
)

// sourceDirs 按照分号或者冒号分割多个dir
func InitTestDB(name, dsn, sourceDirs string) {
	if err := AddDB(
		name,
		dsn,
		10,
		10,
	); err != nil {
		panic(err)
	}
	for _, dir := range strings.FieldsFunc(sourceDirs, func(r rune) bool {
		switch r {
		case ':', ';':
			return true
		}
		return false
	}) {
		log.Debug("source dir", name, dir)
		sourceSqlFiles(name, dir)
	}
}

func sourceSqlFiles(dbname, dir string) {
	sqlDir, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, fi := range sqlDir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), ".SQL") {
			sqlfile := dir + string(os.PathSeparator) + fi.Name()
			sourceSql, err := ioutil.ReadFile(sqlfile)
			if err != nil {
				panic(err)
			}
			for _, sql := range strings.Split(string(sourceSql), ";") {
				if len(strings.TrimSpace(sql)) > 0 {
					if _, err := Get(dbname).Execute(sql); err != nil {
						panic(err)
					}
				}
			}
		}
	}
}
