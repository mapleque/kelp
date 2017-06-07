package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	. "github.com/kelp/log"
)

// DBPool 类型， 是一个database容器，用于存储服务可能用到的所有db连接池
type DBPool struct {
	Pool map[string]*DBQuery
}

// DBQuery 对象，用于直接查询或执行
type DBQuery struct {
	database string
	conn     *sql.DB
}

// DBTransaction 对象，用于事物查询或执行
type DBTransaction struct {
	database string
	conn     *sql.Tx
}

// DB 全局变量，允许用户可以在任何一个方法中调用并操作数据库
var DB *DBPool

func init() {
	if DB != nil {
		return
	}
	Info("init db module...")
	DB = &DBPool{}
	DB.Pool = make(map[string]*DBQuery)
}

/**
 * AddDB 方法，添加一个database，并且在启动前验证其连通性
 */
func AddDB(name, dsn string, maxOpenConns, maxIdleConns int) {
	Info("add db", name, dsn)
	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		Error("can not open db", dsn, err.Error())
		panic(err.Error())
	}
	//	defer db.Close()
	// 如果这里defer，这里刚添加完的db就会被关掉
	err = dbConn.Ping()
	if err != nil {
		Error("can not ping db", dsn, err.Error())
		panic(err.Error()) // 直接panic，让server无法启动
	}
	dbConn.SetMaxOpenConns(maxOpenConns)
	dbConn.SetMaxIdleConns(maxIdleConns)
	dbQuery := &DBQuery{}
	dbQuery.conn = dbConn
	dbQuery.database = name
	DB.Pool[name] = dbQuery
}

// UserDB 方法，返回DBQuery对象
func UseDB(database string) *DBQuery {
	return DB.Pool[database]
}

// Begin 方法，返回DBTransaction对象
func Begin(database string) *DBTransaction {
	return DB.Pool[database].Begin()
}

// Select 方法，返回查询结果数组
func Select(
	database, sql string,
	params ...interface{}) []map[string]interface{} {

	return DB.Pool[database].Select(sql, params...)
}

// Update 方法，返回受影响行数
func Update(
	database, sql string,
	params ...interface{}) int64 {

	return DB.Pool[database].Update(sql, params...)
}

// Insert 方法，返回插入id
func Insert(
	database, sql string,
	params ...interface{}) int64 {

	return DB.Pool[database].Insert(sql, params...)
}

// Begin 方法，返回DBTransaction对象
func (dbq *DBQuery) Begin() *DBTransaction {
	Debug("[transaction begin]", "["+dbq.database+"]")
	trans := &DBTransaction{}
	conn, err := dbq.conn.Begin()
	if err != nil {
		Error("db create transaction faild", dbq.database, err.Error())
		return trans
	}
	trans.conn = conn
	return trans
}

// Select 方法，返回查询结果数组
func (dbq *DBQuery) Select(
	sql string,
	params ...interface{}) []map[string]interface{} {

	Debug("[select sql]", "["+dbq.database+"]", sql, params)
	ret, err := dbq.conn.Query(sql, params...)
	return processQueryRet(sql, ret, err)
}

// Update 方法，返回受影响行数
func (dbq *DBQuery) Update(
	sql string,
	params ...interface{}) int64 {

	Debug("[update sql]", "["+dbq.database+"]", sql, params)
	ret, err := dbq.conn.Exec(sql, params...)
	return processUpdateRet(sql, ret, err)
}

// Insert 方法，返回插入id
func (dbq *DBQuery) Insert(
	sql string,
	params ...interface{}) int64 {

	Debug("[insert sql]", "["+dbq.database+"]", sql, params)
	ret, err := dbq.conn.Exec(sql, params...)
	return processInsertRet(sql, ret, err)
}

// Select 方法，返回查询结果数组
func (dbt *DBTransaction) Select(
	sql string,
	params ...interface{}) []map[string]interface{} {

	Debug("[select sql in transaction]", "["+dbt.database+"]", sql, params)
	ret, err := dbt.conn.Query(sql, params...)
	return processQueryRet(sql, ret, err)
}

// Update 方法，返回受影响行数
func (dbt *DBTransaction) Update(
	sql string,
	params ...interface{}) int64 {

	Debug("[update sql in transaction]", "["+dbt.database+"]", sql, params)
	ret, err := dbt.conn.Exec(sql, params...)
	return processUpdateRet(sql, ret, err)
}

// Insert 方法，返回插入id
func (dbt *DBTransaction) Insert(
	sql string,
	params ...interface{}) int64 {

	Debug("[insert sql in transaction]", "["+dbt.database+"]", sql, params)
	ret, err := dbt.conn.Exec(sql, params...)
	return processInsertRet(sql, ret, err)
}

// Commit 方法，提交事物
func (dbt *DBTransaction) Commit() {
	Debug("[transaction commit]", "["+dbt.database+"]")
	err := dbt.conn.Commit()
	if err != nil {
		Error("db transaction commit faild", dbt.database, err.Error())
	}
}

// Rollback 方法，回滚事物
func (dbt *DBTransaction) Rollback() {
	Debug("[transaction rollback]", "["+dbt.database+"]")
	err := dbt.conn.Rollback()
	if err != nil {
		Error("db transaction rollback faild", dbt.database, err.Error())
	}
}

// 返回查询结果数组
func processQueryRet(
	query string, rows *sql.Rows, err error) []map[string]interface{} {
	if err != nil {
		Error("db query error", query, err.Error())
		return nil
	}
	defer rows.Close()
	ret, err := processRows(rows)
	if err != nil {
		Error("db query error", query, err.Error())
		return nil
	}
	return ret
}

// 返回受影响行数
func processUpdateRet(query string, res sql.Result, err error) int64 {
	if err != nil {
		Error("db update error", query, err.Error())
		return -1
	}
	num, err := res.RowsAffected()
	if err != nil {
		Error("db update error", query, err.Error())
		return -1
	}
	return num
}

// 返回插入id
func processInsertRet(query string, res sql.Result, err error) int64 {
	if err != nil {
		Error("db exec error", query, err.Error())
		return -1
	}
	id, err := res.LastInsertId()
	if err != nil {
		Error("db exec error", query, err.Error())
		return -1
	}
	return id
}

/**
 * processRows 方法，将返回的rows封装为字典数组
 */
func processRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var list []map[string]interface{}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		var ret map[string]interface{}
		ret = make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				ret[columns[i]] = "null"
			} else {
				switch val := (*scanArgs[i].(*interface{})).(type) {
				case []byte:
					ret[columns[i]] = string(val)
					break
				default:
					ret[columns[i]] = val
				}
			}
		}
		list = append(list, ret)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}
