package unquery

import (
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
	Convey("Check for error when passing a struct", t, func(c C) {
		v := Sample1{}
		str := "Singleton=1"
		sig, err := Scan(v)
		NoError(c, err)
		err = Unmarshal(str, sig, Sample1{})
		IsError(c, err)
	})
}
