package gform

const (
	// DriverMsSql ...
	DriverMsSql = "mssql"
)

// BuilderMsSql ...
type BuilderMsSql struct {
	FieldQuotesDefault
	//IOrm
	driver string
}

func init() {
	var builder = &BuilderMsSql{driver: DriverMsSql}
	NewBuilderDriver().Register(DriverMsSql, builder)
}

// NewBuilderMsSql ...
func NewBuilderMsSql() *BuilderMsSql {
	return new(BuilderMsSql)
}

// Clone : a new obj
func (b *BuilderMsSql) Clone() IBuilder {
	return &BuilderMsSql{driver: DriverMsSql}
}

// BuildQuery : build query sql string
func (b *BuilderMsSql) BuildQuery(o IOrm) (sqlStr string, args []interface{}, err error) {
	return NewBuilderDefault(o, NewBuilderMsSql()).SetDriver(b.driver).BuildQuery()
}

// BuildExecut : build execute sql string
func (b *BuilderMsSql) BuildExecute(o IOrm, operType string) (sqlStr string, args []interface{}, err error) {
	return NewBuilderDefault(o, NewBuilderMsSql()).SetDriver(b.driver).BuildExecute(operType)
}
