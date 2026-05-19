package gform

import (
	"fmt"
	"regexp"
)

// IBuilder ...
type IBuilder interface {
	IFieldQuotes
	BuildQuery(orm IOrm) (sqlStr string, args []interface{}, err error)
	BuildExecute(orm IOrm, operType string) (sqlStr string, args []interface{}, err error)
	Clone() IBuilder
	//GetIOrm() IOrm
}

type IFieldQuotes interface {
	AddFieldQuotes(field string) string
}

type FieldQuotesDefault struct {
}

func (FieldQuotesDefault) AddFieldQuotes(field string) string {
	reg := regexp.MustCompile(`^\w+$`)
	if reg.MatchString(field) {
		return fmt.Sprintf("`%s`", field)
	}
	return field
}
