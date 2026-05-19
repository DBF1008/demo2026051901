package gform

import (
	"github.com/gohouse/t"
	"testing"
)

type aaa t.MapStringT

func (u *aaa) TableName() string {
	return "users"
}

// type bbb MapRows
type bbb []t.MapStringT

func (u *bbb) TableName() string {
	return "users"
}

type UsersMap Data

func (*UsersMap) TableName() string {
	return "users"
}

type UsersMapSlice []Data

func (u *UsersMapSlice) TableName() string {
	return "users"
}

type Users struct {
	Uid  int64  `orm:"uid"`
	Name string `orm:"name"`
	Age  int64  `orm:"age"`
	Fi   string `orm:"ignore"`
}

func (Users) TableName() string {
	return "users"
}

type Orders struct {
	Id        int     `orm:"id"`
	GoodsName string  `orm:"goodsname"`
	Price     float64 `orm:"price"`
}

func TestEngin(t *testing.T) {
	e := initDB()
	e.SetPrefix("pre_")

	t.Log(e.GetPrefix())

	db := e.GetQueryDB()

	err := db.Ping()

	if err != nil {
		t.Error("gform初始化失败")
	}
	t.Log("gform初始化成功")
	t.Log(e.GetLogger())
}
