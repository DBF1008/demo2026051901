package gform

import (
	"errors"
	"reflect"

	"github.com/gohouse/t"
)

func (dba *Orm) Insert(data ...interface{}) (int64, error) {
	return dba.exec("insert", data...)
}

// insertGetId : insert data and get id
func (dba *Orm) InsertGetId(data ...interface{}) (int64, error) {
	_, err := dba.Insert(data...)
	if err != nil {
		return 0, err
	}
	return dba.GetISession().LastInsertId(), nil
}

// Update : update data
func (dba *Orm) Update(data ...interface{}) (int64, error) {
	return dba.exec("update", data...)
}

func (dba *Orm) Save(data ...interface{}) (int64, error) {
	return dba.exec("update", data...)
}

func (dba *Orm) Force() IOrm {
	dba.force = true
	return dba
}

// Delete : delete data
func (dba *Orm) Delete() (int64, error) {
	return dba.exec("delete")
}

// Delete : delete data
func (dba *Orm) exec(operType string, data ...interface{}) (int64, error) {
	if operType == "insert" || operType == "update" {
		if dba.GetData() == nil {
			if len(data) > 0 {
				dba.Data(data[0])
			} else {
				return 0, GetErr(ERR_PARAMS_MISSING, "Data()")
			}
		}

		//	if dba.GetData() != nil {
		//		dba.Table(dba.GetData())
		rl := reflect.ValueOf(dba.GetData())
		rl2 := reflect.Indirect(rl)

		switch rl2.Kind() {
		case reflect.Struct, reflect.Ptr:
			if tn := rl2.MethodByName("TableName"); tn.IsValid() {
				dba.Table(dba.GetData())
			}
		case reflect.Map:
			if tn := rl2.MethodByName("TableName"); tn.IsValid() {
				dba.Table(dba.GetData())
			}
			if tn := rl.MethodByName("TableName"); tn.IsValid() {
				dba.Table(dba.GetData())
			}
		case reflect.Slice:
			r2 := rl2.Type().Elem()
			r2val := reflect.New(r2)
			switch r2val.Kind() {
			case reflect.Struct, reflect.Ptr:
				if tn := r2val.MethodByName("TableName"); tn.IsValid() {
					dba.Table(dba.GetData())
				}
			case reflect.Map:
				if tn := r2val.MethodByName("TableName"); tn.IsValid() {
					dba.Table(dba.GetData())
				}
			default:
				return 0, errors.New("表名有误")
			}
		}

	}
	sqlStr, args, err := dba.BuildSql(operType)
	if err != nil {
		return 0, err
	}

	return dba.GetISession().Execute(sqlStr, args...)
}

// Increment : auto Increment +1 default
func (dba *Orm) Increment(args ...interface{}) (int64, error) {
	argLen := len(args)
	var field string
	var mode = "+"
	var value = "1"
	switch argLen {
	case 1:
		field = t.New(args[0]).String()
	case 2:
		field = t.New(args[0]).String()
		value = t.New(args[1]).String()
	case 3:
		field = t.New(args[0]).String()
		value = t.New(args[1]).String()
		mode = t.New(args[2]).String()
	default:
		return 0, errors.New("参数数量只允许1个,2个或3个")
	}
	dba.Data(field + "=" + field + mode + value)
	return dba.Update()
}

// Decrement : auto Decrement -1 default
func (dba *Orm) Decrement(args ...interface{}) (int64, error) {
	arglen := len(args)
	switch arglen {
	case 1:
		args = append(args, 1)
		args = append(args, "-")
	case 2:
		args = append(args, "-")
	default:
		return 0, errors.New("Decrement参数个数有误")
	}
	return dba.Increment(args...)
}
