package godump

import (
	"fmt"
)

// Color defines a general interface used for styling.
type Color interface {
	colorize(string) string
}

func __(c Color, v string) string {
	if c == nil {
		return v
	}
	return c.colorize(v)
}

// RGB implements [Color] and allow you to define your style as an RGB value.
type RGB struct {
	R, G, B int
}

func (rgb RGB) colorize(v string) string {
	return fmt.Sprintf("\033[38;2;%v;%v;%vm%s\033[0m", rgb.R, rgb.G, rgb.B, v)
}

// Theme allows you to define your preferred styling for [Dumper].
type Theme struct {
	// String defines the style used for strings
	String Color

	// Quotes defines the style used for quotes (") around strings.
	Quotes Color

	// Bool defines the style used for boolean values.
	Bool Color

	// Number defines the style used for numbers, including all types of integers, floats and complex numbers.
	Number Color

	// StructuralTypes defines the style used for structural types, typically slices, structs and maps.
	StructuralTypes Color

	// Nil defines the style used for nil.
	Nil Color

	// Func defines the style used for functions.
	Func Color

	// Chan defines the style used for channels.
	Chan Color

	// UnsafePointer defines the style used for unsafe pointers.
	UnsafePointer Color

	// PointerSymbol defines the style used for pointer symbol '&'.
	PointerSymbol Color

	// PointerTag defines the style used for pointer tags, typically the pointer id '#x' and the recursive reference '@x'.
	PointerTag Color

	// Fields defines the style used for struct fields.
	Fields Color

	// Braces defines the style used for braces '{}' in structural types.
	Braces Color
}

// DefaultTheme is the default [Theme] used by [Dump].
var DefaultTheme = Theme{
	String:          RGB{138, 201, 38},
	Quotes:          RGB{112, 214, 255},
	Bool:            RGB{249, 87, 56},
	Number:          RGB{10, 178, 242},
	StructuralTypes: RGB{0, 150, 199},
	PointerSymbol:   RGB{205, 93, 0},
	PointerTag:      RGB{110, 110, 110},
	Nil:             RGB{219, 57, 26},
	Func:            RGB{160, 90, 220},
	Fields:          RGB{189, 176, 194},
	Chan:            RGB{195, 154, 76},
	UnsafePointer:   RGB{89, 193, 180},
	Braces:          RGB{185, 86, 86},
}

// DisableColors disables the colors globally.
//
// Deprecated: As of v0.8.0 this function only sets the [DefaultTheme] to a zero value
func DisableColors() { // TODO: deprecate this function
	DefaultTheme = Theme{}
}
