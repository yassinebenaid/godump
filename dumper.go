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
	buf               bytes.Buffer
	dumpPrivateFields bool
	theme             theme
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
		d.buf.WriteString(d.theme.Quotes.__(`"`) +
			d.theme.String.__(val.String()) +
			d.theme.Quotes.__(`"`))
	case reflect.Bool:
		d.buf.WriteString(d.theme.Bool.__(fmt.Sprintf("%t", val.Bool())))
	case reflect.Slice, reflect.Array:
		d.dumpSlice(val)
	case reflect.Map:
		d.dumpMap(val)
	case reflect.Func:
		d.buf.WriteString(d.theme.Func.__(val.Type().String()))
	case reflect.Chan:
		d.buf.WriteString(d.theme.Chan.__(val.Type().String()))
		if cap := val.Cap(); cap > 0 {
			d.buf.WriteString(d.theme.Chan.__(fmt.Sprintf("<%d>", cap)))
		}
	case reflect.Struct:
		d.dumpStruct(val)
	case reflect.Pointer:
		d.dumpPointer(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d.buf.WriteString(d.theme.Number.__(fmt.Sprint(val)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		d.buf.WriteString(d.theme.Number.__(fmt.Sprint(val)))
	case reflect.Float32, reflect.Float64:
		d.buf.WriteString(d.theme.Number.__(fmt.Sprint(val)))
	case reflect.Complex64, reflect.Complex128:
		d.buf.WriteString(d.theme.Number.__(fmt.Sprint(val)))
	case reflect.Uintptr:
		d.buf.WriteString(d.theme.Number.__(fmt.Sprintf("0x%x", val.Uint())))
	case reflect.Invalid:
		d.buf.WriteString(d.theme.Nil.__("nil"))
	case reflect.Interface:
		d.dump(val.Elem(), true)
	case reflect.UnsafePointer:
		d.buf.WriteString(d.theme.UnsafePointer.__(fmt.Sprintf("unsafe.Pointer(0x%x)", uintptr(val.UnsafePointer()))))
	}
}

func (d *Dumper) dumpSlice(v reflect.Value) {
	length := v.Len()

	var tag string
	if d.ptrTag != 0 {
		tag = d.theme.PointerCounter.__(fmt.Sprintf("#%d", d.ptrTag))
		d.ptrTag = 0
	}

	d.buf.WriteString(d.theme.Types.__(fmt.Sprintf("%s:%d:%d", v.Type(), length, v.Cap())))
	d.buf.WriteString(d.theme.Braces.__(fmt.Sprintf(" {%s", tag)))

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

	d.buf.WriteString(d.theme.Braces.__("}"))
}

func (d *Dumper) dumpMap(v reflect.Value) {
	keys := v.MapKeys()

	var tag string
	if d.ptrTag != 0 {
		tag = d.theme.PointerCounter.__(fmt.Sprintf("#%d", d.ptrTag))
		d.ptrTag = 0
	}

	d.buf.WriteString(d.theme.Types.__(fmt.Sprintf("%s:%d", v.Type(), len(keys))))
	d.buf.WriteString(d.theme.Braces.__(fmt.Sprintf(" {%s", tag)))

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

	d.buf.WriteString(d.theme.Braces.__("}"))
}

func (d *Dumper) dumpPointer(v reflect.Value) {
	elem := v.Elem()

	if isPrimitive(elem) {
		if elem.IsValid() {
			d.buf.WriteString(d.theme.PointerSign.__("&"))
		}
		d.dump(elem, true)
		return
	}

	addr := uintptr(v.UnsafePointer())

	if id, ok := d.ptrs[addr]; ok {
		d.buf.WriteString(d.theme.PointerSign.__("&"))
		d.buf.WriteString(d.theme.PointerCounter.__(fmt.Sprintf("@%d", id)))
		return
	}

	d.ptrs[addr] = uint(len(d.ptrs) + 1)

	d.ptrTag = uint(len(d.ptrs))
	d.buf.WriteString(d.theme.PointerSign.__("&"))
	d.dump(elem, true)
	d.ptrTag = 0
}

func (d *Dumper) dumpStruct(v reflect.Value) {
	vtype, numFields := v.Type(), v.NumField()

	var tag string
	if d.ptrTag != 0 {
		tag = fmt.Sprintf("#%d", d.ptrTag)
		d.ptrTag = 0
	}

	if t := vtype.String(); strings.HasPrefix(t, "struct") {
		d.buf.WriteString(d.theme.Types.__("struct"))
	} else {
		d.buf.WriteString(d.theme.Types.__(t))
	}
	d.buf.WriteString(d.theme.Braces.__(" {"))
	d.buf.WriteString(d.theme.PointerCounter.__(tag))

	d.depth++
	for i := 0; i < numFields; i++ {
		d.buf.WriteString("\n")
		d.indent()

		key := vtype.Field(i)
		d.buf.WriteString(d.theme.StructField.__(key.Name))
		d.buf.WriteString((": "))

		if !key.IsExported() && !d.dumpPrivateFields {
			d.buf.WriteString(d.theme.Types.__(key.Type.String()))
		} else {
			d.dump(v.Field(i), true)
		}

		d.buf.WriteString((","))
	}
	d.depth--

	if numFields > 0 {
		d.buf.WriteString("\n")
		d.indent()
	}

	d.buf.WriteString(d.theme.Braces.__("}"))
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
