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
	})
	Convey("Test corner cases and errors", t, func(c C) {
	})
}
