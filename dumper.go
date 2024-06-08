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

func (d *dumper) dump(v any, ignore_depth ...bool) {
	if len(ignore_depth) <= 0 || !ignore_depth[0] {
		d.indent()
	}

	val := reflect.ValueOf(v)

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
		d.dumpMap(v)
	case reflect.Func:
		d.write(d.theme.Func.apply(val.Type().String()))
	case reflect.Chan:
		d.write(d.theme.VarType.apply(val.Type().String()))
		if cap := val.Cap(); cap > 0 {
			d.write(d.theme.VarType.apply(fmt.Sprintf("<%d>", cap)))
		}
	case reflect.Struct:
		d.dumpStruct(v)
	case reflect.Pointer:
		d.dumpPointer(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d.write(d.theme.Number.apply(fmt.Sprint(v)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		d.write(d.theme.Number.apply(fmt.Sprint(v)))
	case reflect.Float32, reflect.Float64:
		d.write(d.theme.Number.apply(fmt.Sprint(v)))
	case reflect.Complex64, reflect.Complex128:
		d.write(d.theme.Number.apply(fmt.Sprint(v)))
	case reflect.Invalid:
		d.write(d.theme.Nil.apply("nil"))
	}
}

func (d *dumper) dumpSlice(v reflect.Value) {
	length := v.Len()

	d.write(d.theme.VarType.apply(fmt.Sprintf("%s:%d:%d {", v.Type(), length, v.Cap())))

	d.depth++
	for i := 0; i < length; i++ {
		d.write("\n")
		d.dump(v.Index(i).Interface())
		d.write(",")
	}
	d.depth--

	if length > 0 {
		d.write("\n")
		d.indent()
	}

	d.write(d.theme.VarType.apply("}"))
}

func (d *dumper) dumpMap(v any) {
	value := reflect.ValueOf(v)
	keys := value.MapKeys()

	d.write(d.theme.VarType.apply(fmt.Sprintf("%T:%d {", v, len(keys))))

	d.depth++
	for _, key := range keys {
		d.write("\n")
		d.dump(key.Interface())
		d.write((": "))
		d.dump(value.MapIndex(key).Interface(), true)
		d.write((","))
	}
	d.depth--

	if len(keys) > 0 {
		d.write("\n")
		d.indent()
	}

	d.write(d.theme.VarType.apply("}"))
}

func (d *dumper) dumpPointer(v any) {
	if d.ptrs == nil {
		d.ptrs = make(map[uintptr]*pointer)
	}

	ptr := uintptr(reflect.ValueOf(v).UnsafePointer())

	if p, ok := d.ptrs[ptr]; ok {
		d.write(d.theme.PointerSign.apply("&"))
		d.write(d.theme.PointerCounter.apply(fmt.Sprintf("@%d", p.id)))

		if !p.tagged {
			d.tagPtr(p)
			p.tagged = true
		}
		return
	}

	d.ptrs[ptr] = &pointer{
		id:  len(d.ptrs) + 1,
		pos: len(d.buf),
	}

	d.write(d.theme.PointerSign.apply("&"))

	actual := reflect.ValueOf(v).Elem()
	if actual.IsValid() {
		d.dump(actual.Interface(), true)
	} else {
		d.write(d.theme.Nil.apply("nil"))
	}

}

func (d *dumper) dumpStruct(v any) {
	typ := fmt.Sprintf("%T", v)
	if strings.HasPrefix(typ, "struct") {
		typ = "struct"
	}
	d.write(d.theme.VarType.apply(typ + " {"))

	def := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	d.depth++
	for i := 0; i < def.NumField(); i++ {
		d.write("\n")
		k := def.Field(i)
		d.dumpStructKey(k)
		d.write((": "))

		if !k.IsExported() {
			d.write(d.theme.VarType.apply(fmt.Sprintf("%v", k.Type)))
			continue
		}

		d.dump(value.Field(i).Interface(), true)

		d.write((","))
	}
	d.depth--

	if def.NumField() > 0 {
		d.write("\n")
		d.indent()
	}

	d.write(d.theme.VarType.apply("}"))
}

func (d *dumper) dumpStructKey(key reflect.StructField) {
	d.indent()
	if !key.IsExported() {
		d.write(d.theme.StructFieldHash.apply("#"))
	}
	d.write(d.theme.StructField.apply(key.Name))
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

func (d *dumper) write(s string) {
	d.buf = append(d.buf, []byte(s)...)
}

func (d *dumper) indent() {
	if d.indentation == "" {
		d.indentation = "   "
	}

	d.write(strings.Repeat(d.indentation, d.depth))
}
