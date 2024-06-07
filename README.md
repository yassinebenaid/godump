<div align="center"> 

<div width="50px" height="50px">

![binoculars (3)](https://github.com/yassinebenaid/godump/assets/101285507/f2d40c7a-6f5c-4dd9-9580-3accc74efeb4)

</div>

<h1> godump </h1>
</div>

<div align="center">


[![Tests](https://github.com/yassinebenaid/godump/actions/workflows/test.yml/badge.svg)](https://github.com/yassinebenaid/godump/actions/workflows/test.yml)
[![Version](https://badge.fury.io/gh/yassinebenaid%2Fgodump.svg)](https://badge.fury.io/gh/yassinebenaid%2Fgodump)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENCE)

</div>

A versatile Go library designed to output any Go variable in a structured and colored format.

This library is especially useful for debugging and testing when the standard fmt library falls short in displaying arbitrary data effectively. It can also serve as a powerful logging adapter, providing clear and readable logs for both development and production environments.

## Get Started

Install the library:

```bash
go get -u github.com/yassinebenaid/godump
```

Then use the **Dump** function , you can pass any variable of any type:

```go
package main

import (
	"github.com/yassinebenaid/godump"
)

func main() {
	var anything any = "this can be anything"
	godump.Dump(anything)
}

```

# Demo

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

type Slice []int
type User struct {
	Name   string
	Friend *User
}

func main() {

	person1 := User{"test", nil}
	person2 := User{"test 2", &person1}
	person3 := User{"test 3", &person2}
	person1.Friend = &person3

	godump.Dump(map[any]any{
		"uint":         uint(100),
		"int":          1234,
		"signed-int":   -1234,
		"float":        1234.5678,
		"signed-float": -1234.5678,
		"slice":        []int{1, 2, 3},
		"typed-slice":  Slice{1, 2, 3},
		"channel":      make(chan string, 10),
		"map": map[complex64]bool{
			0xf4a5c5d: true,
			0xa0bff6e: false,
		},
		"inline-struct": struct {
			Pointer      *User
			privateField string
			Err          SomeError
			SomeType     DefinedType
		}{
			Pointer:      &person3,
			privateField: "private",
			Err:          SomeError{"some random error"},
			SomeType: DefinedType{
				"hello world",
				45324,
				true,
				698.4521,
			},
		},
		"func": func(string) string { return "" },
	})

}

```

Output:

![demo](./demo/demo.png)
