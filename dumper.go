package godump

import (
	"fmt"
	"reflect"
	"strings"
)

type pointer struct {
	id     int
	pos    int
	tagged bool
}

type dumper struct {
	buf         []byte
	indentation string
	theme       theme
	depth       int
	ptrs        map[uintptr]*pointer
}

func (d *dumper) dump(val reflect.Value, ignore_depth ...bool) {
	if len(ignore_depth) <= 0 || !ignore_depth[0] {
		d.indent()
	}

	switch val.Kind() {
	case reflect.String:
		d.write(d.theme.Quotes.apply(`"`) +
			d.theme.String.apply(val.String()) +
			d.theme.Quotes.apply(`"`))
	case reflect.Bool:
		d.write(d.theme.Bool.apply(fmt.Sprintf("%t", val.Bool())))
	case reflect.Slice, reflect.Array:
		d.dumpSlice(val)
	case reflect.Map:
		d.dumpMap(val)
	case reflect.Func:
		d.write(d.theme.Func.apply(val.Type().String()))
	case reflect.Chan:
		d.write(d.theme.VarType.apply(val.Type().String()))
		if cap := val.Cap(); cap > 0 {
			d.write(d.theme.VarType.apply(fmt.Sprintf("<%d>", cap)))
		}
	case reflect.Struct:
		d.dumpStruct(val)
	case reflect.Pointer:
		d.dumpPointer(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d.write(d.theme.Number.apply(fmt.Sprint(val)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		d.write(d.theme.Number.apply(fmt.Sprint(val)))
	case reflect.Float32, reflect.Float64:
		d.write(d.theme.Number.apply(fmt.Sprint(val)))
	case reflect.Complex64, reflect.Complex128:
		d.write(d.theme.Number.apply(fmt.Sprint(val)))
	case reflect.Invalid:
		d.write(d.theme.Nil.apply("nil"))
	case reflect.Interface:
		d.dump(val.Elem(), true)
	}
}

func (d *dumper) dumpSlice(v reflect.Value) {
	length := v.Len()

	d.write(d.theme.VarType.apply(fmt.Sprintf("%s:%d:%d {", v.Type(), length, v.Cap())))

	d.depth++
	for i := 0; i < length; i++ {
		d.write("\n")
		d.dump(v.Index(i))
		d.write(",")
	}
	d.depth--

	if length > 0 {
		d.write("\n")
		d.indent()
	}

	d.write(d.theme.VarType.apply("}"))
}

func (d *dumper) dumpMap(v reflect.Value) {
	keys := v.MapKeys()
	d.write(d.theme.VarType.apply(fmt.Sprintf("%s:%d {", v.Type(), len(keys))))

	d.depth++
	for _, key := range keys {
		d.write("\n")
		d.dump(key)
		d.write((": "))
		d.dump(v.MapIndex(key), true)
		d.write((","))
	}
	d.depth--

	if len(keys) > 0 {
		d.write("\n")
		d.indent()
	}

	d.write(d.theme.VarType.apply("}"))
}

func (d *dumper) dumpPointer(v reflect.Value) {
	if d.ptrs == nil {
		d.ptrs = make(map[uintptr]*pointer)
	}

	addr := uintptr(v.UnsafePointer())

	if p, ok := d.ptrs[addr]; ok {
		d.write(d.theme.PointerSign.apply("&"))
		d.write(d.theme.PointerCounter.apply(fmt.Sprintf("@%d", p.id)))

		if !p.tagged {
			d.tagPtr(p)
			p.tagged = true
		}
		return
	}

	d.ptrs[addr] = &pointer{
		id:  len(d.ptrs) + 1,
		pos: len(d.buf),
	}

	d.write(d.theme.PointerSign.apply("&"))
	d.dump(v.Elem(), true)
}

func (d *dumper) tagPtr(ptr *pointer) {
	var shifted int

	for _, p := range d.ptrs {
		if ptr.pos > p.pos && p.tagged {
			shifted += len(d.theme.PointerCounter.apply(fmt.Sprintf("#%d", p.id)))
		}
	}

	nbuf := append([]byte{}, d.buf[:ptr.pos+shifted]...)
	nbuf = append(nbuf, []byte(d.theme.PointerCounter.apply(fmt.Sprintf("#%d", ptr.id)))...)
	nbuf = append(nbuf, d.buf[ptr.pos+shifted:]...)
	d.buf = nbuf
}

func (d *dumper) dumpStruct(v reflect.Value) {
	vtype := v.Type()

	if t := vtype.String(); t == "" {
		d.write(d.theme.VarType.apply("struct {"))
	} else {
		d.write(d.theme.VarType.apply(t + " {"))
	}

	d.depth++
	for i := 0; i < v.NumField(); i++ {
		d.write("\n")
		d.indent()

		key := vtype.Field(i)
		d.write(d.theme.StructField.apply(key.Name))
		d.write((": "))

		if !key.IsExported() {
			d.write(d.theme.VarType.apply(key.Type.String()))
		} else {
			d.dump(v.Field(i), true)
		}

		d.write((","))
	}
	d.depth--

	if v.NumField() > 0 {
		d.write("\n")
		d.indent()
	}

	d.write(d.theme.VarType.apply("}"))
}

func (d *dumper) write(s string) {
	d.buf = append(d.buf, []byte(s)...)
}

func (d *dumper) indent() {
	if d.indentation == "" {
		d.indentation = "   "
	}

	d.write(strings.Repeat(d.indentation, d.depth))
}
