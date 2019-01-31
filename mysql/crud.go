package mysql

import (
	"errors"
	"reflect"
	"strings"
)

type CRUDWrapper struct {
	conn   Connector
	model  interface{}
	table  *Table
	option string
}

func CRUD(model interface{}, option string, conn Connector) *CRUDWrapper {
	return &CRUDWrapper{
		conn:  conn,
		model: model,
		table: NewTable(model),
	}
}

func (this *CRUDWrapper) Name() string {
	return this.table.modelType.Name()
}

func (this *CRUDWrapper) Elem() interface{} {
	return reflect.New(this.table.modelType).Interface()
}

func (this *CRUDWrapper) Create(data interface{}) error {
	if !strings.Contains("C", this.option) {
		return errors.New("do not support create")
	}
	if _, err := this.conn.Insert(
		this.table.getInsertSql(),
		this.table.bind(data)...,
	); err != nil {
		return err
	}
	return nil
}

func (this *CRUDWrapper) Update(id int64, data interface{}) error {
	if !strings.Contains("U", this.option) {
		return errors.New("do not support update")
	}
	if _, err := this.conn.Execute(
		this.table.getUpdateSql(),
		append(this.table.bind(data), id)...,
	); err != nil {
		return err
	}
	return nil
}

func (this *CRUDWrapper) Delete(id int64) error {
	if !strings.Contains("D", this.option) {
		return errors.New("do not support delete")
	}
	if _, err := this.conn.Execute(
		this.table.getDeleteSql(),
		id,
	); err != nil {
		return err
	}
	return nil
}

func (this *CRUDWrapper) Retrieve(limit, offset int64) (interface{}, int64, error) {
	if !strings.Contains("R", this.option) {
		return nil, 0, errors.New("do not support retrieve")
	}
	arrType := reflect.SliceOf(reflect.PtrTo(this.table.modelType))
	out := reflect.New(arrType).Interface()
	t := &struct {
		Total int64 `json:"total"`
	}{}
	if err := this.conn.QueryOne(
		t,
		this.table.getRetrieveCountSql(),
	); err != nil {
		return out, t.Total, err
	}
	if err := this.conn.Query(
		out,
		this.table.getRetrieveSql(limit, offset),
	); err != nil {
		return out, t.Total, err
	}
	return out, t.Total, nil
}
