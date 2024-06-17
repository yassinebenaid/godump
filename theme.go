package godump

import "fmt"

type RGB struct {
	R, G, B int
}

func (rgb *RGB) __(v string) string {
	if rgb == nil {
		return v
	}
	return fmt.Sprintf("\033[38;2;%v;%v;%vm%s\033[0m", rgb.R, rgb.G, rgb.B, v)
}

type Theme struct {
	String         *RGB
	Quotes         *RGB
	Bool           *RGB
	Number         *RGB
	Types          *RGB
	Nil            *RGB
	PointerSign    *RGB
	UnsafePointer  *RGB
	PointerCounter *RGB
	Func           *RGB
	StructField    *RGB
	Chan           *RGB
	Braces         *RGB
}

var DefaultTheme = Theme{
	String:         &RGB{138, 201, 38},
	Quotes:         &RGB{112, 214, 255},
	Bool:           &RGB{249, 87, 56},
	Number:         &RGB{10, 178, 242},
	Types:          &RGB{0, 150, 199},
	PointerSign:    &RGB{205, 93, 0},
	PointerCounter: &RGB{110, 110, 110},
	Nil:            &RGB{219, 57, 26},
	Func:           &RGB{160, 90, 220},
	StructField:    &RGB{189, 176, 194},
	Chan:           &RGB{195, 154, 76},
	UnsafePointer:  &RGB{89, 193, 180},
	Braces:         &RGB{185, 86, 86},
}

// DisableColors disables the colors globally.
//
// Deprecated: As of v0.8.0 this function only sets the [godump.DefaultTheme] to a zero value
func DisableColors() { // TODO: deprecate this function
	DefaultTheme = Theme{}
}
