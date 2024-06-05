package godump

import (
	"fmt"
	"os"
)

// Dump the given variable
func Dump(v any) error {
	d := dumper{}
	d.theme = defaultTheme
	d.dump(v)
	_, err := fmt.Fprintln(os.Stdout, string(d.buf))
	if err != nil {
		return err
	}
	return nil
}

// DumpNC is just like Dump but doesn't produce any colors , useful if you want to write to a file or stream.
func DumpNC(v any) error {
	d := dumper{}
	d.dump(v)
	_, err := fmt.Fprintln(os.Stdout, string(d.buf))
	if err != nil {
		return err
	}
	return nil
}

// Sdump is just like Dump but returns the result instead of printing to STDOUT
func Sdump(v any) string {
	d := dumper{}
	d.theme = defaultTheme
	d.dump(v)
	return string(d.buf)
}

// SdumpNC is just like DumpNC but returns the result instead of printing to STDOUT
func SdumpNC(v any) string {
	d := dumper{}
	d.dump(v)
	return string(d.buf)
}
