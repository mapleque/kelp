Database package
========

Connection
--------

Add a db
```
db.AddMysql(
    "db1",
    "user:password@tcp(localhost:3306)/database?charset=utf8",
    1, // max connection
    1, // max idle connection
)
```

Create connection
```
conn := db.UseMysql("db1")
```
Do operations, such as:
- conn.Select
- conn.Insert
- conn.Update
- conn.Execute
- conn.Query
- conn.QueryOne

Use transaction
```
trans := conn.Begin()
```
Do operations, such as:
- conn.Select
- conn.Insert
- conn.Update
- conn.Execute
- conn.Query
- conn.QueryOne
Commit
```
trans.Commit()
```
Rollback
```
trans.Rollback()
```

Api
----

Select return a []map[string]interface{} result or error
```
result, err := conn.Select(`SELECT * FROM table_test WHERE id = ?`, id)
```

Insert return the current id
```
id, err := conn.Insert(`INSERT INTO table_test (v) VALUES (?)`, value)
```

Update return the affected row number
```
eff, err := conn.Update(`UPDATE table_test SET v = ? WHERE id = ?`, value, id)
```

Execute return the affected row number
```
eff, err := conn.Execute(`DELETE table_test WHERE id = ?`, id)
```

Query bind result to the param
```
type TableItem struct {
    Id int64
    Value string
}

list := []*TableItem{}
err := conn.Query(&list, `SELECT * FROM table_test WHERE id < ? ORDER BY id DESC LIMIT 10`, lastId)
```

Query one bind result to struct, err return on can not find data
```
item := &TableItem{}
err := conn.QueryOne(item, `SELECT * FROM TABLE_TEST WHERE id = ?`, id)
```
When db bind data to struct, struct tag will be process as follow:
- column tag
- json tag
- struct name lowercase
field is not matched will be ignore.

For example:
```
type ExampleItem {
    Id
    Id1 `json:"id"`
    Id2 `json:"id2" column:"id"`
}
```
All field will be bind id value
