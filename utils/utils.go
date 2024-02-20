package utils

import (
	"github.com/charmbracelet/bubbletea"
	"net/url"
	"reflect"
)

// Wrap casts a message into a tea.Cmd
func Wrap(msg any) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

// StructToUrlValues uses reflect to read json struct fields and set them as url.Values
// It also checks if omitempty is set and ignores empty fields
// Example:
//
//	type Example struct {
//		Field1 string `json:"field1,omitempty"`
//		Field2 string `json:"field2"`
//	}
func StructToUrlValues(s any) url.Values {
	var urlValues url.Values
	v := reflect.ValueOf(s)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if field.Tag.Get("json") == "" {
			continue
		}
		if field.Tag.Get("json") == "omitempty" && value.String() == "" {
			continue
		}
		urlValues.Add(field.Tag.Get("json"), value.String())
	}
	return urlValues
}
