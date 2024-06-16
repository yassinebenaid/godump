package godump

// Dump the given variable
func Dump(v any) error {
	d := Dumper{}
	d.theme = defaultTheme
	return d.Println(v)
}

// DumpNC is just like Dump but doesn't produce any colors , useful if you want to write to a file or stream.
func DumpNC(v any) error {
	return (&Dumper{}).Println(v)
}

// Sdump is just like Dump but returns the result instead of printing to STDOUT
func Sdump(v any) string {
	d := Dumper{}
	d.theme = defaultTheme
	return d.Sprint(v)
}

// SdumpNC is just like DumpNC but returns the result instead of printing to STDOUT
func SdumpNC(v any) string {
	return (&Dumper{}).Sprint(v)
}
