package utils

import (
	"errors"
	"github.com/charmbracelet/bubbletea"
	"github.com/sahilm/fuzzy"
	"inkbunny/api"
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

// GetSingleUser gets a single user by username, returns an error if no user is found
func GetSingleUser(username string) (api.Autocomplete, error) {
	users, err := api.GetUserID(username)
	if err != nil {
		return api.Autocomplete{}, err
	}
	if len(users.Results) == 0 {
		return api.Autocomplete{}, errors.New("user not found")
	}
	// sort by the closest match using fuzzy
	matches := fuzzy.FindFrom(username, users)
	if len(matches) == 0 {
		return api.Autocomplete{}, errors.New("user not found")
	}
	return users.Results[matches[0].Index], nil
}
