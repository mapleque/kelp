package mysql

import (
	"strconv"
	"strings"
)

type WhereInterface interface {
	Add(field string, value interface{})
	Sql() string
	Bind() []interface{}
}

type InInterface interface {
	Add(field string, value interface{})
	Sql() string
	Bind() []interface{}
}

type SorterInterface interface {
	Add(SorterElement)
	Sql() string
}

type PagerInterface interface {
	Add(PagerElement)
	Sql() string
}

type SorterElement interface {
	GetField() string
	GetReverse() bool
}

type PagerElement interface {
	GetSize() int64
	GetOffset() int64
}

type Where struct {
	condition []string
	bind      []interface{}
}

func NewWhere() *Where {
	return &Where{
		condition: []string{},
		bind:      []interface{}{},
	}
}

func (this *Where) Add(field string, value interface{}) {
	switch value.(type) {
	case int,
		int8,
		int16,
		int32,
		int64,
		uint,
		uint8,
		uint16,
		uint32,
		uint64,
		float32,
		float64,
		bool:
		// continue
	case string:
		if value == "" {
			return
		}
	default:
		if value == nil {
			return
		}
	}
	this.condition = append(this.condition, field+" = ?")
	this.bind = append(this.bind, value)
}

func (this *Where) Sql() string {
	if len(this.condition) == 0 {
		return ""
	}
	return " WHERE " + strings.Join(this.condition, " AND ")
}

func (this *Where) Bind() []interface{} {
	return this.bind
}

type In struct {
	fields []string
	bind   map[string][]interface{}
}

func NewIn() *In {
	return &In{
		fields: []string{},
		bind:   map[string][]interface{}{},
	}
}

func (this *In) Add(field string, value interface{}) {
	if _, ok := this.bind[field]; !ok {
		this.bind[field] = []interface{}{}
		this.fields = append(this.fields, field)
	}
	this.bind[field] = append(this.bind[field], value)
}

func (this *In) Sql() string {
	sqls := []string{}
	if len(this.fields) == 0 {
		return ""
	}
	for _, field := range this.fields {
		values := this.bind[field]
		sqls = append(sqls,
			field+" IN ("+strings.TrimSuffix(strings.Repeat("?,", len(values)), ",")+")")
	}
	return " WHERE " + strings.Join(sqls, " AND ")
}

func (this *In) Bind() []interface{} {
	bind := []interface{}{}
	for _, field := range this.fields {
		bind = append(bind, this.bind[field]...)
	}
	return bind
}

type Sorter struct {
	sqls []string
}

func NewSorter() *Sorter {
	return &Sorter{
		sqls: []string{},
	}
}

func (this *Sorter) Add(sorter SorterElement) {
	sql := sorter.GetField()
	if sorter.GetReverse() {
		sql += " DESC"
	} else {
		sql += " AES"
	}
	this.sqls = append(this.sqls, sql)
}

func (this *Sorter) Sql() string {
	if len(this.sqls) == 0 {
		return ""
	}
	return " ORDER BY " + strings.Join(this.sqls, ", ")
}

type Pager struct {
	size   int64
	offset int64
}

func NewPager() *Pager {
	return &Pager{}
}

func (this *Pager) Add(pager PagerElement) {
	this.size = pager.GetSize()
	this.offset = pager.GetOffset()
}

func (this *Pager) Sql() string {
	if this.size <= 0 {
		return ""
	}
	return " LIMIT " + strconv.FormatInt(this.size, 10) +
		" OFFSET " + strconv.FormatInt(this.offset, 10)
}
