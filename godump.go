package godump

import "os"

// Dump the given variable
func Dump(v any) error {
	d := dumper{}
	d.dump(v)
	d.buf.WriteByte(0xa)
	_, err := d.buf.WriteTo(os.Stdout)
	if err != nil {
		return err
	}
	return nil
}

// DumpNC is just like Dump but doesn't produce any colors , useful if you want to write to a file or stream.
func DumpNC(v any) error {
	d := dumper{}
	d.c.disabled = true
	d.dump(v)
	d.buf.WriteByte(0xa)
	_, err := d.buf.WriteTo(os.Stdout)
	if err != nil {
		return err
	}
	return nil
}

// Sdump is just like Dump but returns the result instead of prining to STDOUT
func Sdump(v any) string {
	d := dumper{}
	d.dump(v)
	d.buf.WriteByte(0xa)
	return d.buf.String()
}

// Sdump is just like DumpNC but returns the result instead of prining to STDOUT
func SdumpNC(v any) string {
	d := dumper{}
	d.c.disabled = true
	d.dump(v)
	d.buf.WriteByte(0xa)
	return d.buf.String()
}
