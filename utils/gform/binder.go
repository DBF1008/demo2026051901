package gform

import (
	"errors"
	"fmt"
	"github.com/gohouse/t"
	"reflect"
)

// Map ...
type Map t.MapStringT

// Data ...
type Data map[string]interface{}

// BindType ...
type BindType int

const (
	OBJECT_STRUCT BindType = iota
	OBJECT_STRUCT_SLICE
	OBJECT_MAP
	OBJECT_MAP_SLICE
	OBJECT_STRING
	OBJECT_MAP_T
	OBJECT_MAP_SLICE_T
	OBJECT_NIL
)

// BindString ...
var BindString = map[BindType]string{
	OBJECT_STRUCT:       "OBJECT_STRUCT",
	OBJECT_STRUCT_SLICE: "OBJECT_STRUCT_SLICE",
	OBJECT_MAP:          "OBJECT_MAP",
	OBJECT_MAP_SLICE:    "OBJECT_MAP_SLICE",
	OBJECT_STRING:       "OBJECT_STRING",
	OBJECT_MAP_T:        "OBJECT_MAP_T",
	OBJECT_MAP_SLICE_T:  "OBJECT_MAP_SLICE_T",
	OBJECT_NIL:          "OBJECT_NIL",
}

// BindType.String ...
func (b BindType) String() string {
	return BindString[b]
}

// Binder ...
type Binder struct {
	BindOrigin interface{}
	//BindOriginTableName []string
	BindName string
	BindResult interface{}
	BindResultSlice reflect.Value
	BindFields []string
	BindType BindType
	BindLimit  int
	BindPrefix string
	BindAll []Data
}

var _ IBinder = &Binder{}

// NewBinder ...
func NewBinder(o ...interface{}) *Binder {
	var binder = new(Binder)
	if len(o) > 0 {
		binder.SetBindOrigin(o[0])
	} else {
		binder.BindType = OBJECT_NIL
	}
	return binder
}

// BindParse ...
func (o *Binder) BindParse(prefix string) error {
	if o.GetBindOrigin() == nil {
		return nil
	}
	var BindName string
	switch o.GetBindOrigin().(type) {
	case string:
		o.SetBindType(OBJECT_STRING)
		BindName = o.GetBindOrigin().(string)
		//o.SetBindAll([]Map{})

	default:
		o.SetBindFields([]string{})
		// make sure dst is an appropriate type
		dstVal := reflect.ValueOf(o.GetBindOrigin())
		sliceVal := reflect.Indirect(dstVal)

		switch sliceVal.Kind() {
		case reflect.Struct: // struct
			o.SetBindType(OBJECT_STRUCT)
			BindName = sliceVal.Type().Name()
			o.SetBindResult(o.GetBindOrigin())
			//o.SetBindLimit(1)
			o.parseFields()
			switch dstVal.Kind() {
			case reflect.Ptr, reflect.Struct:
				if tn := dstVal.MethodByName("TableName"); tn.IsValid() {
					BindName = tn.Call(nil)[0].String()
				}
			default:
				return errors.New("传入的对象有误,示例:var user User,传入 &user{}")
			}
		case reflect.Map: // map
			o.SetBindType(OBJECT_MAP)
			//o.SetBindLimit(1)
			o.SetBindResult(o.GetBindOrigin())
			if sliceVal.Type().Elem() == reflect.ValueOf(map[string]t.Type{}).Type().Elem() {
				o.SetBindType(OBJECT_MAP_T)
			}
			if dstVal.Kind() != reflect.Ptr {
				return errors.New("传入的不是map指针,如:var user gform.Map,传入 &user{}")
			}
			if tn := dstVal.MethodByName("TableName"); tn.IsValid() {
				BindName = tn.Call(nil)[0].String()
			}

		case reflect.Slice: // []struct,[]map
			eltType := sliceVal.Type().Elem()

			switch eltType.Kind() {
			case reflect.Map:
				o.SetBindType(OBJECT_MAP_SLICE)
				o.SetBindResult(reflect.MakeMap(eltType).Interface())
				o.SetBindResultSlice(sliceVal)
				if eltType.Elem() == reflect.ValueOf(map[string]t.Type{}).Type().Elem() {
					o.SetBindType(OBJECT_MAP_SLICE_T)
				}
				if dstVal.Kind() != reflect.Ptr {
					return errors.New("传入的不是map指针,如:var user gform.Map,传入 &user{}")
				}
				r2val := reflect.New(eltType)
				if tn := r2val.MethodByName("TableName"); tn.IsValid() {
					BindName = tn.Call(nil)[0].String()
				}

			case reflect.Struct:
				o.SetBindType(OBJECT_STRUCT_SLICE)
				BindName = eltType.Name()
				br := reflect.New(eltType)
				o.SetBindResult(br.Interface())
				o.SetBindResultSlice(sliceVal)
				o.parseFields()

				switch dstVal.Kind() {
				case reflect.Ptr, reflect.Struct:
					if tn := br.MethodByName("TableName"); tn.IsValid() {
						BindName = tn.Call(nil)[0].String()
					}
				default:
					return errors.New("传入的对象有误,示例:var user User,传入 &user{}")
				}
			default:
				return fmt.Errorf("table只接收 struct,[]struct,map[string]interface{},[]map[string]interface{}的对象和地址, 但是传入的是: %T", o.GetBindOrigin())
			}
			if tn := dstVal.MethodByName("TableName"); tn.IsValid() {
				BindName = tn.Call(nil)[0].String()
			}
		default:
			return fmt.Errorf("table只接收 struct,[]struct,map[string]interface{},[]map[string]interface{}, 但是传入的是: %T", o.GetBindOrigin())
		}
	}

	o.SetBindName(prefix + BindName)
	o.SetBindPrefix(prefix)
	return nil
}

func (o *Binder) parseFields() {
	if len(o.GetBindFields()) == 0 {
		o.SetBindFields(getTagName(o.GetBindResult(), TAGNAME))
	}
}

// ResetBindResultSlice ...
func (o *Binder) ResetBindResultSlice() {
	if o.BindType == OBJECT_MAP_SLICE_T {
		o.BindResultSlice = reflect.New(o.BindResultSlice.Type())
	}
}

// SetBindPrefix ...
func (o *Binder) SetBindPrefix(arg string) {
	o.BindPrefix = arg
}

// GetBindPrefix ...
func (o *Binder) GetBindPrefix() string {
	return o.BindPrefix
}

// SetBindOrigin ...
func (o *Binder) SetBindOrigin(arg interface{}) {
	o.BindOrigin = arg
}

// GetBindOrigin ...
func (o *Binder) GetBindOrigin() interface{} {
	return o.BindOrigin
}

// SetBindName ...
func (o *Binder) SetBindName(arg string) {
	o.BindName = arg
}

// GetBindName ...
func (o *Binder) GetBindName() string {
	return o.BindName
}

// SetBindResult ...
func (o *Binder) SetBindResult(arg interface{}) {
	o.BindResult = arg
}

// GetBindResult ...
func (o *Binder) GetBindResult() interface{} {
	return o.BindResult
}

// SetBindResultSlice ...
func (o *Binder) SetBindResultSlice(arg reflect.Value) {
	o.BindResultSlice = arg
}

// GetBindResultSlice ...
func (o *Binder) GetBindResultSlice() reflect.Value {
	return o.BindResultSlice
}

// SetBindFields ...
func (o *Binder) SetBindFields(arg []string) {
	o.BindFields = arg
}

// GetBindFields ...
func (o *Binder) GetBindFields() []string {
	return o.BindFields
}

// SetBindType ...
func (o *Binder) SetBindType(arg BindType) {
	o.BindType = arg
}

// GetBindType ...
func (o *Binder) GetBindType() BindType {
	return o.BindType
}

// SetBindAll ...
func (o *Binder) SetBindAll(arg []Data) {
	o.BindAll = arg
}

// GetBindAll ...
func (o *Binder) GetBindAll() []Data {
	return o.BindAll
}

// ResetBinder ...
func (o *Binder) ResetBinder() {
	switch o.BindType {
	case OBJECT_STRUCT, OBJECT_MAP, OBJECT_MAP_T:
		o.SetBindOrigin(nil)
	case OBJECT_STRUCT_SLICE, OBJECT_MAP_SLICE, OBJECT_MAP_SLICE_T:
		var rvResult = o.GetBindResultSlice()
		rvResult.Set(rvResult.Slice(0, 0))
	default:
		o.SetBindAll([]Data{})
	}
}
