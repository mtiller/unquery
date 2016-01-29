package unquery

import (
	"reflect"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
	"testing"
)

func Check(c C, sig Signature, name string, pname string,
	kind reflect.Kind, min int, max int) {
	Equals(c, len(sig.Parameters), 1)
	details, exists := sig.Parameters[ParameterName(pname)]
	IsTrue(c, exists)

	Equals(c, details.Kind, kind)
	Equals(c, details.FieldName, name)
	Equals(c, details.Min, min)
	Equals(c, details.Max, max)
}

func TestScanning(t *testing.T) {
	Convey("Test signature generation", t, func(c C) {
		Convey("Test non-structs", func(c C) {
			_, err := Scan(5)
			IsError(c, err)

			_, err = Scan("hello")
			IsError(c, err)

			_, err = Scan(&Sample1{})
			IsError(c, err)

			_, err = Scan(map[string]int{})
			IsError(c, err)

			_, err = Scan(Sample6{})
			IsError(c, err)
		})
		Convey("Test hidden fields", func(c C) {
			v := Sample1{}
			t := reflect.TypeOf(v)
			sig, err := Scan(v)
			NoError(c, err)
			Equals(c, sig.Type, t)

			Check(c, sig, "Singleton", "Singleton", reflect.Int, 1, 1)
		})
		Convey("Test optional fields", func(c C) {
			v := Sample2{}
			t := reflect.TypeOf(v)
			sig, err := Scan(v)
			NoError(c, err)
			Equals(c, sig.Type, t)

			Check(c, sig, "Optional", "Optional", reflect.String, 0, 1)
		})
		Convey("Test slice fields", func(c C) {
			v := Sample3{}
			t := reflect.TypeOf(v)
			sig, err := Scan(v)
			NoError(c, err)
			Equals(c, sig.Type, t)

			Check(c, sig, "Multiple", "Multiple", reflect.Bool, 0, UpperLimit)
		})
		Convey("Test array fields", func(c C) {
			v := Sample4{}
			t := reflect.TypeOf(v)
			sig, err := Scan(v)
			NoError(c, err)
			Equals(c, sig.Type, t)

			Check(c, sig, "Fixed", "Fixed", reflect.Uint, 5, 5)
		})
		Convey("Test tagged fields", func(c C) {
			v := Sample5{}
			t := reflect.TypeOf(v)
			sig, err := Scan(v)
			NoError(c, err)
			Equals(c, sig.Type, t)

			Check(c, sig, "Tagged", "tagged", reflect.String, 1, 1)
		})
	})
}
