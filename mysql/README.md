Mysql Package
====
[![godoc reference](https://godoc.org/github.com/mapleque/kelp/mysql?status.svg)](http://godoc.lcgc.work/pkg/github.com/mapleque/kelp/mysql)

本组件基于`go-sql-driver/mysql`封装了在服务中使用数据库的简单接口。

除了基本的数据库操作接口封装外，本组件还实现了：
- 多数据库的注册使用
- 数据与实体的绑定
- 简单的SQL生成
- 用于测试的数据库初始化


初始化
----

`mysql.AddDB`方法支持通过简单配置参数，初始化数据库：
- name 用于在服务内获取数据库连接
- dsn 包含数据库用户名、密码、主机、端口、数据库名等
- max connection 最大连接数，在当前数据库连接池内可能建立的最大连接数
- max idle connection 最大空闲连接数，当连接空闲时，将根据这个设置决定是否立即回收连接

例如：
```
mysql.AddDB(
    "db1", // name
    "user:password@tcp(localhost:3306)/database?charset=utf8", // dsn
    1, // max connection
    1, // max idle connection
)

conn := mysql.GetConnector("db1")

```

在一个服务内，可以多次调用`AddDB`方法初始化多个数据库，这样每个数据库都有独立的连接池，且不互相影响。
在使用时通过调用`GetCOnnector`方法即可获得对应数据库的连接。

Api
----
通过下面定义的API，实现数据库的操作：

```
type Connector interface {
	Begin() (Connector, error)  // start an transaction
	Commit() error              // commit transaction
	Rollback() error            // rollback transaction
	Query(destList interface{}, sql string, params ...interface{}) error
	QueryOne(destObject interface{}, sql string, params ...interface{}) error
	Insert(sql string, params ...interface{}) (lastInsertId int64, err error)
	Execute(sql string, params ...interface{}) (affectRows int64, err error)
}
```

其中`Query`和`QueryOne`方法都需要传入一个目标对象实例指针用于存储返回的数据。

辅助方法
----

- where 用于生成`WHERE`子句

```
// 首先创建一个where对象，并添加对应的条件
where := mysql.NewWhere()
where.Add("id", 1)
where.Add("name", "kelp")

// 然后使用where对象提供的`Sql`方法获得`WHERE`子句，使用`Bind`方法获得条件参数
conn.Query(dest, "SELECT * FROM kelp"+where.Sql(), where.Bind()...)

// 该查询的sql为：SELECT * FROM kelp WHERE id = ? AND name = ?
// 参数为：1, kelp
```

- in 用于生成`IN`子句

```
// 首先创建一个in对象，并添加对应的条件
in := mysql.NewIn()
in.Add("id", 1)
in.Add("id", 2)

// 然后使用in对象提供的`Sql`方法和`Bind`方法查询
conn.Query(dest, "SELECT * FROM kelp"+in.Sql(), in.Bind()...)

// 该查询的sql为：SELECT * FROM kelp WHERE id IN (?,?)
// 参数为：1, 2
```

- sorter 用于生成排序子句
```
// 首先需要自己定义一个实现了`SorterElement`接口的对象
type MySorter struct {
	field   string
	reverse bool
}

func (this *MySorter) GetField() string {
	return this.field
}
func (this *MySorter) GetReverse() bool {
	return this.reverse
}

// 然后创建sorter对象，并添加条件
sorter := mysql.NewSorter()
sorter.Add(&MySorter{"id", true})
sorter.Add(&MySorter{"name", false})

// 然后使用sorter对象提供的`Sql`方法查询
conn.Query(dest, "SELECT * FROM kelp"+sorter.Sql())

// 该查询的sql为：SELECT * FROM kelp ORDER BY id DESC, name ASC
```

- pager 用于生成分页子句（注意当目标数据量较大时谨慎使用）

```
// 可以自定义实现了`PagerElement`接口的对象，也可以直接使用`mysql.Pager`对象
type MyPager struct {
	size   int64
	offset int64
}

func (this *MyPager) GetSize() int64 {
	return this.size
}

func (this *MyPager) GetOffset() int64 {
	return this.offset
}

pager := mysql.NewPager()
pager.Add(&MyPager{20,40})

// 然后使用sorter对象提供的`Sql`方法查询
// conn.Query(dest, "SELECT * FROM kelp"+pager.Sql())
```

测试数据库初始化
----
`mysql.InitTestDB`方法用于初始化一个测试数据库，在初始化过程中会根据`sourceDirs`所指定的路径读取初始化sql文件，并执行之。

例如：
```
mysql.InitTestDB(
  "testdb", // name
  "user:password@tcp(localhost:3306)/database?charset=utf8", // dsn
  "./sql:./test/sql", // sourceDirs
)
```
这样就会在初始化的时候读取并执行`./sql`和`./test/sql`下的所有`.sql`文件。
