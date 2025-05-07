package utils

import (
	"encoding"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type fieldInfo struct {
	name      string
	noTag     bool
	omitEmpty bool
	index     []int
	kind      reflect.Kind
}

var fieldCache sync.Map // map[reflect.Type][]fieldInfo

func getFieldInfos(t reflect.Type) []fieldInfo {
	if cached, ok := fieldCache.Load(t); ok {
		return cached.([]fieldInfo)
	}

	var infos []fieldInfo
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		// skip unexported
		if f.PkgPath != "" {
			continue
		}

		tag := f.Tag.Get("json")
		name, omit, skip, noTag := "", false, false, false

		if tag == "" {
			// no explicit tag: use field name lowercased
			name = strings.ToLower(f.Name)
			noTag = true
		} else {
			parts := strings.Split(tag, ",")
			// tag "-" means skip
			if parts[0] == "-" {
				skip = true
			} else {
				name = parts[0]
				if len(parts) > 1 && parts[1] == "omitempty" {
					omit = true
				}
			}
		}
		if skip {
			continue
		}

		// anonymous struct → recurse into it instead of treating as a field
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			// grab its fields
			for _, sub := range getFieldInfos(f.Type) {
				// but adjust their index path to go through this embedding
				sub.index = append(f.Index, sub.index...)
				infos = append(infos, sub)
			}
			continue
		}

		infos = append(infos, fieldInfo{
			name:      name,
			noTag:     noTag,
			omitEmpty: omit,
			index:     f.Index,
			kind:      f.Type.Kind(),
		})
	}

	fieldCache.Store(t, infos)
	return infos
}

// StructToUrlValues converts a struct to url.Values.
func StructToUrlValues(s any) url.Values {
	vals := make(url.Values)
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return vals
	}
	t := v.Type()
	infos := getFieldInfos(t)

	for _, fi := range infos {
		fv := v.FieldByIndex(fi.index)

		// handle omitempty
		if fi.omitEmpty && isEmptyValue(fv) {
			continue
		}

		// unwrap pointers
		if fv.Kind() == reflect.Pointer {
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
		}

		// ——— recursion for embedded structs ———
		if fi.noTag && fv.Kind() == reflect.Struct {
			sub := StructToUrlValues(fv.Interface())
			for k, vs := range sub {
				for _, v := range vs {
					vals.Add(k, v)
				}
			}
			continue
		}

		// ——— scalar / slice / interface marshalling ———
		var str string
		iface := fv.Interface()

		switch i := iface.(type) {
		case fmt.Stringer:
			str = i.String()
		case encoding.BinaryMarshaler:
			b, err := i.MarshalBinary()
			if err != nil {
				continue
			}
			str = string(b)
		case encoding.TextMarshaler:
			b, err := i.MarshalText()
			if err != nil {
				continue
			}
			str = string(b)
		case json.Marshaler:
			b, err := i.MarshalJSON()
			if err != nil {
				continue
			}
			str = strings.Trim(string(b), `"`)
		default:
			switch fv.Kind() {
			case reflect.Bool:
				str = strconv.FormatBool(fv.Bool())
			case reflect.Slice:
				parts := make([]string, fv.Len())
				for i := 0; i < fv.Len(); i++ {
					parts[i] = fmt.Sprint(fv.Index(i).Interface())
				}
				str = strings.Join(parts, ",")
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				str = strconv.FormatInt(fv.Int(), 10)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				str = strconv.FormatUint(fv.Uint(), 10)
			case reflect.Float32, reflect.Float64:
				str = strconv.FormatFloat(fv.Float(), 'g', -1, 64)
			case reflect.String:
				str = fv.String()
			default:
				str = fmt.Sprint(iface)
			}
		}

		if str != "" {
			vals.Add(fi.name, str)
		}
	}
	return vals
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		return false
	}
}
