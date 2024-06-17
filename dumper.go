package godump

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

type Dumper struct {
	Indentation       string
	HidePrivateFields bool
	Theme             Theme
	buf               bytes.Buffer
	depth             uint
	ptrs              map[uintptr]uint
	ptrTag            uint
}

func (d *Dumper) Print(v any) error {
	d.init()
	d.dump(reflect.ValueOf(v))
	if _, err := d.buf.WriteTo(os.Stdout); err != nil {
		return fmt.Errorf("dumper error: encountered unexpected error while writing to STDOUT, %v", err)
	}
	return nil
}

func (d *Dumper) Println(v any) error {
	d.init()
	d.dump(reflect.ValueOf(v))
	d.buf.WriteString("\n")
	if _, err := d.buf.WriteTo(os.Stdout); err != nil {
		return fmt.Errorf("dumper error: encountered unexpected error while writing to STDOUT, %v", err)
	}
	return nil
}

func (d *Dumper) Fprint(dst io.Writer, v any) error {
	d.init()
	d.dump(reflect.ValueOf(v))
	if _, err := d.buf.WriteTo(dst); err != nil {
		return fmt.Errorf("dumper error: encountered unexpected error while writing to dst, %v", err)
	}
	return nil
}

func (d *Dumper) Fprintln(dst io.Writer, v any) error {
	d.init()
	d.dump(reflect.ValueOf(v))
	d.buf.WriteString("\n")
	if _, err := d.buf.WriteTo(dst); err != nil {
		return fmt.Errorf("dumper error: encountered unexpected error while writing to dst, %v", err)
	}
	return nil
}

func (d *Dumper) Sprint(v any) string {
	d.init()
	d.dump(reflect.ValueOf(v))
	return d.buf.String()
}

func (d *Dumper) Sprintln(v any) string {
	d.init()
	d.dump(reflect.ValueOf(v))
	d.buf.WriteString("\n")
	return d.buf.String()
}

func (d *Dumper) init() {
	d.buf.Reset()
	d.ptrs = make(map[uintptr]uint)
	if d.Indentation == "" {
		d.Indentation = "   "
	}
}

func (d *Dumper) dump(val reflect.Value, ignore_depth ...bool) {
	if len(ignore_depth) <= 0 || !ignore_depth[0] {
		d.indent()
	}

	switch val.Kind() {
	case reflect.String:
		d.buf.WriteString(__(d.Theme.Quotes, `"`) +
			__(d.Theme.String, val.String()) +
			__(d.Theme.Quotes, `"`))
	case reflect.Bool:
		d.buf.WriteString(__(d.Theme.Bool, fmt.Sprintf("%t", val.Bool())))
	case reflect.Slice, reflect.Array:
		d.dumpSlice(val)
	case reflect.Map:
		d.dumpMap(val)
	case reflect.Func:
		d.buf.WriteString(__(d.Theme.Func, val.Type().String()))
	case reflect.Chan:
		d.buf.WriteString(__(d.Theme.Chan, val.Type().String()))
		if cap := val.Cap(); cap > 0 {
			d.buf.WriteString(__(d.Theme.Chan, fmt.Sprintf("<%d>", cap)))
		}
	case reflect.Struct:
		d.dumpStruct(val)
	case reflect.Pointer:
		d.dumpPointer(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d.buf.WriteString(__(d.Theme.Number, fmt.Sprint(val)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		d.buf.WriteString(__(d.Theme.Number, fmt.Sprint(val)))
	case reflect.Float32, reflect.Float64:
		d.buf.WriteString(__(d.Theme.Number, fmt.Sprint(val)))
	case reflect.Complex64, reflect.Complex128:
		d.buf.WriteString(__(d.Theme.Number, fmt.Sprint(val)))
	case reflect.Uintptr:
		d.buf.WriteString(__(d.Theme.Number, fmt.Sprintf("0x%x", val.Uint())))
	case reflect.Invalid:
		d.buf.WriteString(__(d.Theme.Nil, "nil"))
	case reflect.Interface:
		d.dump(val.Elem(), true)
	case reflect.UnsafePointer:
		d.buf.WriteString(__(d.Theme.UnsafePointer, fmt.Sprintf("unsafe.Pointer(0x%x)", uintptr(val.UnsafePointer()))))
	}
}

func (d *Dumper) dumpSlice(v reflect.Value) {
	length := v.Len()

	var tag string
	if d.ptrTag != 0 {
		tag = __(d.Theme.PointerCounter, fmt.Sprintf("#%d", d.ptrTag))
		d.ptrTag = 0
	}

	d.buf.WriteString(__(d.Theme.Types, fmt.Sprintf("%s:%d:%d", v.Type(), length, v.Cap())))
	d.buf.WriteString(__(d.Theme.Braces, fmt.Sprintf(" {%s", tag)))

	d.depth++
	for i := 0; i < length; i++ {
		d.buf.WriteString("\n")
		d.dump(v.Index(i))
		d.buf.WriteString(",")
	}
	d.depth--

	if length > 0 {
		d.buf.WriteString("\n")
		d.indent()
	}

	d.buf.WriteString(__(d.Theme.Braces, "}"))
}

func (d *Dumper) dumpMap(v reflect.Value) {
	keys := v.MapKeys()

	var tag string
	if d.ptrTag != 0 {
		tag = __(d.Theme.PointerCounter, fmt.Sprintf("#%d", d.ptrTag))
		d.ptrTag = 0
	}

	d.buf.WriteString(__(d.Theme.Types, fmt.Sprintf("%s:%d", v.Type(), len(keys))))
	d.buf.WriteString(__(d.Theme.Braces, fmt.Sprintf(" {%s", tag)))

	d.depth++
	for _, key := range keys {
		d.buf.WriteString("\n")
		d.dump(key)
		d.buf.WriteString((": "))
		d.dump(v.MapIndex(key), true)
		d.buf.WriteString((","))
	}
	d.depth--

	if len(keys) > 0 {
		d.buf.WriteString("\n")
		d.indent()
	}

	d.buf.WriteString(__(d.Theme.Braces, "}"))
}

func (d *Dumper) dumpPointer(v reflect.Value) {
	elem := v.Elem()

	if isPrimitive(elem) {
		if elem.IsValid() {
			d.buf.WriteString(__(d.Theme.PointerSign, "&"))
		}
		d.dump(elem, true)
		return
	}

	addr := uintptr(v.UnsafePointer())

	if id, ok := d.ptrs[addr]; ok {
		d.buf.WriteString(__(d.Theme.PointerSign, "&"))
		d.buf.WriteString(__(d.Theme.PointerCounter, fmt.Sprintf("@%d", id)))
		return
	}

	d.ptrs[addr] = uint(len(d.ptrs) + 1)

	d.ptrTag = uint(len(d.ptrs))
	d.buf.WriteString(__(d.Theme.PointerSign, "&"))
	d.dump(elem, true)
	d.ptrTag = 0
}

func (d *Dumper) dumpStruct(v reflect.Value) {
	vtype := v.Type()

	var tag string
	if d.ptrTag != 0 {
		tag = fmt.Sprintf("#%d", d.ptrTag)
		d.ptrTag = 0
	}

	if t := vtype.String(); strings.HasPrefix(t, "struct") {
		d.buf.WriteString(__(d.Theme.Types, "struct"))
	} else {
		d.buf.WriteString(__(d.Theme.Types, t))
	}
	d.buf.WriteString(__(d.Theme.Braces, " {"))
	d.buf.WriteString(__(d.Theme.PointerCounter, tag))

	var has_fields bool

	d.depth++
	for i := 0; i < v.NumField(); i++ {
		key := vtype.Field(i)
		if !key.IsExported() && d.HidePrivateFields {
			continue
		}

		has_fields = true

		d.buf.WriteString("\n")
		d.indent()

		d.buf.WriteString(__(d.Theme.StructField, key.Name))
		d.buf.WriteString((": "))
		d.dump(v.Field(i), true)
		d.buf.WriteString((","))
	}
	d.depth--

	if has_fields {
		d.buf.WriteString("\n")
		d.indent()
	}

	d.buf.WriteString(__(d.Theme.Braces, "}"))
}

func (d *Dumper) indent() {
	d.buf.WriteString(strings.Repeat(d.Indentation, int(d.depth)))
}

func isPrimitive(val reflect.Value) bool {
	v := val
	for {
		switch v.Kind() {
		case reflect.String, reflect.Bool, reflect.Func, reflect.Chan,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Invalid, reflect.UnsafePointer:
			return true
		case reflect.Pointer:
			v = v.Elem()
		default:
			return false
		}
	}
}
