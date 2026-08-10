package main

import (
	"fmt"
	"gopurs/output/Foreign"
	"gopurs/output/gopurs_runtime"
	"reflect"
)

// Simulate what gopurs compiler generates for Newtype Year
type Year struct {
	V0 gopurs_runtime.Value
}

func main() {
    // Record { year: 1885, approximately: false }
	rec := gopurs_runtime.Record(
		gopurs_runtime.Str("approximately"), gopurs_runtime.Bool(false),
		gopurs_runtime.Str("year"), gopurs_runtime.Int(1885),
	)
	
	// Newtype Year over the Record
	y := Year{V0: rec}
	
	// deepUnbox it
	res := Foreign.UnboxForJSON(y)
	fmt.Printf("%#v\n", res)
}
