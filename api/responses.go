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

type Submission struct {
	SubmissionID     string `json:"submission_id"`
	Hidden           bool   `json:"hidden"`
	Username         string `json:"username"`
	UserID           string `json:"user_id"`
	CreateTimeSystem string `json:"create_datetime"`
	CreateTimeUser   string `json:"create_datetime_usertime"`
	UpdateTimeSystem string `json:"last_file_update_datetime,omitempty"`
	UpdateTimeUser   string `json:"last_file_update_datetime_usertime,omitempty"`
	Title            string `json:"title"`
	Deleted          bool   `json:"deleted"`
	Public           bool   `json:"public"`
	MimeType         string `json:"mimetype"`
	PageCount        int    `json:"pagecount"`
	LatestMimeType   string `json:"latest_mimetype"`
	RatingID         int    `json:"rating_id"`
	RatingName       string `json:"rating_name"`
	ThumbnailDimensions
	SubmissionTypeID int    `json:"submission_type_id"`
	TypeName         string `json:"type_name"`
	Digitalsales     bool   `json:"digitalsales"`
	Printsales       bool   `json:"printsales"`
	FriendsOnly      bool   `json:"friends_only"`
	GuestBlock       bool   `json:"guest_block"`
	Scraps           bool   `json:"scraps"`
}

const (
	SubmissionTypes = iota
	SubmissionTypePicturePinup
	SubmissionTypeSketch
	SubmissionTypePictureSeries
	SubmissionTypeComic
	SubmissionTypePortfolio
	SubmissionTypeShockwaveFlashAnimation
	SubmissionTypeShockwaveFlashInteractive
	SubmissionTypeVideoFeatureLength
	SubmissionTypeVideoAnimation3DCGI
	SubmissionTypeMusicSingleTrack
	SubmissionTypeMusicAlbum
	SubmissionTypeWritingDocument
	SubmissionTypeCharacterSheet
	SubmissionTypePhotography
)

type ThumbnailDimensions struct {
	ThumbMediumX                int `json:"thumb_medium_x,omitempty"`
	ThumbLargeX                 int `json:"thumb_large_x,omitempty"`
	ThumbHugeX                  int `json:"thumb_huge_x,omitempty"`
	ThumbMediumY                int `json:"thumb_medium_y,omitempty"`
	ThumbLargeY                 int `json:"thumb_large_y,omitempty"`
	ThumbHugeY                  int `json:"thumb_huge_y,omitempty"`
	ThumbMediumNonCustomX       int `json:"thumb_medium_noncustom_x,omitempty"`
	ThumbLargeNonCustomX        int `json:"thumb_large_noncustom_x,omitempty"`
	ThumbHugeNonCustomX         int `json:"thumb_huge_noncustom_x,omitempty"`
	ThumbMediumNonCustomY       int `json:"thumb_medium_noncustom_y,omitempty"`
	ThumbLargeNonCustomY        int `json:"thumb_large_noncustom_y,omitempty"`
	ThumbHugeNonCustomY         int `json:"thumb_huge_noncustom_y,omitempty"`
	LatestThumbMediumX          int `json:"latest_thumb_medium_x,omitempty"`
	LatestThumbLargeX           int `json:"latest_thumb_large_x,omitempty"`
	LatestThumbHugeX            int `json:"latest_thumb_huge_x,omitempty"`
	LatestThumbMediumY          int `json:"latest_thumb_medium_y,omitempty"`
	LatestThumbLargeY           int `json:"latest_thumb_large_y,omitempty"`
	LatestThumbHugeY            int `json:"latest_thumb_huge_y,omitempty"`
	LatestThumbMediumNonCustomX int `json:"latest_thumb_medium_noncustom_x,omitempty"`
	LatestThumbLargeNonCustomX  int `json:"latest_thumb_large_noncustom_x,omitempty"`
	LatestThumbHugeNonCustomX   int `json:"latest_thumb_huge_noncustom_x,omitempty"`
	LatestThumbMediumNonCustomY int `json:"latest_thumb_medium_noncustom_y,omitempty"`
	LatestThumbLargeNonCustomY  int `json:"latest_thumb_large_noncustom_y,omitempty"`
	LatestThumbHugeNonCustomY   int `json:"latest_thumb_huge_noncustom_y,omitempty"`
}

type SubmissionDetailsResponse struct {
	Sid          string       `json:"sid"`
	ResultsCount int          `json:"results_count"`
	UserLocation string       `json:"user_location"`
	Submissions  []Submission `json:"submissions"`
}

type SubmissionFavoritesResponse struct {
	Sid   string      `json:"sid"`
	Users []BasicUser `json:"favingusers"`
}
