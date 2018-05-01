package mysql

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	METHOD_NOT_ALLOW = errors.New("kelp.mysql: this method is not allow")
	NO_DATA_TO_BIND  = errors.New("kelp.mysql: no data to bind")
)

// db is sql.DB connector implement Connector
type db struct {
	name string
	conn *sql.DB
}

// tx is sql.Tx connector implement Connector
type tx struct {
	name string
	conn *sql.Tx
}

// AddDB opens a mysql connection and store it into pool
func AddDB(name, dsn string, maxOpen, maxIdle int) error {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	if err := conn.Ping(); err != nil {
		return err
	}
	p.store[name] = &db{name: name, conn: conn}
	return nil
}

// Begin start an transaction
func (this *db) Begin() (Connector, error) {
	name := this.name + "-" + token()
	log.Debug(this.name, "begin", name)
	conn, err := this.conn.Begin()
	if err != nil {
		return nil, err
	}
	tx := &tx{name: name, conn: conn}
	return tx, nil
}

// Commit is not allow to query connector
func (this *db) Commit() error {
	log.Error(this.name, "commit", METHOD_NOT_ALLOW)
	return METHOD_NOT_ALLOW
}

// Rollback is not allow to qurey connector.
func (this *db) Rollback() error {
	log.Error(this.name, "rollback", METHOD_NOT_ALLOW)
	return METHOD_NOT_ALLOW
}

// Query select a set of data and bind into a dest list.
// The destList should be an pointor of slice assembled by data model.
func (this *db) Query(destList interface{}, sql string, params ...interface{}) error {
	log.Debug(this.name, "query", sql, params)
	stmt, err := this.conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return err
	}
	return scanQueryRows(destList, rows)
}

// QueryOne select one data and bind into a dest object.
// The destOjbect should be an pointer of data model.
func (this *db) QueryOne(destObject interface{}, sql string, params ...interface{}) error {
	log.Debug(this.name, "queryone", sql, params)
	stmt, err := this.conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return err
	}
	return scanQueryOne(destObject, rows)
}

// Insert executes an insert sql and returns last insert id.
func (this *db) Insert(sql string, params ...interface{}) (int64, error) {
	log.Debug(this.name, "insert", sql, params)
	ret, err := this.conn.Exec(sql, params...)
	if err != nil {
		return 0, err
	}
	return ret.LastInsertId()
}

// Execute executes a sql and returns effected rows.
func (this *db) Execute(sql string, params ...interface{}) (int64, error) {
	log.Debug(this.name, "execute", sql, params)
	ret, err := this.conn.Exec(sql, params...)
	if err != nil {
		return 0, err
	}
	return ret.RowsAffected()
}

// Begin is not allow to transaction connector.
func (this *tx) Begin() (Connector, error) {
	log.Error(this.name, "begin", METHOD_NOT_ALLOW)
	return nil, METHOD_NOT_ALLOW
}

// Commit commits a transaction
func (this *tx) Commit() error {
	log.Debug(this.name, "commit")
	return this.conn.Commit()
}

// Rollback rollback a transaction
func (this *tx) Rollback() error {
	log.Debug(this.name, "rollback")
	return this.conn.Rollback()
}

// Query select a set of data and bind into a dest list.
// The destList should be an pointor of slice assembled by data model.
func (this *tx) Query(destList interface{}, sql string, params ...interface{}) error {
	log.Debug(this.name, "query", sql, params)
	stmt, err := this.conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return err
	}
	return scanQueryRows(destList, rows)
}

// QueryOne select one data and bind into a dest object.
// The destOjbect should be an pointer of data model.
func (this *tx) QueryOne(destObject interface{}, sql string, params ...interface{}) error {
	log.Debug(this.name, "queryone", sql, params)
	stmt, err := this.conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return err
	}
	return scanQueryOne(destObject, rows)
}

// Insert executes an insert sql and returns last insert id.
func (this *tx) Insert(sql string, params ...interface{}) (lastInsertId int64, err error) {
	log.Debug(this.name, "insert", sql, params)
	ret, err := this.conn.Exec(sql, params...)
	if err != nil {
		return 0, err
	}
	return ret.LastInsertId()
}

// Execute executes a sql and returns effected rows.
func (this *tx) Execute(sql string, params ...interface{}) (int64, error) {
	log.Debug(this.name, "execute", sql, params)
	ret, err := this.conn.Exec(sql, params...)
	if err != nil {
		return 0, err
	}
	return ret.RowsAffected()
}

func scanQueryRows(dest interface{}, rows *sql.Rows) error {
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
				fieldName, ok := field.Tag.Lookup("column")
				if !ok {
					fieldName, ok = field.Tag.Lookup("json")
					if ok {
						fieldName = strings.Split(fieldName, ",")[0]
					}
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

func scanQueryOne(dest interface{}, rows *sql.Rows) error {
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
				fieldName, ok := field.Tag.Lookup("column")
				if !ok {
					fieldName, ok = field.Tag.Lookup("json")
					if ok {
						fieldName = strings.Split(fieldName, ",")[0]
					}
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
		return NO_DATA_TO_BIND
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// ToInt convert type to int
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
		return 0
	}
}

// ToInt64 convert type to int64
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
		return 0
	}
}

// ToFloat convert type to float64
func ToFloat(param interface{}) float64 {
	switch ret := param.(type) {
	case int64:
		return float64(ret)
	case float64:
		return ret
	case []byte:
		r, _ := strconv.ParseFloat(string(ret), 64)
		return r
	case string:
		r, _ := strconv.ParseFloat(ret, 64)
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
		return 0.0
	}
}

// ToBool convert type to bool
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
		}
		return false
	case string:
		switch ret {
		case "1", "true", "y", "on", "yes":
			return true
		case "0", "false", "n", "off", "no":
			return false
		default:
		}
		return false
	case nil:
		return false
	default:
		return false
	}
}

// ToString convert type to string
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
		return ""
	}
}

func ToTime(param interface{}) time.Time {
	switch ret := param.(type) {
	case []byte:
		r, err := time.ParseInLocation("2006-01-02 15:04:05", string(ret), time.Now().Location())
		if err != nil {
			return time.Now()
		}
		return r
	case string:
		r, err := time.ParseInLocation("2006-01-02 15:04:05", ret, time.Now().Location())
		if err != nil {
			return time.Now()
		}
		return r
	case time.Time:
		return ret
	default:
		return time.Now()
	}
}

func token() string {
	timestamp := []byte(strconv.FormatInt(time.Now().Unix(), 10))
	prefix := []byte(strconv.Itoa(rand.Intn(10000)))
	surfix := []byte(strconv.Itoa(rand.Intn(10000)))
	token := string(bytes.Join([][]byte{prefix, timestamp, surfix}, []byte("")))
	return token
}
