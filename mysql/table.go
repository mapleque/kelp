package mysql

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Table struct {
	model     interface{}
	modelType reflect.Type

	name       string
	charset    string
	additional string
	fields     []*TableField
}

type TableField struct {
	propName  string
	propIndex int

	name   string
	schema string
}

func NewTable(model interface{}) *Table {
	table := &Table{
		model:   model,
		charset: "utf8",
		fields:  []*TableField{},
	}

	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}
	if modelValue.Type().Kind() != reflect.Struct {
		panic("model must be a struct or a struct ptr")
	}
	table.modelType = modelValue.Type()
	table.name = getTableNameFromTypeName(modelValue.Type().Name())

	for i := 0; i < modelValue.NumField(); i++ {
		fieldType := modelValue.Type().Field(i)
		columnName := getColumnName(fieldType)
		if columnName == "-" {
			continue
		}
		fieldSchema, exist := fieldType.Tag.Lookup("column_schema")
		if !exist {
			panic("column field has no schema tag")
		}
		table.fields = append(table.fields, &TableField{
			propName: fieldType.Name,
			name:     columnName,
			schema:   fieldSchema,
		})
	}
	return table
}

func (this *Table) SetCharset(charset string) *Table {
	this.charset = charset
	return this
}

func (this *Table) SetName(name string) *Table {
	this.name = name
	return this
}

func (this *Table) Append(additional string) *Table {
	this.additional += additional
	return this
}

func (this *Table) getSqlFilename() string {
	return this.name + ".sql"
}

func (this *Table) getCreateTableSql() string {
	fields := this.getFieldsSql()
	sql := ""
	sql += fmt.Sprintf("DROP TABLE IF EXISTS `%s`;\n", this.name)
	sql += fmt.Sprintf("CREATE TABLE `%s` (\n", this.name)
	sql += strings.Join(fields, ",\n")
	sql += fmt.Sprintf("\n) DEFAULT CHARSET=%s;\n", this.charset)
	sql += fmt.Sprintf("%s\n", this.additional)
	return sql
}

func (this *Table) getFieldsSql() []string {
	fields := []string{}
	for _, field := range this.fields {
		fields = append(fields, fmt.Sprintf("\t`%s` %s", field.name, field.schema))
	}

	return fields
}

func (this *Table) bind(data interface{}) []interface{} {
	ret := []interface{}{}
	v := reflect.ValueOf(data).Elem()
	for _, field := range this.fields {
		if field.name == "id" {
			continue
		}
		ret = append(ret, v.FieldByName(field.propName).Interface())
	}
	return ret
}

func (this *Table) getInsertSql() string {
	fieldNames := []string{}
	fieldParams := []string{}
	for _, field := range this.fields {
		if field.name == "id" {
			continue
		}
		fieldNames = append(fieldNames, "`"+field.name+"`")
		fieldParams = append(fieldParams, "?")
	}
	return fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES (%s)",
		this.name,
		strings.Join(fieldNames, ","),
		strings.Join(fieldParams, ","),
	)
}

func (this *Table) getUpdateSql() string {
	fields := []string{}
	for _, field := range this.fields {
		if field.name == "id" {
			continue
		}
		fields = append(fields, fmt.Sprintf("`%s`=?", field.name))
	}
	return fmt.Sprintf(
		"UPDATE `%s` SET %s WHERE id = ?",
		this.name,
		strings.Join(fields, ","),
	)
}

func (this *Table) getDeleteSql() string {
	return fmt.Sprintf(
		"DELETE FROM `%s` WHERE id = ? LIMIT 1",
		this.name,
	)
}

func (this *Table) getRetrieveCountSql() string {
	return fmt.Sprintf(
		"SELECT COUNT(*) AS total FROM `%s`",
		this.name,
	)
}

func (this *Table) getRetrieveSql(limit, offset int64) string {
	return fmt.Sprintf(
		"SELECT * FROM `%s` ORDER BY id DESC LIMIT %d OFFSET %d",
		this.name,
		limit,
		offset,
	)
}

func getColumnName(fieldType reflect.StructField) string {
	if tag, exist := fieldType.Tag.Lookup("column"); exist && len(tag) > 0 {
		return tag
	}
	if tag, exist := fieldType.Tag.Lookup("json"); exist && len(tag) > 0 {
		return strings.Split(tag, ",")[0]
	}
	return fieldType.Name
}

func getTableNameFromTypeName(typename string) string {
	r, _ := regexp.Compile("[A-Z]([a-z\\d]+)?")
	ng := []string{}
	for _, g := range r.FindAllString(typename, -1) {
		ng = append(ng, strings.ToLower(g))
	}
	return strings.Join(ng, "_")
}
