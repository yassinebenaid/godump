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
}

func (d *dumper) dump(v any, ignore_depth ...bool) {
	if len(ignore_depth) <= 0 || !ignore_depth[0] {
		d.buf.WriteString(strings.Repeat("   ", d.depth))
	}

	var_kind := reflect.ValueOf(v).Kind()

	switch {
	case var_kind == reflect.String:
		d.dumpString(fmt.Sprint(v))
	case var_kind == reflect.Bool:
		d.buf.WriteString(d.theme.Bool.apply(fmt.Sprintf("%t", v)))
	case strings.HasPrefix(var_kind.String(), "int") ||
		strings.HasPrefix(var_kind.String(), "uint") ||
		strings.HasPrefix(var_kind.String(), "float") ||
		strings.HasPrefix(var_kind.String(), "complex"):
		d.buf.WriteString(d.theme.Number.apply(fmt.Sprint(v)))
	case var_kind == reflect.Slice || var_kind == reflect.Array:
		d.dumpSlice(v)
	case var_kind == reflect.Map:
		d.dumpMap(v)
	case var_kind == reflect.Func:
		d.buf.WriteString(d.theme.Func.apply(fmt.Sprintf("%T", v)))
	case var_kind == reflect.Struct:
		d.dumpStruct(v)
	case var_kind == reflect.Pointer:
		d.dumpPointer(v)
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
	d.buf.WriteString(d.theme.Pointer.apply(fmt.Sprintf("%T:%p", v, v)))
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
