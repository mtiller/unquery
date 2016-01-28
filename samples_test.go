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

type Example1 struct {
	unexportedData bool
	Message        string
	Weight         *int `unq:"weight"`
	Vec            [3]float64
	Names          []string `unq:"names"`
}
