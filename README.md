# godump

A simple GO library to dump any GO variable in a structured and colored way.

Useful during development when the standard `fmt` library doesn't help you extract arbitrary data.

## Example

```go
package main

import (
	"github.com/yassinebenaid/godump"
)

type SomeError struct{ E string }

func (e SomeError) Error() string {
	return e.E
}

type DefinedType struct {
	Field  string
	Field2 int
	Field3 bool
	Field4 float64
}

func main() {

	var a int = 55644133

	godump.Dump(map[any]any{
		"int":          1234,
		"signed-int":   -1234,
		"float":        1234.5678,
		"signed-float": -1234.5678,
		"slice":        []int{1, 2, 3},
		"map": map[complex64]bool{
			0xf4a5c5d: true,
			0xa0bff6e: false,
		},
		"inline-struct": struct {
			Pointer      *int
			privateField string
			Err          SomeError
			SomeType     DefinedType
		}{
			Pointer:      &a,
			privateField: "private",
			Err:          SomeError{"some random error"},
			SomeType: DefinedType{
				"hello world",
				45324,
				true,
				698.4521,
			},
		},
	})

}

```

Output:
![example](./example/example.png)
