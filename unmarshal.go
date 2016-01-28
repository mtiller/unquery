package unquery

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
)

func Unmarshal(query string, sig Signature, pv interface{}) error {
	parsed, err := url.ParseQuery(query)
	if err != nil {
		return fmt.Errorf("Error parsing query string: %v", err)
	}
	pvt := reflect.TypeOf(pv)
	if pvt.Kind() != reflect.Ptr {
		return fmt.Errorf("Unmarshal expected a pointer to a %s", sig.Type.String())
	}
	vt := pvt.Elem()
	if sig.Type != vt {
		return fmt.Errorf("Unmarshal passed a pointer to wrong type, "+
			"expected %s, but got %s (%s)", sig.Type.String(), vt.String(),
			pvt.String())
	}

	for key, values := range parsed {
		log.Printf("Processing %s (%v)", key, values)
	}

	return nil
}
