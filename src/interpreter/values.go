// Implementation of Value... interfaces which are defined in
// Go package "github.com/sbinet/go-eval".

package interpreter

import (
	"bytes"
	"github.com/sbinet/go-eval"
	"strconv"
	"strings"
)

type print_style_t int

const (
	SINGLE_LINE print_style_t = iota
	MULTI_LINE
)

// Used in String() methods
const indentation_step = "  "

type print_config_t struct {
	names_orNil    []string                  // The name/tag of each value
	printers_orNil []func(eval.Value) string // Custom printers. The array may contain 'nil' elements.
	printStyle     print_style_t
}

// Print 'values' to the byte-buffer
func printValues(buf *bytes.Buffer, values []eval.Value, conf print_config_t) {
	buf.WriteString("{")
	indentation := ""
	if (conf.printStyle == MULTI_LINE) && (len(values) > 0) {
		buf.WriteString("\n")
		indentation = indentation_step
	}

	for i, v := range values {
		var s string
		if (conf.printers_orNil == nil) || (conf.printers_orNil[i] == nil) {
			s = v.String()
		} else {
			s = conf.printers_orNil[i](v)
		}

		if conf.printStyle == SINGLE_LINE {
			if conf.names_orNil != nil {
				buf.WriteString(conf.names_orNil[i])
				buf.WriteString(": ")
			}
			buf.WriteString(s)
		} else {
			ss := strings.Split(s, "\n")
			for j, x := range ss {
				buf.WriteString(indentation)
				if (j == 0) && (conf.names_orNil != nil) {
					buf.WriteString(conf.names_orNil[i])
					buf.WriteString(": ")
				}
				buf.WriteString(x)
				if j+1 < len(ss) {
					buf.WriteString("\n")
				}
			}

		}

		if i+1 < len(values) {
			buf.WriteString(",")
			if conf.printStyle == SINGLE_LINE {
				buf.WriteString(" ")
			}
		}
		buf.WriteString("\n")
	}

	buf.WriteString("}")
}

//
// Implementation of 'eval.StringValue'
//

type string_value_t string

func (v *string_value_t) String() string {
	return string(*v)
}

func (v *string_value_t) Assign(t *eval.Thread, o eval.Value) {
	*v = string_value_t(o.(eval.StringValue).Get(t))
}

func (v *string_value_t) Get(*eval.Thread) string {
	return string(*v)
}

func (v *string_value_t) Set(t *eval.Thread, x string) {
	*v = string_value_t(x)
}

//
// Implementation of 'eval.Array'
//

type array_value_t []eval.Value

func (v *array_value_t) String() string {
	var indexes []string
	for i := 0; i < len(*v); i++ {
		indexes = append(indexes, strconv.Itoa(i))
	}

	conf := print_config_t{
		names_orNil: indexes,
		printStyle:  MULTI_LINE,
	}

	var buf bytes.Buffer
	printValues(&buf, *v, conf)
	return buf.String()
}

func (v *array_value_t) Assign(t *eval.Thread, o eval.Value) {
	a := o.(eval.ArrayValue)
	l := len(*v)
	for i := 0; i < l; i++ {
		(*v)[i].Assign(t, a.Elem(t, int64(i)))
	}
}

func (v *array_value_t) Get(*eval.Thread) eval.ArrayValue {
	return v
}

func (v *array_value_t) Elem(t *eval.Thread, i int64) eval.Value {
	return (*v)[i]
}

func (v *array_value_t) Sub(i int64, length int64) eval.ArrayValue {
	if (i == 0) && (length == int64(len(*v))) {
		return v
	}

	res := (*v)[i : i+length]
	return &res
}

//
// Implementation of 'eval.SliceValue'
//

type slice_value_t struct {
	eval.Slice
}

func (v *slice_value_t) String() string {
	if v.Base == nil {
		return "<nil>"
	}
	return v.Base.Sub(0, v.Len).String()
}

func (v *slice_value_t) Assign(t *eval.Thread, o eval.Value) {
	v.Slice = o.(eval.SliceValue).Get(t)
}

func (v *slice_value_t) Get(*eval.Thread) eval.Slice {
	return v.Slice
}

func (v *slice_value_t) Set(t *eval.Thread, x eval.Slice) {
	v.Slice = x
}

//
// Implementation of 'eval.StructValue'
//

type struct_value_t struct {
	fields         []eval.Value
	names          []string
	printers_orNil []func(eval.Value) string // Custom per-field printers. A field printer can be nil.
	hide_orNil     []bool                    // Whether to print the corresponding field. Nil means "print all".
	printStyle     print_style_t
}

func (v *struct_value_t) String() string {
	var buf bytes.Buffer

	if v.hide_orNil == nil {
		conf := print_config_t{
			names_orNil:    v.names,
			printers_orNil: v.printers_orNil,
			printStyle:     v.printStyle,
		}

		printValues(&buf, v.fields, conf)
	} else {
		var visible_fields []eval.Value
		var visible_names []string
		var visible_printers []func(eval.Value) string

		for i, _ := range v.fields {
			if !v.hide_orNil[i] {
				visible_fields = append(visible_fields, v.fields[i])
				visible_names = append(visible_names, v.names[i])
				if v.printers_orNil != nil {
					visible_printers = append(visible_printers, v.printers_orNil[i])
				}
			}
		}

		conf := print_config_t{
			names_orNil:    visible_names,
			printers_orNil: visible_printers,
			printStyle:     v.printStyle,
		}

		printValues(&buf, visible_fields, conf)
	}

	return buf.String()
}

func (v *struct_value_t) Assign(t *eval.Thread, o eval.Value) {
	s := o.(eval.StructValue)
	l := len(v.fields)
	for i := 0; i < l; i++ {
		v.fields[i].Assign(t, s.Field(t, i))
	}
}

func (v *struct_value_t) Get(*eval.Thread) eval.StructValue {
	return v
}

func (v *struct_value_t) Field(t *eval.Thread, i int) eval.Value {
	return v.fields[i]
}
