package gform

// IOrmQuery ...
type IOrmQuery interface {
	Select() error
	First() (Data, error)
	Get() ([]Data, error)
	Value(field string) (v interface{}, err error)
	Pluck(field string, fieldKey ...string) (v interface{}, err error)
	Count(args ...string) (int64, error)
	Sum(sum string) (interface{}, error)
	Avg(avg string) (interface{}, error)
	Max(max string) (interface{}, error)
	Min(min string) (interface{}, error)
	Paginate(page ...int) (res Data, err error)
	Chunk(limit int, callback func([]Data) error) (err error)
	ChunkStruct(limit int, callback func() error) (err error)
	Loop(limit int, callback func([]Data) error) (err error)
}
