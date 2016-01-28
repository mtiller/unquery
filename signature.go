package unquery

import (
	"fmt"
	"math"
	"reflect"
	"unicode"
)

type ParameterDetails struct {
	FieldName  string // Needed, if we have kind and offset?
	FieldIndex int
	Array      bool
	Kind       reflect.Kind
	Offset     uintptr
	Min        int
	Max        int
}

type ParameterName string

type Signature struct {
	Type       reflect.Type
	Original   interface{}
	Value      reflect.Value
	Parameters map[ParameterName]ParameterDetails
}

// This is the maximum size of arrays
const UpperLimit = math.MaxInt8

func Scan(v interface{}) (Signature, error) {
	blank := Signature{}

	vt := reflect.TypeOf(v)
	ret := Signature{
		Type:       vt,
		Original:   v,
		Value:      reflect.ValueOf(v),
		Parameters: map[ParameterName]ParameterDetails{},
	}
	if vt.Kind() != reflect.Struct {
		return blank, fmt.Errorf("Value to scan is not a struct: %v", v)
	}

	for i := 0; i < vt.NumField(); i++ {
		ft := vt.Field(i)
		runes := []rune(ft.Name)
		if unicode.IsUpper(runes[0]) {
			pname := ft.Name
			if ft.Tag.Get("unq") != "" {
				pname = ft.Tag.Get("unq")
			}

			min := 1
			max := 1

			pkind := ft.Type.Kind()
			switch ft.Type.Kind() {
			case reflect.Array:
				min = ft.Type.Len()
				max = ft.Type.Len()
				pkind = ft.Type.Elem().Kind()
			case reflect.Slice:
				min = 0
				max = UpperLimit
				pkind = ft.Type.Elem().Kind()
			case reflect.Ptr:
				min = 0
				max = 1
				pkind = ft.Type.Elem().Kind()
			}

			// These types, we can't handle...
			if pkind == reflect.Complex64 ||
				pkind == reflect.Complex128 ||
				pkind == reflect.Chan ||
				pkind == reflect.Func ||
				pkind == reflect.Interface ||
				pkind == reflect.Map ||
				pkind == reflect.Struct ||
				pkind == reflect.UnsafePointer {
				return blank,
					fmt.Errorf("Cannot handle exported variables of type %s",
						pkind.String())
			}

			details := ParameterDetails{
				FieldName:  ft.Name,
				FieldIndex: i,
				Array:      ft.Type.Kind() == reflect.Array,
				Kind:       pkind,
				Offset:     ft.Offset,
				Min:        min,
				Max:        max,
			}
			ret.Parameters[ParameterName(pname)] = details
		}
	}

	return ret, nil
}
