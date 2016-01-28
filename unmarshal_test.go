package unquery

import (
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	Convey("Test unmarshalling of query strings", t, func(c C) {
		str := "Singleton=1"
		v := Sample1{ignore: "true"}

		sig, err := Scan(v)
		NoError(c, err)
		Equals(c, v.Singleton, 0)

		copy := Sample1{}
		err = Unmarshal(str, sig, &copy)
		NoError(c, err)

		Equals(c, copy.ignore, v.ignore)
		Equals(c, copy.Singleton, 1)
	})
	Convey("Check for error when passing a struct", t, func(c C) {
		v := Sample1{}
		str := "Singleton=1"
		sig, err := Scan(v)
		NoError(c, err)
		err = Unmarshal(str, sig, Sample1{})
		IsError(c, err)
	})
}
