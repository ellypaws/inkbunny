package types

import (
	"strings"
)

// Ratings to use when calling Client.ChangeRatings
//
//	err := client.ChangeRatings(user, General|Nudity)
const (
	General uint8 = 0x10 >> iota
	Nudity
	MildViolence
	Sexual
	StrongViolence
)

// Ratings - Binary string representation of the users Allowed Ratings choice. The bits are in this order left-to-right:
// Eg: A string 11100 means only items rated General, Nudity and Violence are allowed, but Sex and Strong Violence are blocked.
// A string 11111 means items of any rating would be shown. Only 'left-most significant bits' are returned. So 11010 and 1101 are the same, and 10000 and 1 are the same.
// Ratings implements json.Unmarshaler which converts string or uint8 bitmask into Ratings.
// It also implements fmt.Stringer which calls its String method, converting the values into an LSB bitmask.
type Ratings struct {
	General        *BooleanYN `json:"tag[1],omitempty" query:"tag[1]"` // Show images with Rating tag: General - Suitable for all ages.
	Nudity         *BooleanYN `json:"tag[2],omitempty" query:"tag[2]"` // Show images with Rating tag: Nudity - Nonsexual nudity exposing breasts or genitals (must not show arousal).
	MildViolence   *BooleanYN `json:"tag[3],omitempty" query:"tag[3]"` // Show images with Rating tag: MildViolence - Mild violence.
	Sexual         *BooleanYN `json:"tag[4],omitempty" query:"tag[4]"` // Show images with Rating tag: Sexual Themes - Erotic imagery, sexual activity or arousal.
	StrongViolence *BooleanYN `json:"tag[5],omitempty" query:"tag[5]"` // Show images with Rating tag: StrongViolence - Strong violence, blood, serious injury or death.
}

func (r *Ratings) UnmarshalJSON(b []byte) error {
	*r = ParseMask(strings.Trim(string(b), `"`))
	return nil
}

func (r Ratings) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

// ParseMaskU returns a Ratings based on a ratings bitmask. True is 1, false is 0
//
//	ratings := ParseMaskU(General|Nudity)
func ParseMaskU(mask uint8) Ratings {
	return Ratings{
		General:        (*BooleanYN)(Address(mask&General != 0)),
		Nudity:         (*BooleanYN)(Address(mask&Nudity != 0)),
		MildViolence:   (*BooleanYN)(Address(mask&MildViolence != 0)),
		Sexual:         (*BooleanYN)(Address(mask&Sexual != 0)),
		StrongViolence: (*BooleanYN)(Address(mask&StrongViolence != 0)),
	}
}

func Address[T any](v T) *T {
	return &v
}

// ParseMask returns a Ratings based on a ratings bitmask. True is 1, false is 0
//
//	"11010" would set Ratings{General: true, Nudity: true, MildViolence: false, Sexual: true, StrongViolence: false}
func ParseMask(s string) Ratings {
	// RatingsMask - Binary string representation of the users Allowed Ratings choice. The bits are in this order left-to-right:
	// Eg: A string 11100 means only items rated General, Nudity and Violence are allowed, but Sex and Strong Violence are blocked.
	// A string 11111 means items of any rating would be shown. Only 'left-most significant bits' are returned. So 11010 and 1101 are the same, and 10000 and 1 are the same.
	var m uint8
	for i := 0; i < len(s) && i < 5; i++ {
		if s[i] == '1' {
			// bit positions: General → the highest bit (0x10), down to 0x01
			m |= 1 << (4 - i)
		}
	}
	return ParseMaskU(m)
}

// String returns a ratings bitmask based on the boolean values of the Ratings struct.
//
//	Ratings{General: true, Nudity: true, MildViolence: false, Sexual: true, StrongViolence: false}
//	would return "1101"
//
// A ratings bitmask is a binary string representation of the users Allowed Ratings choice.
// A string 11100 means only keywords rated General,
// Nudity and Violence are allowed, but Sex and Strong Violence are blocked.
// String 11111 means keywords of any rating would be shown.
// Only 'left-most significant bits' need to be sent.
// So 11010 and 1101 are the same, and 10000 and 1 are the same.
func (r Ratings) String() string {
	// build a 5-byte buffer of '0' or '1'
	var b [5]byte
	b[0] = r.General.Byte()
	b[1] = r.Nudity.Byte()
	b[2] = r.MildViolence.Byte()
	b[3] = r.Sexual.Byte()
	b[4] = r.StrongViolence.Byte()

	// trim trailing '0's by finding the last '1'
	last := len(b) - 1
	for last >= 0 && b[last] == '0' {
		last--
	}
	if last < 0 {
		return ""
	}
	return string(b[:last+1])
}

// Byte returns a 5-bit mask with General in the MSB and StrongViolence in the LSB.
//
//	bit 4 ── General
//	bit 3 ── Nudity
//	bit 2 ── MildViolence
//	bit 1 ── Sexual
//	bit 0 ── StrongViolence
func (r Ratings) Byte() byte {
	var b byte
	if r.General.Bool() {
		b |= 1 << 4
	}
	if r.Nudity.Bool() {
		b |= 1 << 3
	}
	if r.MildViolence.Bool() {
		b |= 1 << 2
	}
	if r.Sexual.Bool() {
		b |= 1 << 1
	}
	if r.StrongViolence.Bool() {
		b |= 1 << 0
	}
	return b
}
