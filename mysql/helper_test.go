package mysql

import (
	"testing"
)

func TestWhere(t *testing.T) {
	where := NewWhere()
	where.Add("field1", "value1")
	if where.Sql() != " WHERE field1 = ?" {
		t.Error("wrong where sql", where.Sql())
	}
	if len(where.Bind()) != 1 || where.Bind()[0] != "value1" {
		t.Error("wrong where bind", where.Bind())
	}
	where = NewWhere()
	where.Add("field1", "value1")
	where.Add("field2", "value2")
	if where.Sql() != " WHERE field1 = ? AND field2 = ?" {
		t.Error("wrong where sql", where.Sql())
	}
	if len(where.Bind()) != 2 || where.Bind()[0] != "value1" || where.Bind()[1] != "value2" {
		t.Error("wrong where bind", where.Bind())
	}
}

func TestIn(t *testing.T) {
	in := NewIn()
	in.Add("field1", "value1")
	in.Add("field1", "value2")
	if in.Sql() != " WHERE field1 IN (?,?)" {
		t.Error("wrong in sql", in.Sql())
	}
	if len(in.Bind()) != 2 || in.Bind()[0] != "value1" || in.Bind()[1] != "value2" {
		t.Error("wrong in bind", in.Bind())
	}
	in.Add("field2", "value3")
	in.Add("field2", "value4")
	if in.Sql() != " WHERE field1 IN (?,?) AND field2 IN (?,?)" {
		t.Error("wrong in sql", in.Sql())
	}
	if len(in.Bind()) != 4 || in.Bind()[2] != "value3" || in.Bind()[3] != "value4" {
		t.Error("wrong in bind", in.Bind())
	}

}

type SorterForTest struct {
	field   string
	reverse bool
}

func (this *SorterForTest) GetField() string {
	return this.field
}
func (this *SorterForTest) GetReverse() bool {
	return this.reverse
}

func TestSorter(t *testing.T) {
	sorter := NewSorter()
	sorter.Add(&SorterForTest{"field1", true})
	if sorter.Sql() != " ORDER BY field1 DESC" {
		t.Error("wrong sorter sql", sorter.Sql())
	}

	sorter = NewSorter()
	sorter.Add(&SorterForTest{"field1", true})
	sorter.Add(&SorterForTest{"field2", false})
	if sorter.Sql() != " ORDER BY field1 DESC, field2 AES" {
		t.Error("wrong sorter sql", sorter.Sql())
	}
}

type PagerForTest struct {
	size   int64
	offset int64
}

func (this *PagerForTest) GetSize() int64 {
	return this.size
}

func (this *PagerForTest) GetOffset() int64 {
	return this.offset
}

func TestPager(t *testing.T) {
	pager := NewPager()
	pager.Add(&PagerForTest{10, 1})
	if pager.Sql() != " LIMIT 10 OFFSET 1" {
		t.Error("wrong pager sql", pager.Sql())
	}
	pager.Add(&PagerForTest{0, 1})
	if pager.Sql() != "" {
		t.Error("wrong pager sql", pager.Sql())
	}
}
