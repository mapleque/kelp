package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// _MysqlPool 类型， 是一个database容器，用于存储服务可能用到的所有db连接池
type _MysqlPool struct {
	pool map[string]*MysqlQuery
}

// MysqlQuery 对象，用于直接查询或执行
type MysqlQuery struct {
	database string
	conn     *sql.DB
}

// MysqlTransaction 对象，用于事物查询或执行
type MysqlTransaction struct {
	database string
	conn     *sql.Tx
}

// MysqlConnector 接口，用于屏蔽普通连接和事物连接的对象差异
type MysqlConnector interface {
	Select(sql string, param ...interface{}) ([]map[string]interface{}, error)
	Insert(sql string, param ...interface{}) (int64, error)
	Update(sql string, param ...interface{}) (int64, error)
	Execute(sql string, param ...interface{}) (int64, error)
	Query(list interface{}, sql string, params ...interface{}) error
	QueryOne(dest interface{}, sql string, params ...interface{}) error
}

// db全局私有变量，用于保持连接池
var mysqlPool *_MysqlPool

func init() {
	if mysqlPool != nil {
		return
	}
	log.Info("init mysql module...")
	mysqlPool = &_MysqlPool{}
	mysqlPool.pool = make(map[string]*MysqlQuery)
}

// AddDB 方法，添加一个database，并且在启动前验证其连通性
func AddMysql(name, dsn string, maxOpenConns, maxIdleConns int) {
	log.Info("add db", name, dsn)
	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Error("can not open db", dsn, err.Error())
		panic(err.Error()) // 直接panic，让server无法启动
	}
	if err := dbConn.Ping(); err != nil {
		log.Error("can not ping db", dsn, err.Error())
		panic(err.Error()) // 直接panic，让server无法启动
	}
	dbConn.SetMaxOpenConns(maxOpenConns)
	dbConn.SetMaxIdleConns(maxIdleConns)
	dbQuery := &MysqlQuery{}
	dbQuery.conn = dbConn
	dbQuery.database = name
	mysqlPool.pool[name] = dbQuery
}

// Use 方法，返回MysqlQuery对象
func UseMysql(database string) *MysqlQuery {
	return mysqlPool.pool[database]
}

// 获取MysqlQuery的connection
func (dbq *MysqlQuery) GetConn() *sql.DB {
	return dbq.conn
}

// Begin 方法，返回MysqlTransaction对象
func (dbq *MysqlQuery) Begin() (*MysqlTransaction, error) {
	log.Debug("[transaction begin]", "[db:"+dbq.database+"]")
	trans := &MysqlTransaction{}
	conn, err := dbq.conn.Begin()
	if err != nil {
		log.Error("db create transaction faild", dbq.database, err.Error())
		return trans, err
	}
	trans.conn = conn
	return trans, nil
}

func (dbq *MysqlQuery) Query(list interface{}, sql string, params ...interface{}) error {
	log.Debug("[query sql]", "[db:"+dbq.database+"]", sql, params)
	stmt, err := dbq.conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return err
	}
	return scanMysqlQueryRows(list, rows)
}

func (dbq *MysqlQuery) QueryOne(dest interface{}, sql string, params ...interface{}) error {
	log.Debug("[query sql]", "[db:"+dbq.database+"]", sql, params)
	stmt, err := dbq.conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return err
	}
	return scanMysqlQueryOne(dest, rows)
}

func (dbq *MysqlQuery) Query(obj []interface{}, sql string, params ...interface{}) error {
	log.Debug("[query sql]", "["+dbq.database+"]", sql, params)
	rows, err := dbq.conn.Query(sql, params...)
	if err != nil {
		return err
	}
	return scanMysqlQueryRows(obj, rows)
}

// Select 方法，返回查询结果数组
func (dbq *MysqlQuery) Select(
	sql string,
	params ...interface{}) ([]map[string]interface{}, error) {

	log.Debug("[select sql]", "[db:"+dbq.database+"]", sql, params)
	stmt, err := dbq.conn.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, err
	}
	ret, err := processMysqlRows(rows)
	if err != nil {
		log.Error("db query error", sql, params, err.Error())
		return nil, err
	}
	return ret, nil
}

// Update 方法，返回受影响行数
func (dbq *MysqlQuery) Update(
	sql string,
	params ...interface{}) (int64, error) {

	log.Debug("[update sql]", "[db:"+dbq.database+"]", sql, params)
	ret, err := dbq.conn.Exec(sql, params...)
	return processMysqlAffectedRet(sql, ret, err)
}

// Execute 方法，返回受影响行数
func (dbq *MysqlQuery) Execute(
	sql string,
	params ...interface{}) (int64, error) {

	log.Debug("[execute sql]", "[db:"+dbq.database+"]", sql, params)
	ret, err := dbq.conn.Exec(sql, params...)
	return processMysqlAffectedRet(sql, ret, err)
}

// Insert 方法，返回插入id
func (dbq *MysqlQuery) Insert(
	sql string,
	params ...interface{}) (int64, error) {

	log.Debug("[insert sql]", "[db:"+dbq.database+"]", sql, params)
	ret, err := dbq.conn.Exec(sql, params...)
	return processMysqlInsertRet(sql, ret, err)
}

func (dbt *MysqlTransaction) GetConn() *sql.Tx {
	return dbt.conn
}

func (dbt *MysqlTransaction) Query(list interface{}, sql string, params ...interface{}) error {
	log.Debug("[query sql in transaction]", "[db:"+dbt.database+"]", sql, params)
	stmt, err := dbt.conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return err
	}
	return scanMysqlQueryRows(list, rows)
}

func (dbt *MysqlTransaction) QueryOne(dest interface{}, sql string, params ...interface{}) error {
	log.Debug("[query sql in transaction]", "[db:"+dbt.database+"]", sql, params)
	stmt, err := dbt.conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return err
	}
	return scanMysqlQueryOne(dest, rows)
}

// Select 方法，返回查询结果数组
func (dbt *MysqlTransaction) Select(
	sql string,
	params ...interface{}) ([]map[string]interface{}, error) {

	log.Debug("[select sql in transaction]", "[db:"+dbt.database+"]", sql, params)
	stmt, err := dbt.conn.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, err
	}
	ret, err := processMysqlRows(rows)
	if err != nil {
		log.Error("db query error", sql, params, err.Error())
		return nil, err
	}
	return ret, nil
}

// Update 方法，返回受影响行数
func (dbt *MysqlTransaction) Update(
	sql string,
	params ...interface{}) (int64, error) {

	log.Debug("[update sql in transaction]", "[db:"+dbt.database+"]", sql, params)
	ret, err := dbt.conn.Exec(sql, params...)
	return processMysqlAffectedRet(sql, ret, err)
}

// Execute 方法，返回受影响行数
func (dbt *MysqlTransaction) Execute(
	sql string,
	params ...interface{}) (int64, error) {

	log.Debug("[execute sql in transaction]", "[db:"+dbt.database+"]", sql, params)
	ret, err := dbt.conn.Exec(sql, params...)
	return processMysqlAffectedRet(sql, ret, err)
}

// Insert 方法，返回插入id
func (dbt *MysqlTransaction) Insert(
	sql string,
	params ...interface{}) (int64, error) {

	log.Debug("[insert sql in transaction]", "[db:"+dbt.database+"]", sql, params)
	ret, err := dbt.conn.Exec(sql, params...)
	return processMysqlInsertRet(sql, ret, err)
}

// Commit 方法，提交事物
func (dbt *MysqlTransaction) Commit() error {
	log.Debug("[transaction commit]", "[db:"+dbt.database+"]")
	err := dbt.conn.Commit()
	if err != nil {
		log.Error("db transaction commit faild", dbt.database, err.Error())
	}
	return err
}

// Rollback 方法，回滚事物
func (dbt *MysqlTransaction) Rollback() error {
	log.Debug("[transaction rollback]", "[db:"+dbt.database+"]")
	err := dbt.conn.Rollback()
	if err != nil {
		log.Error("db transaction rollback faild", dbt.database, err.Error())
	}
	return err
}

// 返回受影响行数
func processMysqlAffectedRet(query string, res sql.Result, err error) (int64, error) {
	if err != nil {
		log.Error("db affected exec error", query, err.Error())
		return -1, err
	}
	num, err := res.RowsAffected()
	if err != nil {
		log.Error("db affected exec error", query, err.Error())
		return -1, err
	}
	return num, nil
}

// 返回插入id
func processMysqlInsertRet(query string, res sql.Result, err error) (int64, error) {
	if err != nil {
		log.Error("db insert exec error", query, err.Error())
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Error("db insert exec error", query, err.Error())
		return -1, err
	}
	return id, nil
}

// processMysqlRows 方法，将返回的rows封装为字典数组
func processMysqlRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	defer rows.Close()
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

func scanMysqlQueryOne(dest interface{}, rows *sql.Rows) error {
	defer rows.Close()
	// dest 必须是 ptr
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return fmt.Errorf("kelp.db.mysql: dest should be a ptr but %s", destType.Kind())
	}
	destValue := reflect.ValueOf(dest).Elem()
	if !destValue.CanSet() {
		return fmt.Errorf("kelp.db.mysql: dest can not set")
	}
	eleType := destType.Elem()
	// 必须要是struct类型
	if eleType.Kind() != reflect.Struct {
		return fmt.Errorf("kelp.db.mysql: target should be a *struct but *%s", eleType.Kind())
	}
	// 遍历查询结果
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	values := make([]interface{}, len(columnTypes))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// 新建一个元素实例
	ele := destValue
	if rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return err
		}
		for i, col := range values {
			key := columnTypes[i].Name()
			for j := 0; j < ele.NumField(); j++ {
				field := ele.Type().Field(j)
				if _, ignore := field.Tag.Lookup("ignore_db_bind"); ignore {
					continue
				}
				fieldName, ok := field.Tag.Lookup("column")
				if !ok {
					fieldName, ok = field.Tag.Lookup("json")
				}
				if !ok {
					fieldName = strings.ToLower(field.Name)
				}
				if key == fieldName {
					eleField := ele.FieldByName(field.Name)
					if eleField.CanSet() {
						switch field.Type.Kind() {
						case reflect.Int:
							eleField.Set(reflect.ValueOf(ToInt(col)))
						case reflect.Float64:
							eleField.Set(reflect.ValueOf(ToFloat(col)))
						case reflect.String:
							eleField.Set(reflect.ValueOf(ToString(col)))
						case reflect.Bool:
							eleField.Set(reflect.ValueOf(ToBool(col)))
						case reflect.Struct:
							switch {
							case field.Type.Name() == "Time":
								eleField.Set(reflect.ValueOf(ToTime(col)))
							}
						default:
							eleField.Set(reflect.ValueOf(col))
						}
					}
				}
			}
		}
	} else {
		return fmt.Errorf("kelp.db.mysql: no data to bind")
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func scanMysqlQueryRows(dest interface{}, rows *sql.Rows) error {
	// dest 必须是 ptr
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return fmt.Errorf("kelp.db.mysql: dest should be a ptr but %s", destType.Kind())
	}
	destValue := reflect.ValueOf(dest).Elem()
	if !destValue.CanSet() {
		return fmt.Errorf("kelp.db.mysql: dest can not set")
	}
	listType := destType.Elem()

	// list必须是slice
	if listType.Kind() != reflect.Slice {
		return fmt.Errorf("kelp.db.mysql: target should be a slice but %s", listType.Kind())
	}
	// 获取list的元素类型
	eleType := listType.Elem()
	isPointer := false
	// 如果是指针类型，就再取真实类型
	if eleType.Kind() == reflect.Ptr {
		eleType = eleType.Elem()
		isPointer = true
	}

	// 必须要是struct类型
	if eleType.Kind() != reflect.Struct {
		return fmt.Errorf("kelp.db.mysql: target should be a []struct or a []*struct but []%s", eleType.Kind())
	}

	// 遍历查询结果
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	values := make([]interface{}, len(columnTypes))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		// 根据field name和tag给元素赋值
		err = rows.Scan(scanArgs...)
		if err != nil {
			return err
		}
		// 新建一个元素实例
		eleAddr := reflect.New(eleType)
		ele := eleAddr.Elem()
		for i, col := range values {
			key := columnTypes[i].Name()
			for j := 0; j < ele.NumField(); j++ {
				field := ele.Type().Field(j)
				if _, ignore := field.Tag.Lookup("ignore_db_bind"); ignore {
					continue
				}
				fieldName, ok := field.Tag.Lookup("column")
				if !ok {
					fieldName, ok = field.Tag.Lookup("json")
				}
				if !ok {
					fieldName = strings.ToLower(field.Name)
				}
				if key == fieldName {
					eleField := ele.FieldByName(field.Name)
					if eleField.CanSet() {
						switch field.Type.Kind() {
						case reflect.Int:
							eleField.Set(reflect.ValueOf(ToInt(col)))
						case reflect.Int64:
							eleField.Set(reflect.ValueOf(ToInt64(col)))
						case reflect.Float64:
							eleField.Set(reflect.ValueOf(ToFloat(col)))
						case reflect.String:
							eleField.Set(reflect.ValueOf(ToString(col)))
						case reflect.Bool:
							eleField.Set(reflect.ValueOf(ToBool(col)))
						case reflect.Struct:
							switch {
							case field.Type.Name() == "Time":
								eleField.Set(reflect.ValueOf(ToTime(col)))
							}
						default:
							eleField.Set(reflect.ValueOf(col))
						}
					}
				}
			}
		}

		if isPointer {
			// 元素是指针，要往slice里append指针
			destValue.Set(reflect.Append(destValue, ele.Addr()))
		} else {
			destValue.Set(reflect.Append(destValue, ele))
		}
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

// 类型转换，任何类型转成int
func ToInt(param interface{}) int {
	switch ret := param.(type) {
	case int:
		return ret
	case int64:
		return int(ret)
	case float64:
		return int(ret)
	case []byte:
		r, _ := strconv.Atoi(string(ret))
		return r
	case string:
		r, _ := strconv.Atoi(ret)
		return r
	case bool:
		if ret {
			return 1
		} else {
			return 0
		}
	case nil:
		return 0
	default:
		log.Error("param type change to int error",
			ret, fmt.Sprintf("%T", ret))
		return 0
	}
}

// 类型转换，任何类型转成int64
func ToInt64(param interface{}) int64 {
	switch ret := param.(type) {
	case int:
		return int64(ret)
	case int64:
		return ret
	case float64:
		return int64(ret)
	case []byte:
		r, _ := strconv.ParseInt(string(ret), 10, 64)
		return r
	case string:
		r, _ := strconv.ParseInt(ret, 10, 64)
		return r
	case bool:
		if ret {
			return 1
		} else {
			return 0
		}
	case nil:
		return 0
	default:
		log.Error("param type change to int error",
			ret, fmt.Sprintf("%T", ret))
		return 0
	}
}

// 类型转换，类型转换成float
func ToFloat(param interface{}) float64 {
	switch ret := param.(type) {
	case int64:
		return float64(ret)
	case float64:
		return ret
	case []byte:
		r, err := strconv.ParseFloat(string(ret), 64)
		if err != nil {
			log.Error("param type change error", ret, err.Error())
		}
		return r
	case string:
		r, err := strconv.ParseFloat(ret, 64)
		if err != nil {
			log.Error("param type change error", ret, err.Error())
		}
		return r
	case bool:
		if ret {
			return 1.0
		} else {
			return 0.0
		}
	case nil:
		return 0.0
	default:
		log.Error("param type change to int error",
			ret, fmt.Sprintf("%T", ret))
		return 0.0
	}
}

// 类型转换，任何类型转成bool
func ToBool(param interface{}) bool {
	switch ret := param.(type) {
	case bool:
		return ret
	case int64:
		if ret > 0 {
			return true
		} else {
			return false
		}
	case float64:
		if ret > 0 {
			return true
		} else {
			return false
		}
	case []byte:
		switch string(ret) {
		case "1", "true", "y", "on", "yes":
			return true
		case "0", "false", "n", "off", "no":
			return false
		default:
			log.Error("param type change to bool error", ret, "unknown type")
		}
		return false
	case string:
		switch ret {
		case "1", "true", "y", "on", "yes":
			return true
		case "0", "false", "n", "off", "no":
			return false
		default:
			log.Error("param type change to bool error", ret, "unknown type")
		}
		return false
	case nil:
		return false
	default:
		log.Error("param type change to bool error",
			ret, fmt.Sprintf("%T", ret))
		return false
	}
}

// 类型转换，任何类型转成string
func ToString(param interface{}) string {
	switch ret := param.(type) {
	case string:
		return ret
	case []byte:
		return string(ret)
	case int64:
		return strconv.FormatInt(ret, 10)
	case float64:
		return strconv.FormatFloat(ret, 'f', -1, 64)
	case bool:
		if ret {
			return "1"
		} else {
			return "0"
		}
	case time.Time:
		return fmt.Sprint(ret)
	case nil:
		return ""
	default:
		log.Error("param type change to string error",
			ret, fmt.Sprintf("%T", ret))
		return ""
	}
}

func ToTime(param interface{}) time.Time {
	switch ret := param.(type) {
	case []byte:
		r, err := time.ParseInLocation("2006-01-02 15:04:05", string(ret), time.Now().Location())
		if err != nil {
			log.Error("param type change to time error", string(ret), err)
			return time.Now()
		}
		return r
	case string:
		r, err := time.ParseInLocation("2006-01-02 15:04:05", ret, time.Now().Location())
		if err != nil {
			log.Error("param type change to time error", ret, err)
			return time.Now()
		}
		return r
	case time.Time:
		return ret
	default:
		log.Error("param type change to time error", ret, fmt.Sprintf("%T", ret))
		return time.Now()
	}
}
