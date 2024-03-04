package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

func (b BooleanYN) String() string {
	if b {
		return "yes"
	}
	return "no"
}

func (b BooleanYN) Int() int {
	if b {
		return 1
	}
	return 0
}

// IntString is a custom type to handle int values marshaled as strings. Typically only returned by responses.
type IntString int

func (i IntString) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func (i *IntString) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	atoi, err := strconv.Atoi(strings.ReplaceAll(string(data), `"`, ""))
	if err != nil {
		return fmt.Errorf(`failed to convert data %s to int: %w`, data, err)
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
	_, err := fmt.Sscanf(strings.ReplaceAll(string(data), `"`, ""), `$%f`, i)
	if err != nil {
		return fmt.Errorf(`failed to convert data %s to float64: %w`, data, err)
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
	Ratings  `json:"ratings,omitempty"`
}

type Ratings struct {
	// RatingsMask - Binary string representation of the users Allowed Ratings choice. The bits are in this order left-to-right:
	// Eg: A string 11100 means only items rated General, Nudity and Violence are allowed, but Sex and Strong Violence are blocked.
	// A string 11111 means items of any rating would be shown. Only 'left-most significant bits' are returned. So 11010 and 1101 are the same, and 10000 and 1 are the same.
	RatingsMask    string    `json:"-"`
	General        BooleanYN `json:"tag[1],omitempty"` // Show images with Rating tag: General - Suitable for all ages.
	Nudity         BooleanYN `json:"tag[2],omitempty"` // Show images with Rating tag: Nudity - Nonsexual nudity exposing breasts or genitals (must not show arousal).
	MildViolence   BooleanYN `json:"tag[3],omitempty"` // Show images with Rating tag: MildViolence - Mild violence.
	Sexual         BooleanYN `json:"tag[4],omitempty"` // Show images with Rating tag: Sexual Themes - Erotic imagery, sexual activity or arousal.
	StrongViolence BooleanYN `json:"tag[5],omitempty"` // Show images with Rating tag: StrongViolence - Strong violence, blood, serious injury or death.
}

// parseMask sets the Ratings boolean fields based on the RatingsMask field. True is 1, false is 0
//
//	"11010" would set Ratings{General: true, Nudity: true, MildViolence: false, Sexual: true, StrongViolence: false}
func (ratings *Ratings) parseMask() {
	// RatingsMask - Binary string representation of the users Allowed Ratings choice. The bits are in this order left-to-right:
	// Eg: A string 11100 means only items rated General, Nudity and Violence are allowed, but Sex and Strong Violence are blocked.
	// A string 11111 means items of any rating would be shown. Only 'left-most significant bits' are returned. So 11010 and 1101 are the same, and 10000 and 1 are the same.
	set := func(r int32) BooleanYN {
		return r == '1'
	}
	for i, rating := range ratings.RatingsMask {
		switch i {
		case 0:
			ratings.General = set(rating)
		case 1:
			ratings.Nudity = set(rating)
		case 2:
			ratings.MildViolence = set(rating)
		case 3:
			ratings.Sexual = set(rating)
		case 4:
			ratings.StrongViolence = set(rating)
		}
	}
}

// parseBooleans sets the RatingsMask field based on the boolean values of the Ratings struct.
//
//	Ratings{General: true, Nudity: true, MildViolence: false, Sexual: true, StrongViolence: false}
//	would return "1101"
//
// RatingsMask is a binary string representation of the users Allowed Ratings choice.
// A string 11100 means only keywords rated General,
// Nudity and Violence are allowed, but Sex and Strong Violence are blocked.
// String 11111 means keywords of any rating would be shown.
// Only 'left-most significant bits' need to be sent.
// So 11010 and 1101 are the same, and 10000 and 1 are the same.
func (ratings *Ratings) parseBooleans() {
	ratings.RatingsMask = fmt.Sprintf("%d%d%d%d%d",
		ratings.General.Int(),
		ratings.Nudity.Int(),
		ratings.MildViolence.Int(),
		ratings.Sexual.Int(),
		ratings.StrongViolence.Int(),
	)

	ratings.RatingsMask = strings.TrimRight(ratings.RatingsMask, "0")
}

func (r Ratings) String() string {
	r.parseBooleans()
	return r.RatingsMask
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
