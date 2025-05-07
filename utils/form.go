package utils

import (
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
)

// StructToMultipartForm builds a multipart/form-data body from a struct.
// It returns the body buffer and the content type (with boundary).
func StructToMultipartForm(s any) (*bytes.Buffer, string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return body, writer.FormDataContentType(), nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, "", fmt.Errorf("StructToMultipartForm: expected struct, got %s", v.Kind())
	}
	// write all fields
	if err := writeFields(writer, v); err != nil {
		return nil, "", err
	}
	contentType := writer.FormDataContentType()
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return body, contentType, nil
}

// StructToMultipartPipe builds a multipart/form-data body from a struct.
// It returns an io.ReadCloser and the content type (with boundary).
func StructToMultipartPipe(s any) (io.ReadCloser, string, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, "", nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, "", fmt.Errorf("StructToMultipartForm: expected struct, got %s", v.Kind())
	}
	// write all fields
	go func() {
		if err := writeFields(writer, v); err != nil {
			pw.CloseWithError(err)
		}
	}()
	return pr, writer.FormDataContentType(), nil
}

// StructToMultipartWriter builds a multipart/form-data body from a struct.
// It writes to a multipart.Writer directly
func StructToMultipartWriter(writer *multipart.Writer, s any) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("StructToMultipartForm: expected struct, got %s", v.Kind())
	}
	if err := writeFields(writer, v); err != nil {
		return err
	}
	return nil
}

// writeFields iterates over struct fields and writes each as a form field.
func writeFields(writer *multipart.Writer, v reflect.Value) error {
	t := v.Type()
	infos := getFieldInfos(t)
	for _, fi := range infos {
		fv := v.FieldByIndex(fi.index)
		// omit empty
		if fi.omitEmpty && isEmptyValue(fv) {
			continue
		}
		// unwrap pointer
		if fv.Kind() == reflect.Pointer {
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
		}
		// embedded struct → recurse
		if fi.noTag && fv.Kind() == reflect.Struct {
			if err := writeFields(writer, fv); err != nil {
				return err
			}
			continue
		}

		var str string
		iface := fv.Interface()
		switch i := iface.(type) {
		case io.Reader:
			var buf strings.Builder
			if _, err := io.Copy(&buf, i); err != nil {
				return err
			}
			str = buf.String()
		case fmt.Stringer:
			str = i.String()
		case encoding.BinaryMarshaler:
			b, err := i.MarshalBinary()
			if err != nil {
				return err
			}
			str = string(b)
		case json.Marshaler:
			b, err := i.MarshalJSON()
			if err != nil {
				return err
			}
			str = strings.Trim(string(b), `"`)
		default:
			switch fv.Kind() {
			case reflect.Bool:
				str = strconv.FormatBool(fv.Bool())
			case reflect.Slice:
				parts := make([]string, fv.Len())
				for i := range fv.Len() {
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
		// write form field
		if str != "" {
			if err := writer.WriteField(fi.name, str); err != nil {
				return err
			}
		}
	}
	return nil
}
