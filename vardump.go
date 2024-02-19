package vardump

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	color_str            = lipgloss.NewStyle().Foreground(lipgloss.Color("#8ac926")).Render
	color_str_quotes     = lipgloss.NewStyle().Foreground(lipgloss.Color("#70d6ff")).Render
	color_bool           = lipgloss.NewStyle().Foreground(lipgloss.Color("#f95738")).Render
	color_number         = lipgloss.NewStyle().Foreground(lipgloss.Color("#00a8e8")).Render
	color_var_type       = lipgloss.NewStyle().Foreground(lipgloss.Color("#0096c7")).Render
	color_ptr            = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ecdc4")).Render
	color_func           = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff7b00")).Render
	color_struct         = lipgloss.NewStyle().Foreground(lipgloss.Color("#d3d3d3")).Render
	struct_private_field = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff7b00")).Render("#")
)

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

type dumper struct {
	buf  bytes.Buffer
	deep int
}

func (d *dumper) dump(v any, ignore_deep ...bool) {
	if len(ignore_deep) <= 0 || !ignore_deep[0] {
		d.buf.WriteString(strings.Repeat("   ", d.deep))
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
	d.buf.WriteString(color_str_quotes(`"`))
	d.buf.WriteString(color_str(v))
	d.buf.WriteString(color_str_quotes(`"`))
}

func (d *dumper) dumpBool(v bool) {
	if v {
		d.buf.WriteString(color_bool(`true`))
	} else {
		d.buf.WriteString(color_bool(`false`))
	}
}

func (d *dumper) dumpNum(v any) {
	d.buf.WriteString(color_number(fmt.Sprint(v)))
}

func (d *dumper) dumpSlice(v any) {
	value := reflect.ValueOf(v)
	length := value.Len()
	capacity := value.Cap()

	d.buf.WriteString(color_var_type(fmt.Sprintf("%T:%d:%d {", v, length, capacity)))

	d.deep++
	for i := 0; i < length; i++ {
		d.buf.WriteByte(0xa)
		d.dump(value.Index(i).Interface())
		d.buf.WriteString((","))
	}
	d.deep--
	d.buf.WriteString("\n" + strings.Repeat("   ", d.deep) + color_var_type("}"))
}

func (d *dumper) dumpMap(v any) {
	value := reflect.ValueOf(v)
	keys := value.MapKeys()

	d.buf.WriteString(color_var_type(fmt.Sprintf("%T:%d {", v, len(keys))))

	d.deep++
	for _, key := range keys {
		d.buf.WriteByte(0xa)
		d.dump(key.Interface())
		d.buf.WriteString((": "))
		d.dump(value.MapIndex(key).Interface(), true)
		d.buf.WriteString((","))
	}
	d.deep--

	d.buf.WriteString("\n" + strings.Repeat("   ", d.deep) + color_var_type("}"))
}

func (d *dumper) dumpPointer(v any) {
	d.buf.WriteString(color_ptr(fmt.Sprintf("%T:%p", v, v)))
}

func (d *dumper) dumpFunc(fn string) {
	d.buf.WriteString(color_func(fn))
}

func (d *dumper) dumpStruct(v any) {
	typ := fmt.Sprintf("%T", v)
	if strings.HasPrefix(typ, "struct") {
		typ = "struct"
	}
	d.buf.WriteString(color_var_type(typ + " {"))

	def := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	d.deep++
	for i := 0; i < def.NumField(); i++ {
		d.buf.WriteByte(0xa)
		k := def.Field(i)
		d.dumpStructKey(k)
		d.buf.WriteString((": "))

		if !k.IsExported() {
			d.buf.WriteString(color_var_type(fmt.Sprintf("%v", k.Type)))
			continue
		}

		d.dump(value.Field(i).Interface(), true)

		d.buf.WriteString((","))
	}
	d.deep--

	d.buf.WriteString("\n" + strings.Repeat("   ", d.deep) + color_var_type("}"))
}

func (d *dumper) dumpErr(err error) {
	d.buf.WriteString(color_var_type(fmt.Sprintf("%T {", err)))
	d.buf.WriteString("\n")
	d.buf.WriteString(strings.Repeat("   ", d.deep+1))
	d.buf.WriteString(color_str_quotes(`"`))
	d.buf.WriteString(color_str(err.Error()))
	d.buf.WriteString(color_str_quotes(`"`))
	d.buf.WriteString("\n" + strings.Repeat("   ", d.deep) + color_var_type("}"))
}

func (d *dumper) dumpStructKey(key reflect.StructField) {
	d.buf.WriteString(strings.Repeat("   ", d.deep))
	if !key.IsExported() {
		d.buf.WriteString(struct_private_field)
	}
	d.buf.WriteString(color_struct(key.Name))
}
