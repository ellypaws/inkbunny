package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// BooleanYN is a custom type to handle boolean values marshaled as "yes" or "no".
type BooleanYN bool

const (
	Yes BooleanYN = true
	No  BooleanYN = false

	True  BooleanYN = true
	False BooleanYN = false
)

// MarshalJSON converts the BooleanYN boolean into a JSON string of "yes" or "no".
// Typically used for requests as part of url.Values.
func (b BooleanYN) MarshalJSON() ([]byte, error) {
	if b {
		return json.Marshal("yes")
	}
	return json.Marshal("no")
}

// UnmarshalJSON parses string booleans into a BooleanYN type.
// Typically, responses returns "t" or "f" for true and false, while requests use "yes" and "no".
func (b *BooleanYN) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "t", "yes", "true":
		*b = true
	case "f", "no", "false":
		*b = false
	default:
		return fmt.Errorf(`allowed values for Boolean ["t","yes","true","f","no","false"], got "%s"`, s)
	}
	return nil
}

// IntString is a custom type to handle int values marshaled as strings. Typically only returned by responses.
type IntString int

func (i IntString) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.Itoa(int(i)))
}

func (i *IntString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("failed to unmarshal string: %w", err)
	}
	atoi, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf(`failed to convert string "%s" to int: %w`, s, err)
	}
	*i = IntString(atoi)
	return nil
}

func (i IntString) String() string {
	return strconv.Itoa(int(i))
}

func (i IntString) Int() int {
	return int(i)
}

// PriceString is a custom type to handle float64 values marshaled as strings ($USD). Typically only returned by responses.
type PriceString float64

func (i PriceString) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.Itoa(int(i)))
}

func (i *PriceString) UnmarshalJSON(data []byte) error {
	var f string
	if err := json.Unmarshal(data, &f); err != nil {
		return fmt.Errorf("failed to unmarshal float: %w", err)
	}
	// Sscanf is used to parse the string into a float64, as it can handle the $USD format.
	_, err := fmt.Sscanf(f, "$%f", i)
	if err != nil {
		return fmt.Errorf(`failed to convert string "%s" to float64: %w`, f, err)
	}
	return nil
}

func (i PriceString) String() string {
	return fmt.Sprintf("$%.2f", i)
}

func (i PriceString) Float() float64 {
	return float64(i)
}

type Credentials struct {
	Sid      string `json:"sid"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Ratings
}

type Ratings struct {
	// RatingsMask - Binary string representation of the users Allowed Ratings choice. The bits are in this order left-to-right:
	// Eg: A string 11100 means only items rated General, Nudity and Violence are allowed, but Sex and Strong Violence are blocked.
	// A string 11111 means items of any rating would be shown. Only 'left-most significant bits' are returned. So 11010 and 1101 are the same, and 10000 and 1 are the same.
	RatingsMask    string `json:"ratingsmask"`
	General        bool   `json:"1,omitempty"` // Show images with Rating tag: General - Suitable for all ages.
	Nudity         bool   `json:"2,omitempty"` // Show images with Rating tag: Nudity - Nonsexual nudity exposing breasts or genitals (must not show arousal).
	MildViolence   bool   `json:"3,omitempty"` // Show images with Rating tag: MildViolence - Mild violence.
	Sexual         bool   `json:"4,omitempty"` // Show images with Rating tag: Sexual Themes - Erotic imagery, sexual activity or arousal.
	StrongViolence bool   `json:"5,omitempty"` // Show images with Rating tag: StrongViolence - Strong violence, blood, serious injury or death.
}

type LogoutResponse struct {
	Sid    string `json:"sid"`
	Logout string `json:"logout"`
}

type WatchlistResponse struct {
	Watches []UsernameID `json:"watches"`
}

type UsernameID struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

type UsernameAutocomplete struct {
	Results []Autocomplete `json:"results"`
}

type Autocomplete struct {
	ID         string `json:"id"`
	Value      string `json:"value"`
	Icon       string `json:"icon"`
	Info       string `json:"info"`
	SingleWord string `json:"singleword"`
	SearchTerm string `json:"searchterm"`
}

func (u UsernameAutocomplete) String(i int) string {
	return u.Results[i].Value
}

func (u UsernameAutocomplete) Len() int {
	return len(u.Results)
}

type KeywordAutocomplete struct {
	Results []struct {
		Autocomplete
		SubmissionsCount int `json:"submissions_count"`
	} `json:"results"`
}
