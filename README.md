# Unquery

This Go package provides functions that parse query strings and apply
them to a Go struct.

I had very specific needs here.  I needed the unmarshal function to
copy the values of any **unexported** fields from a prototype instance
but then overwrite exported fields with the values of query variables.
In other words, the unmarshal function fuses together unexported
fields form the prototype with exported fields from the query string.

The library handles string, int, uint, bool and float types.  If a
value in a struct is meant to be optional, it should have a pointer
type.  If it is meant to have multiple values, it should be either a
slice (for any arbitrary number of elements) or an array (for a fixed
number of elements).  Finally, if it is just one of the supported
primitives, it is assumed to be required.

Along the way, lots of checking and casting is done to ensure that the
resulting struct is validated.
