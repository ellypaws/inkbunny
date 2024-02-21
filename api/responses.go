package api

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
	Watches []BasicUser `json:"watches"`
}

type BasicUser struct {
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

type KeywordAutocomplete struct {
	Results []struct {
		Autocomplete
		SubmissionsCount int `json:"submissions_count"`
	} `json:"results"`
}

type SubmissionFavoritesResponse struct {
	Sid   string      `json:"sid"`
	Users []BasicUser `json:"favingusers"`
}
