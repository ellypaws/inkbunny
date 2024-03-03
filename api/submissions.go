package api

import (
	"encoding/json"
	"errors"
	"github.com/ellypaws/inkbunny/utils"
	"strings"
)

// BooleanYN is a custom type to handle boolean values marshaled as "yes" or "no".
type BooleanYN bool

const (
	Yes BooleanYN = true
	No  BooleanYN = false
)

// MarshalJSON converts the BooleanYN boolean into a JSON string of "yes" or "no".
func (b BooleanYN) MarshalJSON() ([]byte, error) {
	if b {
		return json.Marshal("yes")
	}
	return json.Marshal("no")
}

// UnmarshalJSON parses a JSON "yes" or "no" into a BooleanYN boolean.
func (b *BooleanYN) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "yes":
		*b = true
	case "no":
		*b = false
	default:
		return errors.New("boolean must be 'yes' or 'no'")
	}
	return nil
}

// SubmissionDetailsRequest is modified to use BooleanYN for fields requiring "yes" or "no" representation.
type SubmissionDetailsRequest struct {
	SID                         string `json:"sid"`
	SubmissionIDs               string `json:"submission_ids"`
	SubmissionIDSlice           []string
	OutputMode                  string    `json:"output_mode"`
	SortKeywordsBy              string    `json:"sort_keywords_by"`
	ShowDescription             BooleanYN `json:"show_description"`
	ShowDescriptionBbcodeParsed BooleanYN `json:"show_description_bbcode_parsed"`
	ShowWriting                 BooleanYN `json:"show_writing"`
	ShowWritingBbcodeParsed     BooleanYN `json:"show_writing_bbcode_parsed"`
	ShowPools                   BooleanYN `json:"show_pools"`
}

type SubmissionBasic struct {
	SubmissionID                string `json:"submission_id"`
	Hidden                      bool   `json:"hidden"`
	Username                    string `json:"username"`
	UserID                      string `json:"user_id"`
	CreateDateSystem            string `json:"create_datetime"`
	CreateDateUser              string `json:"create_datetime_usertime"`
	UpdateDateSystem            string `json:"last_file_update_datetime,omitempty"`
	UpdateDateUser              string `json:"last_file_update_datetime_usertime,omitempty"`
	FileName                    string `json:"file_name"`
	LatestFileName              string `json:"latest_file_name"`
	ThumbnailURL                string `json:"thumbnail_url,omitempty"`
	ThumbnailURLNonCustom       string `json:"thumbnail_url_noncustom,omitempty"`
	LatestThumbnailURL          string `json:"latest_thumbnail_url,omitempty"`
	LatestThumbnailURLNonCustom string `json:"latest_thumbnail_url_noncustom,omitempty"`
	Title                       string `json:"title"`
	Deleted                     bool   `json:"deleted"`
	Public                      bool   `json:"public"`
	MimeType                    string `json:"mimetype"`
	LatestMimeType              string `json:"latest_mimetype"`
	PageCount                   int    `json:"pagecount"`
	RatingID                    int    `json:"rating_id"`
	RatingName                  string `json:"rating_name"`
	ThumbnailDimensions
	SubmissionTypeID int    `json:"submission_type_id"`
	TypeName         string `json:"type_name"`
	Digitalsales     bool   `json:"digitalsales"`
	Printsales       bool   `json:"printsales"`
	FriendsOnly      bool   `json:"friends_only"`
	GuestBlock       bool   `json:"guest_block"`
	Scraps           bool   `json:"scraps"`
}

type SubmissionResponse struct {
	Sid          string `json:"sid"`
	ResultsCount int    `json:"results_count"`
	UserLocation string `json:"user_location"`
	Submissions  []struct {
		SubmissionBasic
		Keywords []struct {
			KeywordID   string `json:"keyword_id"`
			KeywordName string `json:"keyword_name"`
			Suggested   bool   `json:"contributed"`
			Count       int    `json:"submissions_count"`
		} `json:"keywords"`
		Favorite         bool   `json:"favorite"`
		FavoritesCount   int    `json:"favorites_count"`
		UserIconFileName string `json:"user_icon_file_name"`
		UserIconURL      struct {
			Large  string `json:"user_icon_url_large,omitempty"`
			Medium string `json:"user_icon_url_medium,omitempty"`
			Small  string `json:"user_icon_url_small,omitempty"`
		}
		Files []struct {
			FileID                string `json:"file_id"`
			FileName              string `json:"file_name"`
			ThumbnailURL          string `json:"thumbnail_url,omitempty"`
			ThumbnailURLNonCustom string `json:"thumbnail_url_noncustom,omitempty"`
			FileURL               string `json:"file_url,omitempty"`
			MimeType              string `json:"mimetype"`
			SubmissionID          string `json:"submission_id"`
			UserID                string `json:"user_id"`
			SubmissionFileOrder   int    `json:"submission_file_order"`
			SizeX                 int    `json:"full_size_x"`
			SizeY                 int    `json:"full_size_y"`
			ThumbX                int    `json:"thumb_huge_x,omitempty"`
			ThumbY                int    `json:"thumb_huge_y,omitempty"`
			ThumbNonCustomX       int    `json:"thumb_huge_noncustom_x,omitempty"`
			ThumbNonCustomY       int    `json:"thumb_huge_noncustom_y,omitempty"`
			InitialFileMD5        string `json:"initial_file_md5"`
			FullFileMD5           string `json:"full_file_md5"`
			LargeFileMD5          string `json:"large_file_md5"`
			SmallFileMD5          string `json:"small_file_md5"`
			ThumbnailMD5          string `json:"thumbnail_md5"`
			Deleted               bool   `json:"deleted"`
			CreateDateTime        string `json:"create_datetime"`
			CreateDateTimeUser    string `json:"create_datetime_usertime"`
		} `json:"files"`
		Pools []struct {
			PoolID                     string `json:"pool_id"`
			Name                       string `json:"name"`
			Description                string `json:"description"`
			Count                      int    `json:"count"`
			LeftSubmissionID           string `json:"submission_left_submission_id"`
			RightSubmissionID          string `json:"submission_right_submission_id"`
			LeftSubmissionFileName     string `json:"submission_left_file_name"`
			RightSubmissionFileName    string `json:"submission_right_file_name"`
			LeftThumbnailURL           string `json:"submission_left_thumbnail_url,omitempty"`
			RightThumbnailURL          string `json:"submission_right_thumbnail_url,omitempty"`
			LeftThumbnailURLNonCustom  string `json:"submission_left_thumbnail_url_noncustom,omitempty"`
			RightThumbnailURLNonCustom string `json:"submission_right_thumbnail_url_noncustom,omitempty"`
			LeftThumbX                 int    `json:"submission_left_thumb_huge_x,omitempty"`
			LeftThumbY                 int    `json:"submission_left_thumb_huge_y,omitempty"`
			RightThumbX                int    `json:"submission_right_thumb_huge_x,omitempty"`
			RightThumbY                int    `json:"submission_right_thumb_huge_y,omitempty"`
			LeftThumbNonCustomX        int    `json:"submission_left_thumb_huge_noncustom_x,omitempty"`
			LeftThumbNonCustomY        int    `json:"submission_left_thumb_huge_noncustom_y,omitempty"`
			RightThumbNonCustomX       int    `json:"submission_right_thumb_huge_noncustom_x,omitempty"`
			RightThumbNonCustomY       int    `json:"submission_right_thumb_huge_noncustom_y,omitempty"`
		} `json:"pools"`
		Description             string `json:"description"`
		DescriptionBBCodeParsed string `json:"description_bbcode_parsed"`
		Writing                 string `json:"writing"`
		WritingBBCodeParsed     string `json:"writing_bbcode_parsed"`
		PoolsCount              int    `json:"pools_count"`
		Ratings                 []struct {
			ContentTagID int    `json:"content_tag_id"`
			Name         string `json:"name"`
			Description  string `json:"description"`
			RatingID     int    `json:"rating_id"`
		} `json:"ratings"`
		CommentsCount    int    `json:"comments_count"`
		Views            int    `json:"views"`
		SalesDescription string `json:"sales_description"`
		ForSale          bool   `json:"forsale"`
		DigitalPrice     int    `json:"digital_price"`
		Prints           []struct {
			PrintSizeID        int    `json:"print_size_id"`
			Name               string `json:"name"`
			Price              int    `json:"price"`
			PriceOwnerDiscount int    `json:"price_owner_discount,omitempty"`
		} `json:"prints"`
	} `json:"submissions"`
}

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
	Sid          string               `json:"sid"`
	ResultsCount int                  `json:"results_count"`
	UserLocation string               `json:"user_location"`
	Submissions  []SubmissionResponse `json:"submissions"`
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

func (user Credentials) SubmissionDetails(req SubmissionDetailsRequest) (SubmissionDetailsResponse, error) {
	if !user.LoggedIn() {
		return SubmissionDetailsResponse{}, ErrNotLoggedIn
	}
	if req.SID == "" {
		req.SID = user.Sid
	}

	urlValues := utils.StructToUrlValues(req)

	if !urlValues.Has("submission_ids") && len(req.SubmissionIDSlice) > 0 {
		urlValues.Set("submission_ids", strings.Join(req.SubmissionIDSlice, ","))
	}

	resp, err := user.Get(apiURL("submission", urlValues))
	if err != nil {
		return SubmissionDetailsResponse{}, err
	}
	defer resp.Body.Close()

	var submission SubmissionDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&submission); err != nil {
		return SubmissionDetailsResponse{}, err
	}
	return submission, nil

}

type SubmissionRequest struct {
	SID          string `json:"sid"`
	SubmissionID string `json:"submission_id"`
	OutputMode   string `json:"output_mode,omitempty"`
}

type SubmissionFavoritesResponse struct {
	Sid   string       `json:"sid"`
	Users []UsernameID `json:"favingusers"`
}

func (user Credentials) SubmissionFavorites(req SubmissionRequest) (SubmissionFavoritesResponse, error) {
	if !user.LoggedIn() {
		return SubmissionFavoritesResponse{}, ErrNotLoggedIn
	}
	if req.SID == "" {
		req.SID = user.Sid
	}

	resp, err := user.Get(apiURL("submissionfavingusers", utils.StructToUrlValues(req)))
	if err != nil {
		return SubmissionFavoritesResponse{}, err
	}
	defer resp.Body.Close()

	var favorites SubmissionFavoritesResponse
	if err := json.NewDecoder(resp.Body).Decode(&favorites); err != nil {
		return SubmissionFavoritesResponse{}, err
	}
	return favorites, nil
}
