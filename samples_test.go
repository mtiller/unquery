package unquery

type Sample1 struct {
	ignore    string
	Singleton int
}

type Sample2 struct {
	Optional *string
}

type Sample3 struct {
	Multiple []bool
}

type Sample4 struct {
	Fixed [5]uint
}

type Sample5 struct {
	Tagged string `unq:"tagged"`
}

type Sample6 struct {
	NestedStruct Sample5 // Not allowed
}

type Sample7 struct {
	Int    int
	Int8   int8
	Int16  int16
	Int32  int32
	Int64  int64
	UInt   uint
	UInt8  uint8
	UInt16 uint16
	UInt32 uint32
	UInt64 uint64
}

type Example1 struct {
	unexportedData bool
	Message        string
	Weight         *int `unq:"weight"`
	Vec            [3]float64
	Names          []string `unq:"names"`
}
