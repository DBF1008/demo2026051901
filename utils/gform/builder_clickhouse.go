package gform

const (
	// DriverClickhouse ...
	DriverClickhouse = "clickhouse"
)

// BuilderClickhouse ...
type BuilderClickhouse struct {
	FieldQuotesDefault
	//IOrm
	driver string
}

func init() {
	var builder = &BuilderClickhouse{}
	NewBuilderDriver().Register(DriverClickhouse, builder)
}

// NewBuilderClickhouse ...
func NewBuilderClickhouse() *BuilderClickhouse {
	return new(BuilderClickhouse)
}

// Clone : a new obj
func (b *BuilderClickhouse) Clone() IBuilder {
	return &BuilderClickhouse{}
}

// BuildQuery : build query sql string
func (b *BuilderClickhouse) BuildQuery(o IOrm) (sqlStr string, args []interface{}, err error) {
	return NewBuilderDefault(o, NewBuilderClickhouse()).SetDriver(DriverClickhouse).BuildQuery()
}

// BuildExecut : build execute sql string
func (b *BuilderClickhouse) BuildExecute(o IOrm, operType string) (sqlStr string, args []interface{}, err error) {
	return NewBuilderDefault(o, NewBuilderClickhouse()).SetDriver(DriverClickhouse).BuildExecute(operType)
}
