package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// DBPool 类型， 是一个database容器，用于存储服务可能用到的所有db连接池
type DBPool struct {
	pool     map[string]*DBQuery
	openFunc func(string) Connector
}

type Connector interface {
	Begin() (*sql.Tx, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	Exec(string, ...interface{}) (sql.Result, error)
	Ping() error
	Close() error
	SetMaxOpenConns(int)
	SetMaxIdleConns(int)
}

type TransConnector interface {
	Commit() error
	Rollback() error
	Query(string, ...interface{}) (*sql.Rows, error)
	Exec(string, ...interface{}) (sql.Result, error)
}

// DBQuery 对象，用于直接查询或执行
type DBQuery struct {
	database string
	conn     Connector
}

// DBTransaction 对象，用于事物查询或执行
type DBTransaction struct {
	database string
	conn     TransConnector
}

// db全局私有变量，用于保持连接池
var db *DBPool

func init() {
	if db != nil {
		return
	}
	log.Info("init db module...")
	db = &DBPool{}
	db.pool = make(map[string]*DBQuery)
	db.openFunc = func(dsn string) Connector {
		dbConn, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Error("can not open db", dsn, err.Error())
			panic(err.Error())
		}
		return dbConn
	}
}

// 可以重定向打开db的方法
func SetOpenFunc(openFunc func(dsn string) Connector) {
	db.openFunc = openFunc
}

/**
 * AddDB 方法，添加一个database，并且在启动前验证其连通性
 */
func AddDB(name, dsn string, maxOpenConns, maxIdleConns int) {
	log.Info("add db", name, dsn)
	dbConn := db.openFunc(dsn)
	err := dbConn.Ping()
	if err != nil {
		log.Error("can not ping db", dsn, err.Error())
		panic(err.Error()) // 直接panic，让server无法启动
	}
	dbConn.SetMaxOpenConns(maxOpenConns)
	dbConn.SetMaxIdleConns(maxIdleConns)
	dbQuery := &DBQuery{}
	dbQuery.conn = dbConn
	dbQuery.database = name
	db.pool[name] = dbQuery
}

// User 方法，返回DBQuery对象
func Use(database string) *DBQuery {
	return db.pool[database]
}

// Begin 方法，返回DBTransaction对象
func Begin(database string) *DBTransaction {
	return db.pool[database].Begin()
}

// Select 方法，返回查询结果数组
func Select(
	database, sql string,
	params ...interface{}) []map[string]interface{} {

	return db.pool[database].Select(sql, params...)
}

// Update 方法，返回受影响行数
func Update(
	database, sql string,
	params ...interface{}) int64 {

	return db.pool[database].Update(sql, params...)
}

// Execute 方法，返回受影响行数
func Execute(
	database, sql string,
	params ...interface{}) int64 {

	return db.pool[database].Execute(sql, params...)
}

// Insert 方法，返回插入id
func Insert(
	database, sql string,
	params ...interface{}) int64 {

	return db.pool[database].Insert(sql, params...)
}

// 获取DBQuery的connection
func (dbq *DBQuery) GetConn() Connector {
	return dbq.conn
}

// Begin 方法，返回DBTransaction对象
func (dbq *DBQuery) Begin() *DBTransaction {
	log.Debug("[transaction begin]", "["+dbq.database+"]")
	trans := &DBTransaction{}
	conn, err := dbq.conn.Begin()
	if err != nil {
		log.Error("db create transaction faild", dbq.database, err.Error())
		return trans
	}
	trans.conn = conn
	return trans
}

// Select 方法，返回查询结果数组
func (dbq *DBQuery) Select(
	sql string,
	params ...interface{}) []map[string]interface{} {

	log.Debug("[select sql]", "["+dbq.database+"]", sql, params)
	ret, err := dbq.conn.Query(sql, params...)
	return processQueryRet(sql, ret, err)
}

// Update 方法，返回受影响行数
func (dbq *DBQuery) Update(
	sql string,
	params ...interface{}) int64 {

	log.Debug("[update sql]", "["+dbq.database+"]", sql, params)
	ret, err := dbq.conn.Exec(sql, params...)
	return processAffectedRet(sql, ret, err)
}

// Execute 方法，返回受影响行数
func (dbq *DBQuery) Execute(
	sql string,
	params ...interface{}) int64 {

	log.Debug("[execute sql]", "["+dbq.database+"]", sql, params)
	ret, err := dbq.conn.Exec(sql, params...)
	return processAffectedRet(sql, ret, err)
}

// Insert 方法，返回插入id
func (dbq *DBQuery) Insert(
	sql string,
	params ...interface{}) int64 {

	log.Debug("[insert sql]", "["+dbq.database+"]", sql, params)
	ret, err := dbq.conn.Exec(sql, params...)
	return processInsertRet(sql, ret, err)
}

// 获取transaction的connection
func (dbt *DBTransaction) GetConn() TransConnector {
	return dbt.conn
}

// Select 方法，返回查询结果数组
func (dbt *DBTransaction) Select(
	sql string,
	params ...interface{}) []map[string]interface{} {

	log.Debug("[select sql in transaction]", "["+dbt.database+"]", sql, params)
	ret, err := dbt.conn.Query(sql, params...)
	return processQueryRet(sql, ret, err)
}

// Update 方法，返回受影响行数
func (dbt *DBTransaction) Update(
	sql string,
	params ...interface{}) int64 {

	log.Debug("[update sql in transaction]", "["+dbt.database+"]", sql, params)
	ret, err := dbt.conn.Exec(sql, params...)
	return processAffectedRet(sql, ret, err)
}

// Execute 方法，返回受影响行数
func (dbt *DBTransaction) Execute(
	sql string,
	params ...interface{}) int64 {

	log.Debug("[execute sql in transaction]", "["+dbt.database+"]", sql, params)
	ret, err := dbt.conn.Exec(sql, params...)
	return processAffectedRet(sql, ret, err)
}

// Insert 方法，返回插入id
func (dbt *DBTransaction) Insert(
	sql string,
	params ...interface{}) int64 {

	log.Debug("[insert sql in transaction]", "["+dbt.database+"]", sql, params)
	ret, err := dbt.conn.Exec(sql, params...)
	return processInsertRet(sql, ret, err)
}

// Commit 方法，提交事物
func (dbt *DBTransaction) Commit() {
	log.Debug("[transaction commit]", "["+dbt.database+"]")
	err := dbt.conn.Commit()
	if err != nil {
		log.Error("db transaction commit faild", dbt.database, err.Error())
	}
}

// Rollback 方法，回滚事物
func (dbt *DBTransaction) Rollback() {
	log.Debug("[transaction rollback]", "["+dbt.database+"]")
	err := dbt.conn.Rollback()
	if err != nil {
		log.Error("db transaction rollback faild", dbt.database, err.Error())
	}
}

// 返回查询结果数组
func processQueryRet(
	query string, rows *sql.Rows, err error) []map[string]interface{} {
	if err != nil {
		log.Error("db query error", query, err.Error())
		return nil
	}
	defer rows.Close()
	ret, err := processRows(rows)
	if err != nil {
		log.Error("db query error", query, err.Error())
		return nil
	}
	return ret
}

// 返回受影响行数
func processAffectedRet(query string, res sql.Result, err error) int64 {
	if err != nil {
		log.Error("db affected exec error", query, err.Error())
		return -1
	}
	num, err := res.RowsAffected()
	if err != nil {
		log.Error("db affected exec error", query, err.Error())
		return -1
	}
	return num
}

// 返回插入id
func processInsertRet(query string, res sql.Result, err error) int64 {
	if err != nil {
		log.Error("db insert exec error", query, err.Error())
		return -1
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Error("db insert exec error", query, err.Error())
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
	list := []map[string]interface{}{}
	// 这里需要初始化为空数组，否则在查询结果为空的时候，返回的会是一个未初始化的指针
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		var ret map[string]interface{}
		ret = make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				ret[columns[i]] = nil
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
