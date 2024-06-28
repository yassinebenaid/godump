<div align="center">

<div width="50px" height="50px">

![binoculars (3)](https://github.com/yassinebenaid/godump/assets/101285507/f2d40c7a-6f5c-4dd9-9580-3accc74efeb4)

</div>

<h1> godump </h1>
</div>

<div align="center">

[![Tests](https://github.com/yassinebenaid/godump/actions/workflows/test.yml/badge.svg)](https://github.com/yassinebenaid/godump/actions/workflows/test.yml)
[![codecov](https://codecov.io/github/yassinebenaid/godump/graph/badge.svg?token=EAZNA85AIS)](https://codecov.io/github/yassinebenaid/godump)
[![Go Report Card](https://goreportcard.com/badge/github.com/yassinebenaid/godump)](https://goreportcard.com/report/github.com/yassinebenaid/godump)
[![Version](https://badge.fury.io/gh/yassinebenaid%2Fgodump.svg)](https://badge.fury.io/gh/yassinebenaid%2Fgodump)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENCE)
[![Go Reference](https://pkg.go.dev/badge/github.com/yassinebenaid/godump.svg)](https://pkg.go.dev/github.com/yassinebenaid/godump)


</div>

A versatile Go library designed to output any Go variable in a structured and colored format.

This library is especially useful for debugging and testing when the standard fmt library falls short in displaying arbitrary data effectively. It can also serve as a powerful logging adapter, providing clear and readable logs for both development and production environments.

`godump` is not here to replace the `fmt` package. Instead, it provides an extension to what the `fmt.Printf("%#v")` can do.

## Why godump

- ability to pretty print values of all types
- coloreful output
- unexported structs are dumped too
- pointers are followed and recursive pointers are taken in mind ([see examples](#example-3))
- customizable
- zero dependencies

## Get Started

Install the library:

```bash
go get -u github.com/yassinebenaid/godump
```

Then use the **Dump** function:

```go
package main

import (
	"github.com/yassinebenaid/godump"
)

func main() {
	godump.Dump("Anything")
}

```

## Cusomization

If you need more control over the output. Use the `Dumper`

```go
package main

import (
	"os"

	"github.com/yassinebenaid/godump"
)

func main() {

	var v = "Foo Bar"
	var d = godump.Dumper{
		Indentation:       "  ",
		HidePrivateFields: false,
		Theme: godump.Theme{
			String: godump.RGB{R: 138, G: 201, B: 38},
			// ...
		},
	}

	d.Print(v)
	d.Println(v)
	d.Fprint(os.Stdout, v)
	d.Fprintln(os.Stdout, v)
	d.Sprint(v)
	d.Sprintln(v)
}

```

## Demo

### Example 1.

```go
package main

import (
	"os"

	"github.com/yassinebenaid/godump"
)

func main() {
	godump.Dump(os.Stdout)
}

```

Output:

![stdout](./demo/stdout.png)

### Example 2.

```go
package main

import (
	"net"

	"github.com/yassinebenaid/godump"
)

func main() {
	godump.Dump(net.Dialer{})
}

```

Output:

![dialer](./demo/dialer.png)

### Example 3.

This example shows how recursive pointers are handeled

```go
package main

import (
	"github.com/yassinebenaid/godump"
)

func main() {
	type User struct {
		Name       string
		age        int
		BestFriend *User
	}

	me := User{
		Name: "yassinebenaid",
		age:  22,
	}

    // This creates a ring
	me.BestFriend = &me

	godump.Dump(me)
}
```

Output:

![pointer](./demo/pointer.png)

**Note**: the `&@1` syntax means that this is an address of the value with the id of 1. which is marked as `#1`.

For more examples, please have a look at [dumper_test](./dumper_test.go) along with [testdata](./testdata)

## Contribution

Please read [CONTRIBUTING guildlines](.github/CONTRIBUTING.md)
