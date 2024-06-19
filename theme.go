package godump

import (
	"fmt"
)

// Style defines a general interface used for styling.
type Style interface {
	apply(string) string
}

func __(s Style, v string) string {
	if s == nil {
		return v
	}
	return s.apply(v)
}

// RGB implements [Style] and allow you to define your style as an RGB value.
type RGB struct {
	R, G, B int
}

func (rgb RGB) apply(v string) string {
	return fmt.Sprintf("\033[38;2;%v;%v;%vm%s\033[0m", rgb.R, rgb.G, rgb.B, v)
}

// Theme allows you to define your preferred styling for [Dumper].
type Theme struct {
	// String defines the style used for strings
	String Style

	// Quotes defines the style used for quotes (") around strings.
	Quotes Style

	// Bool defines the style used for boolean values.
	Bool Style

	// Number defines the style used for numbers, including all types of integers, floats and complex numbers.
	Number Style

	// Types defines the style used for defined and/or structural types, eg. slices, structs, maps...
	Types Style

	// Nil defines the style used for nil.
	Nil Style

	// Func defines the style used for functions.
	Func Style

	// Chan defines the style used for channels.
	Chan Style

	// UnsafePointer defines the style used for unsafe pointers.
	UnsafePointer Style

	// Address defines the style used for address symbol '&'.
	Address Style

	// PointerTag defines the style used for pointer tags, typically the pointer id '#x' and the recursive reference '@x'.
	PointerTag Style

	// Fields defines the style used for struct fields.
	Fields Style

	// Braces defines the style used for braces '{}' in structural types.
	Braces Style
}

// DefaultTheme is the default [Theme] used by [Dump].
var DefaultTheme = Theme{
	String:        RGB{138, 201, 38},
	Quotes:        RGB{112, 214, 255},
	Bool:          RGB{249, 87, 56},
	Number:        RGB{10, 178, 242},
	Types:         RGB{0, 150, 199},
	Address:       RGB{205, 93, 0},
	PointerTag:    RGB{110, 110, 110},
	Nil:           RGB{219, 57, 26},
	Func:          RGB{160, 90, 220},
	Fields:        RGB{189, 176, 194},
	Chan:          RGB{195, 154, 76},
	UnsafePointer: RGB{89, 193, 180},
	Braces:        RGB{185, 86, 86},
}

// DisableColors disables the colors globally.
//
// Deprecated: As of v0.8.0 this function only sets the [DefaultTheme] to a zero value
func DisableColors() {
	DefaultTheme = Theme{}
}
