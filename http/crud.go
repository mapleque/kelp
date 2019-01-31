package http

import (
	"reflect"
)

type CRUDHandlerer interface {
	Name() string
	Elem() interface{}
	Create(data interface{}) error
	Update(id int64, data interface{}) error
	Delete(id int64) error
	Retrieve(limit, offset int64) (interface{}, int64, error)
}

func (this *Router) CRUD(h CRUDHandlerer) {
	this.Handle("Create "+h.Name(), "/create", func() interface{} {
		handler := reflect.FuncOf(
			[]reflect.Type{reflect.TypeOf(h.Elem())},
			[]reflect.Type{reflect.TypeOf(&Status{})},
			false,
		)

		return reflect.MakeFunc(handler, func(args []reflect.Value) []reflect.Value {
			if err := h.Create(args[0].Interface()); err != nil {
				Error(err)
				return []reflect.Value{reflect.ValueOf(STATUS_ERROR_DB)}
			}
			var nilret *Status
			return []reflect.Value{reflect.ValueOf(nilret)}
		}).Interface()
	}())
	this.Handle("Update "+h.Name(), "/update", func() interface{} {
		inStruct := reflect.StructOf([]reflect.StructField{
			{
				Name: "Id",
				Type: reflect.TypeOf(int64(0)),
				Tag:  `json:"id" valid:"(0,),message=invalid id"`,
			},
			{
				Name: "Data",
				Type: reflect.TypeOf(h.Elem()),
				Tag:  `json:"data" valid:"message=invalid data"`,
			},
		})
		handler := reflect.FuncOf(
			[]reflect.Type{reflect.PtrTo(inStruct)},
			[]reflect.Type{reflect.TypeOf(&Status{})},
			false,
		)
		return reflect.MakeFunc(handler, func(args []reflect.Value) []reflect.Value {
			if err := h.Update(
				args[0].Elem().FieldByName("Id").Int(),
				args[0].Elem().FieldByName("Data").Interface(),
			); err != nil {
				return []reflect.Value{reflect.ValueOf(STATUS_ERROR_DB)}
			}
			var nilret *Status
			return []reflect.Value{reflect.ValueOf(nilret)}
		}).Interface()
	}())
	this.Handle("Delete "+h.Name(), "/delete", func(in *struct {
		Id int64 `json:"id" valid:"(0,),message=invalid id"`
	}) *Status {
		if err := h.Delete(in.Id); err != nil {
			return STATUS_ERROR_DB
		}
		return nil
	})
	this.Handle("Retrieve "+h.Name(), "/retrieve", func() interface{} {
		var in *struct {
			Size   int64 `json:"size" valid:"message=invalid size"`
			Offset int64 `json:"offset" valid:"message=invalid offset"`
		}
		outStruct := reflect.StructOf([]reflect.StructField{
			{
				Name: "Total",
				Type: reflect.TypeOf(int64(0)),
				Tag:  `json:"total" comment:"数据总量，用于分页，根据约定返回"`,
			},
			{
				Name: "List",
				Type: reflect.SliceOf(reflect.TypeOf(h.Elem())),
				Tag:  `json:"list" comment:"数据列表"`,
			},
		})
		handler := reflect.FuncOf(
			[]reflect.Type{reflect.TypeOf(in), reflect.PtrTo(outStruct)},
			[]reflect.Type{reflect.TypeOf(&Status{})},
			false,
		)
		return reflect.MakeFunc(handler, func(args []reflect.Value) []reflect.Value {
			if data, lastIdOrOffset, err := h.Retrieve(
				args[0].Elem().FieldByName("Size").Int(),
				args[0].Elem().FieldByName("Offset").Int(),
			); err != nil {
				return []reflect.Value{reflect.ValueOf(STATUS_ERROR_DB)}
			} else {
				args[1].Elem().FieldByName("Total").Set(reflect.ValueOf(lastIdOrOffset))
				args[1].Elem().FieldByName("List").Set(reflect.ValueOf(data).Elem())
			}
			var nilret *Status
			return []reflect.Value{reflect.ValueOf(nilret)}
		}).Interface()
	}())
}
