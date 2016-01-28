package unquery

import (
	"fmt"
	"reflect"
	"unicode"
)

type ParameterDetails struct {
	FieldName string // Needed, if we have kind and offset?
	Kind      reflect.Kind
	Offset    uintptr
	Min       int
	Max       int
}

type ParameterName string

type Signature struct {
	Type       reflect.Type
	Parameters map[ParameterName]ParameterDetails
}

func Scan(v interface{}) (Signature, error) {
	blank := Signature{}

	vt := reflect.TypeOf(v)
	ret := Signature{
		Type:       vt,
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
				max = -1
				pkind = ft.Type.Elem().Kind()
			case reflect.Ptr:
				min = 0
				max = 1
				pkind = ft.Type.Elem().Kind()
			}

			// These types, we can't handle...
			if pkind == reflect.Complex64 {
				continue
			}
			if pkind == reflect.Complex128 {
				continue
			}
			if pkind == reflect.Chan {
				continue
			}
			if pkind == reflect.Func {
				continue
			}
			if pkind == reflect.Interface {
				continue
			}
			if pkind == reflect.Map {
				continue
			}
			if pkind == reflect.Struct {
				continue
			}
			if pkind == reflect.UnsafePointer {
				continue
			}

			details := ParameterDetails{
				FieldName: ft.Name,
				Kind:      pkind,
				Offset:    ft.Offset,
				Min:       min,
				Max:       max,
			}
			ret.Parameters[ParameterName(pname)] = details
		}
	}

	return ret, nil
}
