package unquery

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func Unmarshal(query string, sig Signature, p interface{}) error {
	parsed, err := url.ParseQuery(query)
	if err != nil {
		return fmt.Errorf("Error parsing query string: %v", err)
	}
	pt := reflect.TypeOf(p)
	if pt.Kind() != reflect.Ptr {
		return fmt.Errorf("Unmarshal expected a pointer to a %s", sig.Type.String())
	}
	pit := pt.Elem()
	if sig.Type != pit {
		return fmt.Errorf("Unmarshal passed a pointer to wrong type, "+
			"expected %s, but got %s (%s)", sig.Type.String(), pit.String(),
			pt.String())
	}

	pv := reflect.ValueOf(p)
	piv := pv.Elem()

	// Initialize the unmarshalled value to be equal to the default
	// value of the signature.
	piv.Set(sig.Value)

	for pname, details := range sig.Parameters {
		values, exists := parsed[string(pname)]
		if !exists {
			if details.Min > 0 {
				return fmt.Errorf("No value given for %s", pname)
			}
			continue
		}

		nv := len(values)
		if nv < details.Min {
			return fmt.Errorf("Not enough values given for %s", pname)
		}
		if nv > details.Max {
			return fmt.Errorf("Too many values given for %s", pname)
		}

		// Optional
		if details.Min == 0 && details.Max == 1 {
		}
		// Slice
		if details.Min == 0 && details.Max == UpperLimit {
		}
		// Array
		if details.Array {
		}
		// Scalar
		if !details.Array && details.Min == 1 && details.Max == 1 {
			parseAs(values[0], details.Kind, piv.Field(details.FieldIndex))
		}
	}

	return nil
}

func parseAs(str string, kind reflect.Kind, dst reflect.Value) error {
	// Handle strings
	if kind == reflect.String {
		dst.SetString(str)
		if !dst.IsValid() {
			return fmt.Errorf("Error setting field of type %s to %s",
				kind.String(), str)
		}
		return nil
	}

	// Handle booleans
	if kind == reflect.Bool {
		if strings.ToLower(str) == "true" ||
			strings.ToLower(str) == "yes" ||
			str == "1" {
			dst.SetBool(true)
			return nil
		}
		if strings.ToLower(str) == "false" ||
			strings.ToLower(str) == "no" ||
			str == "0" {
			dst.SetBool(false)
			return nil
		}
		return fmt.Errorf("String '%s' not recognized as a boolean value", str)
	}

	// Handle integers
	if kind >= reflect.Int && kind <= reflect.Int64 {
		size := 0
		if kind == reflect.Int8 {
			size = 8
		}
		if kind == reflect.Int16 {
			size = 16
		}
		if kind == reflect.Int32 {
			size = 32
		}
		if kind == reflect.Int64 {
			size = 64
		}
		i, err := strconv.ParseInt(str, 10, size)
		if err != nil {
			return fmt.Errorf("Error parsing '%s' as an integer: %v", str, err)
		}
		dst.SetInt(i)
		if !dst.IsValid() {
			return fmt.Errorf("Error setting field of type %s to %v",
				kind.String(), i)
		}
		return nil
	}

	// Handle unsigned integers
	if kind >= reflect.Uint && kind <= reflect.Uint64 {
		size := 0
		if kind == reflect.Uint8 {
			size = 8
		}
		if kind == reflect.Uint16 {
			size = 16
		}
		if kind == reflect.Uint32 {
			size = 32
		}
		if kind == reflect.Uint64 {
			size = 64
		}
		i, err := strconv.ParseUint(str, 10, size)
		if err != nil {
			return fmt.Errorf("Error parsing '%s' as an unsigned integer: %v",
				str, err)
		}
		dst.SetUint(i)
		if !dst.IsValid() {
			return fmt.Errorf("Error setting field of type %s to %v",
				kind.String(), i)
		}
		return nil
	}

	// Handle floats
	if kind >= reflect.Float32 || kind <= reflect.Float64 {
		size := 32
		if kind == reflect.Float64 {
			size = 64
		}
		i, err := strconv.ParseFloat(str, size)
		if err != nil {
			return fmt.Errorf("Error parsing '%s' as a float: %v", str, err)
		}
		dst.SetFloat(i)
		if !dst.IsValid() {
			return fmt.Errorf("Error setting field of type %s to %v",
				kind.String(), i)
		}
		return nil
	}

	return fmt.Errorf("Unable to convert strings to %s", kind.String())
}
