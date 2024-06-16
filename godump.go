package godump

import (
	"reflect"
)

// Dump the given variable
func Dump(v any) error {
	d := Dumper{}
	d.dumpPrivateFields = true
	d.theme = defaultTheme
	return d.Println(v)
}

// DumpNC is just like Dump but doesn't produce any colors , useful if you want to write to a file or stream.
func DumpNC(v any) error {
	d := Dumper{}
	d.dump(reflect.ValueOf(v))
	return d.Println(v)
}

// Sdump is just like Dump but returns the result instead of printing to STDOUT
func Sdump(v any) string {
	d := Dumper{}
	d.theme = defaultTheme
	d.dump(reflect.ValueOf(v))
	return d.buf.String()
}

// SdumpNC is just like DumpNC but returns the result instead of printing to STDOUT
func SdumpNC(v any) string {
	d := Dumper{}
	d.dump(reflect.ValueOf(v))
	return d.buf.String()
}
