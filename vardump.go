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
	color_str        = lipgloss.NewStyle().Foreground(lipgloss.Color("#8ac926")).Render
	color_str_quotes = lipgloss.NewStyle().Foreground(lipgloss.Color("#70d6ff")).Render
	color_bool       = lipgloss.NewStyle().Foreground(lipgloss.Color("#f95738")).Render
	color_number     = lipgloss.NewStyle().Foreground(lipgloss.Color("#00a8e8")).Render
	color_var_type   = lipgloss.NewStyle().Foreground(lipgloss.Color("#1982c4")).Render
	color_ptr        = lipgloss.NewStyle().Foreground(lipgloss.Color("#7678ed")).Render
	color_func       = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff7b00")).Render
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

	var_t := fmt.Sprintf("%T", v)

	if strings.HasPrefix(var_t, "int") || strings.HasPrefix(var_t, "float") {
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
		d.dumpPointer(var_t)
		return
	}

	if strings.HasPrefix(var_t, "func") {
		d.dumpFunc(var_t)
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
	d.buf.WriteString(color_var_type(fmt.Sprintf("%T {", v)))

	value := reflect.ValueOf(v)

	d.deep++
	for i := 0; i < value.Len(); i++ {
		d.buf.WriteByte(0xa)
		d.dump(value.Index(i).Interface())
		d.buf.WriteString((","))
	}
	d.deep--
	d.buf.WriteString("\n" + strings.Repeat("   ", d.deep) + color_var_type("}"))
}

func (d *dumper) dumpMap(v any) {
	d.buf.WriteString(color_var_type(fmt.Sprintf("%T {", v)))

	value := reflect.ValueOf(v)
	keys := value.MapKeys()

	d.deep++
	for _, key := range keys {
		d.buf.WriteByte(0xa)
		d.dump(key.Interface())
		d.buf.WriteString((": "))

		v := value.MapIndex(key).Interface()
		if isprim(v) {
			deep := d.deep
			d.deep = 0
			d.dump(v)
			d.deep = deep
		} else {
			d.dump(v, true)
		}

		d.buf.WriteString((","))
	}
	d.deep--

	d.buf.WriteString("\n" + strings.Repeat("   ", d.deep) + color_var_type("}"))
}

func (d *dumper) dumpPointer(ptr string) {
	d.buf.WriteString(color_ptr(`&` + strings.TrimPrefix(ptr, "*")))
}

func (d *dumper) dumpFunc(fn string) {
	d.buf.WriteString(color_func(fn))
}

func isprim(v any) bool {
	t := fmt.Sprintf("%T", v)
	return t == "string" ||
		t == "bool" ||
		strings.HasPrefix(t, "int") ||
		strings.HasPrefix(t, "float")
}
