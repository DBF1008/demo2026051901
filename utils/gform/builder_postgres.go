package gform

import (
	"fmt"
	"regexp"
)

const (
	// DriverPostgres ...
	DriverPostgres = "postgres"
)

// BuilderPostgres ...
type BuilderPostgres struct {
	//IOrm
	driver string
}

func init() {
	var builder = &BuilderPostgres{driver: DriverPostgres}
	NewBuilderDriver().Register(DriverPostgres, builder)
}

// NewBuilderPostgres ...
func NewBuilderPostgres() *BuilderPostgres {
	return new(BuilderPostgres)
}

// Clone : a new obj
func (b *BuilderPostgres) Clone() IBuilder {
	return &BuilderPostgres{driver: DriverPostgres}
}

// BuildQuery : build query sql string
func (b *BuilderPostgres) BuildQuery(o IOrm) (sqlStr string, args []interface{}, err error) {
	return NewBuilderDefault(o, NewBuilderPostgres()).SetDriver(b.driver).BuildQuery()
}

// BuildExecut : build execute sql string
func (b *BuilderPostgres) BuildExecute(o IOrm, operType string) (sqlStr string, args []interface{}, err error) {
	return NewBuilderDefault(o, NewBuilderPostgres()).SetDriver(b.driver).BuildExecute(operType)
}

func (*BuilderPostgres) AddFieldQuotes(field string) string {
	reg := regexp.MustCompile(`^\w+$`)
	if reg.MatchString(field) {
		return fmt.Sprintf(`"%s"`, field)
	}
	return field
}
