package pp

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

func Print(v interface{}) {
	Fprint(os.Stdout, v)
}

func Println(v interface{}) {
	Fprintln(os.Stdout, v)
}

func Sprint(v interface{}) string {
	var sb strings.Builder
	Fprint(&sb, v)
	return sb.String()
}

func Sprintln(v interface{}) string {
	var sb strings.Builder
	Fprintln(&sb, v)
	return sb.String()
}

func Fprint(w io.Writer, v interface{}) {
	if v == nil {
		io.WriteString(w, "nil")
	}
	fprint(w, reflect.ValueOf(v))
}

func Fprintln(w io.Writer, v interface{}) {
	Fprint(w, v)
	io.WriteString(w, "\n")
}

func PrintIndent(v interface{}, indent string) {
	FprintIndent(os.Stdout, v, indent)
}

func PrintIndentln(v interface{}, indent string) {
	FprintIndentln(os.Stdout, v, indent)
}

func SprintIndent(v interface{}, indent string) string {
	var sb strings.Builder
	FprintIndent(&sb, v, indent)
	return sb.String()
}

func SprintIndentln(v interface{}, indent string) string {
	var sb strings.Builder
	FprintIndentln(&sb, v, indent)
	return sb.String()
}

func FprintIndent(w io.Writer, v interface{}, indent string) {
	if v == nil {
		io.WriteString(w, "nil")
	}
	fprintIndent(w, reflect.ValueOf(v), indent, 0)
}

func FprintIndentln(w io.Writer, v interface{}, indent string) {
	FprintIndent(w, v, indent)
	io.WriteString(w, "\n")
}

func fprint(w io.Writer, rv reflect.Value) {
	switch v := rv.Interface().(type) {
	case string:
		fmt.Fprintf(w, "\"%s\"", v)
	case []byte:
		fmt.Fprintf(w, "\"%x\"", v)
	case fmt.Stringer:
		fmt.Fprintf(w, "\"%s\"", v.String())
	default:
		switch rv.Kind() {
		case reflect.Ptr:
			if rv.IsNil() {
				io.WriteString(w, "nil")
			} else {
				io.WriteString(w, "&")
				fprint(w, rv.Elem())
			}
		case reflect.Interface:
			if rv.IsNil() {
				io.WriteString(w, "nil")
			} else {
				fprint(w, rv.Elem())
			}
		case reflect.Array, reflect.Slice:
			io.WriteString(w, "[")
			for i := 0; i < rv.Len()-1; i++ {
				fprint(w, rv.Index(i))
				io.WriteString(w, ",")
			}
			if rv.Len() > 0 {
				fprint(w, rv.Index(rv.Len()-1))
			}
			io.WriteString(w, "]")
		case reflect.Map:
			isFirst := true
			iter := rv.MapRange()
			io.WriteString(w, "{")
			for iter.Next() {
				if isFirst {
					isFirst = false
				} else {
					io.WriteString(w, ",")
				}
				fprint(w, iter.Key())
				io.WriteString(w, ":")
				fprint(w, iter.Value())
			}
			io.WriteString(w, "}")
		case reflect.Struct:
			t := rv.Type()
			isFirst := true
			fmt.Fprintf(w, "%s{", t)
			for i := 0; i < rv.NumField(); i++ {
				f := rv.Field(i)
				if k := f.Kind(); (k == reflect.Chan || k == reflect.Func || k == reflect.Map || k == reflect.Pointer || k == reflect.UnsafePointer || k == reflect.Interface || k == reflect.Slice) && f.IsNil() {
					continue
				}
				if !f.CanInterface() {
					if !f.CanAddr() {
						continue
					}
					f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
				}
				if isFirst {
					isFirst = false
				} else {
					io.WriteString(w, ",")
				}
				fmt.Fprintf(w, "%s:", t.Field(i).Name)
				fprint(w, f)
			}
			io.WriteString(w, "}")
		default:
			fmt.Fprintf(w, "%v", v)
		}
	}
}

func fprintIndent(w io.Writer, rv reflect.Value, indent string, depth int) {
	switch v := rv.Interface().(type) {
	case string:
		fmt.Fprintf(w, "\"%s\"", v)
	case []byte:
		fmt.Fprintf(w, "\"%x\"", v)
	case fmt.Stringer:
		fmt.Fprintf(w, "\"%s\"", v.String())
	default:
		switch rv.Kind() {
		case reflect.Ptr:
			if rv.IsNil() {
				io.WriteString(w, "nil")
			} else {
				io.WriteString(w, "&")
				fprintIndent(w, rv.Elem(), indent, depth)
			}
		case reflect.Interface:
			if rv.IsNil() {
				io.WriteString(w, "nil")
			} else {
				fprintIndent(w, rv.Elem(), indent, depth)
			}
		case reflect.Array, reflect.Slice:
			if rv.Len() == 0 {
				io.WriteString(w, "[]")
			} else {
				io.WriteString(w, "[\n")
				for i := 0; i < rv.Len(); i++ {
					writeIndent(w, indent, depth+1)
					fprintIndent(w, rv.Index(i), indent, depth+1)
					io.WriteString(w, ",\n")
				}
				writeIndent(w, indent, depth)
				io.WriteString(w, "]")
			}
		case reflect.Map:
			if rv.Len() == 0 {
				io.WriteString(w, "{}")
			} else {
				iter := rv.MapRange()
				io.WriteString(w, "{\n")
				for iter.Next() {
					writeIndent(w, indent, depth+1)
					fprint(w, iter.Key())
					io.WriteString(w, ": ")
					fprintIndent(w, iter.Value(), indent, depth+1)
					io.WriteString(w, ",\n")
				}
				writeIndent(w, indent, depth)
				io.WriteString(w, "}")
			}
		case reflect.Struct:
			t := rv.Type()
			fmt.Fprintf(w, "%s{\n", t)
			for i := 0; i < rv.NumField(); i++ {
				f := rv.Field(i)
				if k := f.Kind(); (k == reflect.Chan || k == reflect.Func || k == reflect.Map || k == reflect.Pointer || k == reflect.UnsafePointer || k == reflect.Interface || k == reflect.Slice) && f.IsNil() {
					continue
				}
				if !f.CanInterface() {
					if !f.CanAddr() {
						continue
					}
					f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
				}
				writeIndent(w, indent, depth+1)
				fmt.Fprintf(w, "%s: ", t.Field(i).Name)
				fprintIndent(w, f, indent, depth+1)
				io.WriteString(w, ",\n")
			}
			writeIndent(w, indent, depth)
			io.WriteString(w, "}")
		default:
			fmt.Fprintf(w, "%v", v)
		}
	}
}

func writeIndent(w io.Writer, indent string, depth int) {
	for i := 0; i < depth; i++ {
		io.WriteString(w, indent)
	}
}
