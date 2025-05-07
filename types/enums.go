package types

import (
	"encoding/json"
	"fmt"
	"iter"
	"strconv"
	"strings"
)

// OutputMode is a custom type to handle the output mode of the API response. Valid values are ["json","xml"]. Defaults to "json".
type OutputMode = string

const (
	JSON OutputMode = "json"
	XML  OutputMode = "xml" // Deprecated: Do not use it with JSON parsing, it will most likely error.
)

// BooleanYN is a custom type to handle boolean values marshaled as "yes" or "no".
// UnmarshalJSON can handle receiving "t", "f", "yes", "no", true, false (both as strings or booleans).
// Typically, responses return "t" or "f" for true and false, while requests use Yes and No.
type BooleanYN bool

var (
	Yes BooleanYN = true
	No  BooleanYN = false
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
	var d any
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}
	switch d := d.(type) {
	case string:
		switch d {
		case "t", "yes", "true":
			*b = true
		case "f", "no", "false":
			*b = false
		default:
			return fmt.Errorf(`allowed values for Boolean [t, f], [yes, no], [true, false], got %s`, d)
		}
	case bool:
		*b = BooleanYN(d)
	default:
		return fmt.Errorf("invalid type for boolean: %T", d)
	}
	return nil
}

func (b BooleanYN) String() string {
	if b {
		return "yes"
	}
	return "no"
}

func (b BooleanYN) Byte() byte {
	if b {
		return '1'
	}
	return '0'
}

func (b BooleanYN) Int() int {
	if b {
		return 1
	}
	return 0
}

func (b BooleanYN) Bool() bool {
	return bool(b)
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
	if string(data) == "null" {
		return nil
	}
	atoi, err := strconv.Atoi(strings.ReplaceAll(string(data), `"`, ""))
	if err != nil {
		return fmt.Errorf("failed to convert data %s to int: %w", data, err)
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

func (i IntString) Iter() iter.Seq[IntString] {
	return func(yield func(IntString) bool) {
		for i := range i.Int() {
			if !yield(IntString(i)) {
				return
			}
		}
	}
}

// PriceString is a custom type to handle float64 values marshaled as strings ($USD). Typically only returned by responses.
type PriceString float64

func (i PriceString) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.Itoa(int(i)))
}

func (i *PriceString) UnmarshalJSON(data []byte) error {
	_, err := fmt.Sscanf(strings.ReplaceAll(string(data), `"`, ""), `$%f`, i)
	if err != nil {
		return fmt.Errorf("failed to convert data %s to float64: %w", data, err)
	}
	return nil
}

func (i PriceString) String() string {
	return fmt.Sprintf("$%.2f", i)
}

func (i PriceString) Float() float64 {
	return float64(i)
}

// FieldJoinType Defines the union between keywords, description, writing and title search fields. Possible values are "or", "and".
//   - "or" will return submissions found that have the search text in any one of the chosen fields (The default and recommended settings).
//   - "and" will ONLY return submissions that have the search text found in ALL of the chosen fields (unusual and not recommended).
type FieldJoinType = string

const (
	FieldJoinTypeAnd FieldJoinType = "and" // "or" will return submissions found that have the search text in any one of the chosen fields (The default and recommended settings).
	FieldJoinTypeOr  FieldJoinType = "or"  // "and" will ONLY return submissions that have the search text found in ALL of the chosen fields (unusual and not recommended).
)

type JoinType = string

const (
	JoinTypeAnd   JoinType = "and"
	JoinTypeOr    JoinType = "or"
	JoinTypeExact JoinType = "exact"
)

// OrderBy for SubmissionSearchRequest. Default is OrderByCreateDatetime = "create_datetime".
type OrderBy = string

const (
	OrderByDefault        OrderBy = "create_datetime"
	OrderByCreateDatetime OrderBy = "create_datetime"
	OrderByUnreadDatetime OrderBy = "unread_datetime"
	OrderByViews          OrderBy = "views"
	OrderByTotalPrint     OrderBy = "total_print_sales"
	OrderByTotalDigital   OrderBy = "total_digital_sales"
	OrderByTotalSales     OrderBy = "total_sales"
	OrderByUsername       OrderBy = "username"
	OrderByFavDatetime    OrderBy = "fav_datetime"
	OrderByFavStars       OrderBy = "fav_stars"
	OrderByFavs           OrderBy = "favs" // undocumented but exists
	OrderByPoolOrder      OrderBy = "pool_order"
)

// SalesFilter for SubmissionSearchRequest.
//
//	Filter by sales status. Possible options are "forsale" (for sale by any method), "digital" (digital sales), "prints" (print sales).
type SalesFilter = string

const (
	SalesFilterAll     SalesFilter = "forsale"
	SalesFilterDigital SalesFilter = "digital"
	SalesFilterPrints  SalesFilter = "prints"
)
