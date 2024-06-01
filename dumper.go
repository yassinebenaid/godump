package godump

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

type dumper struct {
	buf   bytes.Buffer
	theme theme
	depth int
	ptrs  map[uintptr]uint
}

func (d *dumper) dump(v any, ignore_depth ...bool) {
	if len(ignore_depth) <= 0 || !ignore_depth[0] {
		d.buf.WriteString(strings.Repeat("   ", d.depth))
	}

	switch reflect.ValueOf(v).Kind() {
	case reflect.String:
		d.dumpString(fmt.Sprint(v))
	case reflect.Bool:
		d.buf.WriteString(d.theme.Bool.apply(fmt.Sprintf("%t", v)))
	case reflect.Slice, reflect.Array:
		d.dumpSlice(v)
	case reflect.Map:
		d.dumpMap(v)
	case reflect.Func:
		d.buf.WriteString(d.theme.Func.apply(fmt.Sprintf("%T", v)))
	case reflect.Chan:
		d.buf.WriteString(d.theme.VarType.apply(fmt.Sprintf("%T", v)))
		cap := reflect.ValueOf(v).Cap()
		if cap > 0 {
			d.buf.WriteString(d.theme.VarType.apply(fmt.Sprintf("<%d>", cap)))
		}
	case reflect.Struct:
		d.dumpStruct(v)
	case reflect.Pointer:
		d.dumpPointer(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = fmt.Sprint(reflect.ValueOf(v).Int())
		d.buf.WriteString(d.theme.Number.apply(v.(string)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = fmt.Sprint(reflect.ValueOf(v).Uint())
		d.buf.WriteString(d.theme.Number.apply(v.(string)))
	case reflect.Float32, reflect.Float64:
		v = fmt.Sprint(reflect.ValueOf(v).Float())
		d.buf.WriteString(d.theme.Number.apply(v.(string)))
	case reflect.Complex64, reflect.Complex128:
		v = fmt.Sprint(reflect.ValueOf(v).Complex())
		d.buf.WriteString(d.theme.Number.apply(v.(string)))
	case reflect.Invalid:
		d.buf.WriteString(d.theme.Nil.apply("nil"))

	}
}

func (d *dumper) dumpString(v string) {
	d.buf.WriteString(d.theme.Quotes.apply(`"`))
	d.buf.WriteString(d.theme.String.apply(v))
	d.buf.WriteString(d.theme.Quotes.apply(`"`))
}

func (d *dumper) dumpSlice(v any) {
	value := reflect.ValueOf(v)
	length := value.Len()
	capacity := value.Cap()

	d.buf.WriteString(d.theme.VarType.apply(fmt.Sprintf("%T:%d:%d {", v, length, capacity)))

	d.depth++
	for i := 0; i < length; i++ {
		d.buf.WriteByte(0xa)
		d.dump(value.Index(i).Interface())
		d.buf.WriteString((","))
	}
	d.depth--
	d.buf.WriteString("\n" + strings.Repeat("   ", d.depth) + d.theme.VarType.apply("}"))
}

func (d *dumper) dumpMap(v any) {
	value := reflect.ValueOf(v)
	keys := value.MapKeys()

	d.buf.WriteString(d.theme.VarType.apply(fmt.Sprintf("%T:%d {", v, len(keys))))

	d.depth++
	for _, key := range keys {
		d.buf.WriteByte(0xa)
		d.dump(key.Interface())
		d.buf.WriteString((": "))
		d.dump(value.MapIndex(key).Interface(), true)
		d.buf.WriteString((","))
	}
	d.depth--

	d.buf.WriteString("\n" + strings.Repeat("   ", d.depth) + d.theme.VarType.apply("}"))
}

func (d *dumper) dumpPointer(v any) {
	if d.ptrs == nil {
		d.ptrs = make(map[uintptr]uint)
	}

	ptr := uintptr(reflect.ValueOf(v).UnsafePointer())

	if ctr, ok := d.ptrs[ptr]; ok {
		d.buf.WriteString(d.theme.PointerSign.apply("&"))
		d.buf.WriteString(d.theme.PointerCounter.apply(fmt.Sprintf("#%d", ctr)))
		return
	}

	d.ptrs[ptr] = uint(len(d.ptrs) + 1)
	d.buf.WriteString(d.theme.PointerCounter.apply(fmt.Sprintf("#%d ", d.ptrs[ptr])))
	d.buf.WriteString(d.theme.PointerSign.apply("&"))

	actual := reflect.ValueOf(v).Elem()
	if actual.IsValid() {
		d.dump(actual.Interface(), true)
	} else {
		d.buf.WriteString(d.theme.Nil.apply("nil"))
	}

}

func (d *dumper) dumpStruct(v any) {
	typ := fmt.Sprintf("%T", v)
	if strings.HasPrefix(typ, "struct") {
		typ = "struct"
	}
	d.buf.WriteString(d.theme.VarType.apply(typ + " {"))

	def := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	d.depth++
	for i := 0; i < def.NumField(); i++ {
		d.buf.WriteByte(0xa)
		k := def.Field(i)
		d.dumpStructKey(k)
		d.buf.WriteString((": "))

		if !k.IsExported() {
			d.buf.WriteString(d.theme.VarType.apply(fmt.Sprintf("%v", k.Type)))
			continue
		}

		d.dump(value.Field(i).Interface(), true)

		d.buf.WriteString((","))
	}
	d.depth--

	d.buf.WriteString("\n" + strings.Repeat("   ", d.depth) + d.theme.VarType.apply("}"))
}

func (d *dumper) dumpStructKey(key reflect.StructField) {
	d.buf.WriteString(strings.Repeat("   ", d.depth))
	if !key.IsExported() {
		d.buf.WriteString(d.theme.StructFieldHash.apply("#"))
	}
	d.buf.WriteString(d.theme.StructField.apply(key.Name))
}
