package godump

import "fmt"

type rgb struct {
	R, G, B int
}

func (rgb *rgb) apply(v string) string {
	if rgb == nil {
		return v
	}
	return fmt.Sprintf("\033[38;2;%v;%v;%vm%s\033[0m", rgb.R, rgb.G, rgb.B, v)
}

type theme struct {
	String          *rgb
	Quotes          *rgb
	Bool            *rgb
	Number          *rgb
	VarType         *rgb
	Pointer         *rgb
	Func            *rgb
	StructField     *rgb
	StructFieldHash *rgb
}

var defaultTheme = theme{
	String:          &rgb{138, 201, 38},
	Quotes:          &rgb{112, 214, 255},
	Bool:            &rgb{249, 87, 56},
	Number:          &rgb{0, 168, 232},
	VarType:         &rgb{0, 150, 199},
	Pointer:         &rgb{78, 205, 196}, // TODO: get rid of this
	Func:            &rgb{160, 90, 220},
	StructField:     &rgb{211, 211, 211},
	StructFieldHash: &rgb{255, 123, 0},
}

// DisableColors disables the colors globally.
func DisableColors() { // TODO: deprecate this function
	defaultTheme = theme{}
}
