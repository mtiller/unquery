package unquery

import (
	"log"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	Convey("Test unmarshalling of singleton", t, func(c C) {
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
	Convey("Test unmarshalling of optional", t, func(c C) {
		v := Sample2{}

		sig, err := Scan(v)
		NoError(c, err)
		Equals(c, v.Optional, nil)

		copy := Sample2{}

		err = Unmarshal("", sig, &copy)
		NoError(c, err)
		Equals(c, copy.Optional, nil)

		err = Unmarshal("Optional=1", sig, &copy)
		NoError(c, err)
		NotNil(c, copy.Optional)
		Equals(c, *copy.Optional, "1")
	})
	Convey("Test unmarshalling of multiple", t, func(c C) {
		v := Sample3{}

		sig, err := Scan(v)
		NoError(c, err)
		Equals(c, v.Multiple, nil)

		copy := Sample3{}

		err = Unmarshal("", sig, &copy)
		NoError(c, err)
		Equals(c, copy.Multiple, nil)

		copy.Multiple = []bool{false, true}

		err = Unmarshal("", sig, &copy)
		NoError(c, err)
		Resembles(c, copy.Multiple, []bool{false, true})

		err = Unmarshal("Multiple=true&Multiple=yes&Multiple=0&Multiple=No",
			sig, &copy)
		NoError(c, err)
		NotNil(c, copy.Multiple)
		Resembles(c, copy.Multiple, []bool{true, true, false, false})
	})
	Convey("Test unmarshalling of fixed", t, func(c C) {
		v := Sample4{}

		sig, err := Scan(v)
		NoError(c, err)
		Resembles(c, v.Fixed, [5]uint{0, 0, 0, 0, 0})

		copy := Sample4{}

		err = Unmarshal("", sig, &copy)
		IsError(c, err)

		err = Unmarshal("Fixed=1&Fixed=2&Fixed=3&Fixed=4&Fixed=5",
			sig, &copy)
		NoError(c, err)
		NotNil(c, copy.Fixed)
		Resembles(c, copy.Fixed, [5]uint{1, 2, 3, 4, 5})
	})
	Convey("Test unmarshalling of tagged", t, func(c C) {
		v := Sample5{}

		sig, err := Scan(v)
		NoError(c, err)
		Equals(c, v.Tagged, "")

		copy := Sample5{}

		err = Unmarshal("", sig, &copy)
		IsError(c, err)

		copy.Tagged = "IsTagged"

		err = Unmarshal("", sig, &copy)
		IsError(c, err)

		err = Unmarshal("Tagged=Ignore", sig, &copy)
		IsError(c, err)

		err = Unmarshal("tagged=ItWorked", sig, &copy)
		NoError(c, err)
		Resembles(c, copy.Tagged, "ItWorked")
	})
	Convey("Test different int sizes", t, func(c C) {
		v := Sample7{}

		sig, err := Scan(v)
		NoError(c, err)

		copy := Sample7{}

		err = Unmarshal("Int=2&Int8=3&Int16=4&Int32=5&Int64=6&UInt=7&UInt8=8&UInt16=9&UInt32=10&UInt64=11", sig, &copy)
		NoError(c, err)
		Resembles(c, copy, Sample7{
			Int:    2,
			Int8:   3,
			Int16:  4,
			Int32:  5,
			Int64:  6,
			UInt:   7,
			UInt8:  8,
			UInt16: 9,
			UInt32: 10,
			UInt64: 11,
		})
	})
	Convey("Check example", t, func(c C) {
		v := Example1{
			unexportedData: true,
		}

		sig, err := Scan(v)
		NoError(c, err)

		copy1 := Example1{}

		err = Unmarshal("Message=Hello&Vec=.1&Vec=.2&Vec=.3&names=bill", sig, &copy1)
		NoError(c, err)
		Resembles(c, copy1, Example1{
			unexportedData: true,
			Message:        "Hello",
			Weight:         nil,
			Vec:            [3]float64{0.1, 0.2, 0.3},
			Names:          []string{"bill"},
		})

		copy2 := Example1{}

		err = Unmarshal("Message=Hello&Vec=.1&Vec=.2&Vec=.3&weight=120", sig, &copy2)
		NoError(c, err)
		Equals(c, copy2.unexportedData, true)
		Equals(c, copy2.Message, "Hello")
		Equals(c, *copy2.Weight, 120)
		Equals(c, copy2.Vec, [3]float64{0.1, 0.2, 0.3})
		IsNil(c, copy2.Names)
	})
}

func TestErrors(t *testing.T) {
	Convey("Check for error when passing a struct", t, func(c C) {
		v := Sample1{}
		str := "Singleton=1"
		sig, err := Scan(v)
		NoError(c, err)
		err = Unmarshal(str, sig, Sample1{})
		IsError(c, err)
	})
	Convey("Check for error when passing a bogus query string", t, func(c C) {
		v := Sample1{}
		str := "Singleton=%5"
		sig, err := Scan(v)
		NoError(c, err)
		err = Unmarshal(str, sig, &Sample1{})
		log.Printf("err = %v", err)
		IsError(c, err)
	})
	Convey("Check for error when passing wrong type", t, func(c C) {
		str := "Singleton=1"
		v := Sample1{ignore: "true"}

		sig, err := Scan(v)
		NoError(c, err)
		Equals(c, v.Singleton, 0)

		copy := Sample2{}
		err = Unmarshal(str, sig, &copy)
		IsError(c, err)
	})
	Convey("Check for bogus integer values", t, func(c C) {
		v := Sample7{}

		sig, err := Scan(v)
		NoError(c, err)

		copy := Sample7{}

		err = Unmarshal("Int=2.5&Int8=3&Int16=4&Int32=5&Int64=6&UInt=7&UInt8=8&UInt16=9&UInt32=10&UInt64=11", sig, &copy)
		IsError(c, err)
	})
	Convey("Check for errors when passing too few values", t, func(c C) {
		v := Sample4{}

		sig, err := Scan(v)
		NoError(c, err)

		copy := Sample4{}

		err = Unmarshal("Fixed=1&Fixed=2&Fixed=3&Fixed=4", sig, &copy)
		IsError(c, err)
	})
	Convey("Check for errors when passing too many values", t, func(c C) {
		v := Sample4{}

		sig, err := Scan(v)
		NoError(c, err)

		copy := Sample4{}

		err = Unmarshal("Fixed=1&Fixed=2&Fixed=3&Fixed=4&Fixed=5&Fixed=6",
			sig, &copy)
		IsError(c, err)
	})
	Convey("Check for errors when passing wrong type to array", t, func(c C) {
		v := Sample4{}

		sig, err := Scan(v)
		NoError(c, err)

		copy := Sample4{}

		err = Unmarshal("Fixed=1&Fixed=2&Fixed=3&Fixed=Four&Fixed=5",
			sig, &copy)
		IsError(c, err)
	})
	Convey("Check for errors when passing wrong type to slice", t, func(c C) {
		v := Sample3{}

		sig, err := Scan(v)
		NoError(c, err)
		Equals(c, v.Multiple, nil)

		copy := Sample3{}

		err = Unmarshal("Multiple=true&Multiple=5&Multiple=0&Multiple=No",
			sig, &copy)
		IsError(c, err)
	})
	Convey("Check for error when parsing bad float", t, func(c C) {
		v := Example1{
			unexportedData: true,
		}

		sig, err := Scan(v)
		NoError(c, err)

		copy1 := Example1{}

		err = Unmarshal("Message=Hello&Vec=.1&Vec=seven&Vec=.3&names=bill",
			sig, &copy1)
		IsError(c, err)
	})
}
