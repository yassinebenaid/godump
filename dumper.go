package godump

import (
	"fmt"
	"reflect"
	"strings"
)

type dumper struct {
	buf               []byte
	indentation       string
	dumpPrivateFields bool
	theme             theme
	depth             uint
	ptrs              map[uintptr]uint
	ptrTag            uint
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
		d.write(d.theme.Chan.apply(val.Type().String()))
		if cap := val.Cap(); cap > 0 {
			d.write(d.theme.Chan.apply(fmt.Sprintf("<%d>", cap)))
		}
	case reflect.Struct:
		d.dumpStruct(val)
	case reflect.Pointer:
		d.dumpPointer(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d.write(d.theme.Number.apply(fmt.Sprint(val)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		d.write(d.theme.Number.apply(fmt.Sprint(val)))
	case reflect.Float32, reflect.Float64:
		d.write(d.theme.Number.apply(fmt.Sprint(val)))
	case reflect.Complex64, reflect.Complex128:
		d.write(d.theme.Number.apply(fmt.Sprint(val)))
	case reflect.Uintptr:
		d.write(d.theme.Number.apply(fmt.Sprintf("0x%x", val.Uint())))
	case reflect.Invalid:
		d.write(d.theme.Nil.apply("nil"))
	case reflect.Interface:
		d.dump(val.Elem(), true)
	case reflect.UnsafePointer:
		d.write(d.theme.VarType.apply(fmt.Sprintf("unsafe.Pointer(%v)", val.UnsafePointer())))
	}
}

func (d *dumper) dumpSlice(v reflect.Value) {
	length := v.Len()

	var tag string
	if d.ptrTag != 0 {
		tag = d.theme.PointerCounter.apply(fmt.Sprintf("#%d", d.ptrTag))
		d.ptrTag = 0
	}

	d.write(d.theme.VarType.apply(fmt.Sprintf("%s:%d:%d {%s", v.Type(), length, v.Cap(), tag)))

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

	var tag string
	if d.ptrTag != 0 {
		tag = d.theme.PointerCounter.apply(fmt.Sprintf("#%d", d.ptrTag))
		d.ptrTag = 0
	}

	d.write(d.theme.VarType.apply(fmt.Sprintf("%s:%d {%s", v.Type(), len(keys), tag)))

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
		d.ptrs = make(map[uintptr]uint)
	}

	elem := v.Elem()

	if isPrimitive(elem) {
		d.write(d.theme.PointerSign.apply("&"))
		d.dump(elem, true)
		return
	}

	addr := uintptr(v.UnsafePointer())

	if id, ok := d.ptrs[addr]; ok {
		d.write(d.theme.PointerSign.apply("&"))
		d.write(d.theme.PointerCounter.apply(fmt.Sprintf("@%d", id)))
		return
	}

	d.ptrs[addr] = uint(len(d.ptrs) + 1)

	d.ptrTag = uint(len(d.ptrs))
	d.write(d.theme.PointerSign.apply("&"))
	d.dump(elem, true)
	d.ptrTag = 0
}

func (d *dumper) dumpStruct(v reflect.Value) {
	vtype, numFields := v.Type(), v.NumField()

	var tag string
	if d.ptrTag != 0 {
		tag = fmt.Sprintf("#%d", d.ptrTag)
		d.ptrTag = 0
	}

	if t := vtype.String(); strings.HasPrefix(t, "struct") {
		d.write(d.theme.VarType.apply("struct {"))
	} else {
		d.write(d.theme.VarType.apply(t + " {"))
	}
	d.write(d.theme.PointerCounter.apply(tag))

	d.depth++
	for i := 0; i < numFields; i++ {
		d.write("\n")
		d.indent()

		key := vtype.Field(i)
		d.write(d.theme.StructField.apply(key.Name))
		d.write((": "))

		if !key.IsExported() && !d.dumpPrivateFields {
			d.write(d.theme.VarType.apply(key.Type.String()))
		} else {
			d.dump(v.Field(i), true)
		}

		d.write((","))
	}
	d.depth--

	if numFields > 0 {
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

	d.write(strings.Repeat(d.indentation, int(d.depth)))
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
			fmt.Println(v.Kind())

			return false
		}
	}
}
