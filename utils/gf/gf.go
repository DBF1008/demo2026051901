package gf

type (
	Map        = map[string]interface{}
	MapAnyAny  = map[interface{}]interface{}
	MapAnyStr  = map[interface{}]string
	MapAnyInt  = map[interface{}]int
	MapStrAny  = map[string]interface{}
	MapStrStr  = map[string]string
	MapStrInt  = map[string]int
	MapIntAny  = map[int]interface{}
	MapIntStr  = map[int]string
	MapIntInt  = map[int]int                 // MapIntInt is alias of frequently-used map type map[int]int.
	MapAnyBool = map[interface{}]bool
	MapStrBool = map[string]bool
	MapIntBool = map[int]bool
)

type (
	List        = []Map        // List is alias of frequently-used slice type []Map.
	ListAnyAny  = []MapAnyAny
	ListAnyStr  = []MapAnyStr
	ListAnyInt  = []MapAnyInt
	ListStrAny  = []MapStrAny
	ListStrStr  = []MapStrStr
	ListStrInt  = []MapStrInt
	ListIntAny  = []MapIntAny
	ListIntStr  = []MapIntStr
	ListIntInt  = []MapIntInt
	ListAnyBool = []MapAnyBool
	ListStrBool = []MapStrBool
	ListIntBool = []MapIntBool
)

type (
	Slice    = []interface{} // Slice is alias of frequently-used slice type []interface{}.
	SliceAny = []interface{}
	SliceStr = []string      // SliceStr is alias of frequently-used slice type []string.
	SliceInt = []int         // SliceInt is alias of frequently-used slice type []int.
)

type (
	Array    = []interface{} // Array is alias of frequently-used slice type []interface{}.
	ArrayAny = []interface{}
	ArrayStr = []string      // ArrayStr is alias of frequently-used slice type []string.
	ArrayInt = []int         // ArrayInt is alias of frequently-used slice type []int.
)
