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
	return UnmarshalValues(parsed, sig, p)
}

func UnmarshalValues(parsed url.Values, sig Signature, p interface{}) error {
	// Get the type of p
	pt := reflect.TypeOf(p)

	// Make sure that the type of p (pt) is a pointer type
	if pt.Kind() != reflect.Ptr {
		return fmt.Errorf("Unmarshal expected a pointer to a %s", sig.Type.String())
	}

	// Get the type that p points to...
	pet := pt.Elem()

	// If it isn't the same type as what our signature is for, error out
	if sig.Type != pet {
		return fmt.Errorf("Unmarshal passed a pointer to wrong type, "+
			"expected %s, but got %s (%s)", sig.Type.String(), pet.String(),
			pt.String())
	}

	// Now get get Value reprentations of both p and the what p is pointing to...
	pv := reflect.ValueOf(p)
	pev := pv.Elem()

	// Create a new instance of the structure we are working with here
	newv := reflect.New(pev.Type())

	// Then give it an initial value that matches the default value associated
	// with the signature.  This is mainly a way to copy the *unexported*
	// fields
	newv.Elem().Set(sig.Value)

	// Now copy all the **exported** field values from p into this new
	// value.  This gives us a value that has the default values for
	// the unexported variables but doesn't touch the exported field
	// values from p
	for _, details := range sig.Parameters {
		newv.Elem().Field(details.FieldIndex).Set(pev.Field(details.FieldIndex))
	}

	// Now set the value that p points to be this new hybrid initialized value
	pev.Set(newv.Elem())

	// Now lets loop over the exported fields again and given them
	// new values based on the contents of the query string values...
	for pname, details := range sig.Parameters {
		// Get the value and type for the field associated with the
		// current parameter
		fv := pev.Field(details.FieldIndex)
		ft := pet.Field(details.FieldIndex)

		//fv.Set(pev.Field(details.FieldIndex))

		// Is this parameter in the query string?
		values, exists := parsed[string(pname)]

		// If not, just check to see if one is required.  If not,
		// then just skip this parameter and leave its current
		// value alone.
		if !exists || len(values) == 0 {
			if details.Min > 0 {
				return fmt.Errorf("Expected at least %d values for %s",
					details.Min, pname)
			}
			continue
		}

		nv := len(values)

		// Check to make sure we have as many entries in values as
		// we expect.  If not, throw an error
		if nv < details.Min {
			return fmt.Errorf("Not enough values given for %s", pname)
		}
		if nv > details.Max {
			return fmt.Errorf("Too many values given for %s", pname)
		}

		// Assume we there is a problem.  This will get overridden
		// if we find a way to use the values...
		var perr error = fmt.Errorf("Unhandled case for parameters %s"+
			"in UnmarshalValues", pname)

		// This parameter is optional...
		if details.Min == 0 && details.Max == 1 {
			fv.Set(reflect.New(ft.Type.Elem()))
			perr = parseAs(values[0], details.Kind, fv.Elem())
		}

		// This parameter is a slice (can be any size)...
		if details.Min == 0 && details.Max == UpperLimit {
			// Allocate a slice big enough to hold the values
			slice := reflect.MakeSlice(ft.Type, nv, nv)
			// Set the field equal to this new slice
			fv.Set(slice)
			// Now loop over the values, parse them and insert
			// them at the appropriate spot in the slice
			for ind, value := range values {
				perr = parseAs(value, details.Kind, fv.Index(ind))
				if perr != nil {
					// If we get an error, just break out.  That error
					// value will be seen at the end
					break
				}
			}
		}

		// This parameter is an array (i.e., fixed size)
		if details.Array {
			// Make sure we have exactly the number of values we need
			// to fill the array
			if nv != details.Min {
				return fmt.Errorf("Expected %d values for %s, but got %d: %v",
					details.Min, pname, nv, values)
			}
			// Now loop over all values in the query string, parse them
			// and then insert them into the array at the approriate index
			for ind, value := range values {
				perr = parseAs(value, details.Kind, fv.Index(ind))
				if perr != nil {
					break
				}
			}
		}

		// This parameter is a simple scalar
		if !details.Array && details.Min == 1 && details.Max == 1 {
			// Just parse it and assign it to the current field
			perr = parseAs(values[0], details.Kind, pev.Field(details.FieldIndex))
		}

		// If any errors occured parsing any of the query string values for
		// the expected kind, report an error
		if perr != nil {
			return fmt.Errorf("Error processing parameter %s: %v", pname, perr)
		}
	}

	// If we get here, everything worked.
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
