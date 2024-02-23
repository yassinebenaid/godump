package godump

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

type dumper struct {
	buf   bytes.Buffer
	c     colorizer
	depth int
}

func (d *dumper) dump(v any, ignore_depth ...bool) {
	if len(ignore_depth) <= 0 || !ignore_depth[0] {
		d.buf.WriteString(strings.Repeat("   ", d.depth))
	}

	if s, ok := v.(string); ok {
		d.dumpString(s)
		return
	}

	if b, ok := v.(bool); ok {
		d.dumpBool(b)
		return
	}

	if b, ok := v.(error); ok {
		d.dumpErr(b)
		return
	}

	var_t := fmt.Sprintf("%T", v)

	if strings.HasPrefix(var_t, "int") || strings.HasPrefix(var_t, "float") || strings.HasPrefix(var_t, "complex") {
		d.dumpNum(v)
		return
	}

	if strings.HasPrefix(var_t, "[]") {
		d.dumpSlice(v)
		return
	}

	if strings.HasPrefix(var_t, "map") {
		d.dumpMap(v)
		return
	}

	if strings.HasPrefix(var_t, "*") {
		d.dumpPointer(v)
		return
	}

	if strings.HasPrefix(var_t, "func") {
		d.dumpFunc(var_t)
		return
	}

	if reflect.TypeOf(v).Kind() == reflect.Struct {
		d.dumpStruct(v)
		return
	}
}

func (d *dumper) dumpString(v string) {
	d.buf.WriteString(d.c.quote(`"`))
	d.buf.WriteString(d.c.str(v))
	d.buf.WriteString(d.c.quote(`"`))
}

func (d *dumper) dumpBool(v bool) {
	if v {
		d.buf.WriteString(d.c.bool(`true`))
	} else {
		d.buf.WriteString(d.c.bool(`false`))
	}
}

func (d *dumper) dumpNum(v any) {
	d.buf.WriteString(d.c.num(fmt.Sprint(v)))
}

func (d *dumper) dumpSlice(v any) {
	value := reflect.ValueOf(v)
	length := value.Len()
	capacity := value.Cap()

	d.buf.WriteString(d.c.vtype(fmt.Sprintf("%T:%d:%d {", v, length, capacity)))

	d.depth++
	for i := 0; i < length; i++ {
		d.buf.WriteByte(0xa)
		d.dump(value.Index(i).Interface())
		d.buf.WriteString((","))
	}
	d.depth--
	d.buf.WriteString("\n" + strings.Repeat("   ", d.depth) + d.c.vtype("}"))
}

func (d *dumper) dumpMap(v any) {
	value := reflect.ValueOf(v)
	keys := value.MapKeys()

	d.buf.WriteString(d.c.vtype(fmt.Sprintf("%T:%d {", v, len(keys))))

	d.depth++
	for _, key := range keys {
		d.buf.WriteByte(0xa)
		d.dump(key.Interface())
		d.buf.WriteString((": "))
		d.dump(value.MapIndex(key).Interface(), true)
		d.buf.WriteString((","))
	}
	d.depth--

	d.buf.WriteString("\n" + strings.Repeat("   ", d.depth) + d.c.vtype("}"))
}

func (d *dumper) dumpPointer(v any) {
	d.buf.WriteString(d.c.ptr(fmt.Sprintf("%T:%p", v, v)))
}

func (d *dumper) dumpFunc(fn string) {
	d.buf.WriteString(d.c.fn(fn))
}

func (d *dumper) dumpStruct(v any) {
	typ := fmt.Sprintf("%T", v)
	if strings.HasPrefix(typ, "struct") {
		typ = "struct"
	}
	d.buf.WriteString(d.c.vtype(typ + " {"))

	def := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	d.depth++
	for i := 0; i < def.NumField(); i++ {
		d.buf.WriteByte(0xa)
		k := def.Field(i)
		d.dumpStructKey(k)
		d.buf.WriteString((": "))

		if !k.IsExported() {
			d.buf.WriteString(d.c.vtype(fmt.Sprintf("%v", k.Type)))
			continue
		}

		d.dump(value.Field(i).Interface(), true)

		d.buf.WriteString((","))
	}
	d.depth--

	d.buf.WriteString("\n" + strings.Repeat("   ", d.depth) + d.c.vtype("}"))
}

func (d *dumper) dumpErr(err error) {
	d.buf.WriteString(d.c.vtype(fmt.Sprintf("%T {", err)))
	d.buf.WriteString("\n")
	d.buf.WriteString(strings.Repeat("   ", d.depth+1))
	d.buf.WriteString(d.c.quote(`"`))
	d.buf.WriteString(d.c.str(err.Error()))
	d.buf.WriteString(d.c.quote(`"`))
	d.buf.WriteString("\n" + strings.Repeat("   ", d.depth) + d.c.vtype("}"))
}

func (d *dumper) dumpStructKey(key reflect.StructField) {
	d.buf.WriteString(strings.Repeat("   ", d.depth))
	if !key.IsExported() {
		d.buf.WriteString(d.c.strct_prv("#"))
	}
	d.buf.WriteString(d.c.strct(key.Name))
}
