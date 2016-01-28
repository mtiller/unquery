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
