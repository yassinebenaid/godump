package godump

import "fmt"

type rgb struct {
	R, G, B int
}

func (rgb *rgb) __(v string) string {
	if rgb == nil {
		return v
	}
	return fmt.Sprintf("\033[38;2;%v;%v;%vm%s\033[0m", rgb.R, rgb.G, rgb.B, v)
}

type theme struct {
	String         *rgb
	Quotes         *rgb
	Bool           *rgb
	Number         *rgb
	Types          *rgb
	Nil            *rgb
	PointerSign    *rgb
	UnsafePointer  *rgb
	PointerCounter *rgb
	Func           *rgb
	StructField    *rgb
	Chan           *rgb
}

var defaultTheme = theme{
	String:         &rgb{138, 201, 38},
	Quotes:         &rgb{112, 214, 255},
	Bool:           &rgb{249, 87, 56},
	Number:         &rgb{10, 178, 242},
	Types:          &rgb{0, 150, 199},
	PointerSign:    &rgb{205, 93, 0},
	PointerCounter: &rgb{110, 110, 110},
	Nil:            &rgb{219, 57, 26},
	Func:           &rgb{160, 90, 220},
	StructField:    &rgb{189, 176, 194},
	Chan:           &rgb{195, 154, 76},
	UnsafePointer:  &rgb{89, 193, 180},
}

// DisableColors disables the colors globally.
func DisableColors() { // TODO: deprecate this function
	defaultTheme = theme{}
}
