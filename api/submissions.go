package api

import (
	"encoding/json"
	"fmt"
	"github.com/ellypaws/inkbunny/utils"
	"strings"
)

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
	SubmissionID                string    `json:"submission_id"`
	Hidden                      BooleanYN `json:"hidden"`
	Username                    string    `json:"username"`
	UserID                      string    `json:"user_id"`
	CreateDateSystem            string    `json:"create_datetime"`
	CreateDateUser              string    `json:"create_datetime_usertime"`
	UpdateDateSystem            string    `json:"last_file_update_datetime,omitempty"`
	UpdateDateUser              string    `json:"last_file_update_datetime_usertime,omitempty"`
	FileName                    string    `json:"file_name"`
	LatestFileName              string    `json:"latest_file_name"`
	ThumbnailURL                string    `json:"thumbnail_url,omitempty"`
	ThumbnailURLNonCustom       string    `json:"thumbnail_url_noncustom,omitempty"`
	LatestThumbnailURL          string    `json:"latest_thumbnail_url,omitempty"`
	LatestThumbnailURLNonCustom string    `json:"latest_thumbnail_url_noncustom,omitempty"`
	Title                       string    `json:"title"`
	Deleted                     BooleanYN `json:"deleted"`
	Public                      BooleanYN `json:"public"`
	MimeType                    string    `json:"mimetype"`
	LatestMimeType              string    `json:"latest_mimetype"`
	PageCount                   IntString `json:"pagecount"`
	RatingID                    IntString `json:"rating_id"`
	RatingName                  string    `json:"rating_name"`
	ThumbnailDimensions
	SubmissionTypeID IntString `json:"submission_type_id"`
	TypeName         string    `json:"type_name"`
	Digitalsales     BooleanYN `json:"digitalsales"`
	Printsales       BooleanYN `json:"printsales"`
	FriendsOnly      BooleanYN `json:"friends_only"`
	GuestBlock       BooleanYN `json:"guest_block"`
	Scraps           BooleanYN `json:"scraps"`
}

type Submission struct {
	SubmissionBasic
	Keywords []struct {
		KeywordID   string    `json:"keyword_id"`
		KeywordName string    `json:"keyword_name"`
		Suggested   BooleanYN `json:"contributed"`
		Count       IntString `json:"submissions_count"`
	} `json:"keywords"`
	Favorite         BooleanYN `json:"favorite"`
	FavoritesCount   IntString `json:"favorites_count"`
	UserIconFileName string    `json:"user_icon_file_name"`
	UserIconURL      struct {
		Large  string `json:"user_icon_url_large,omitempty"`
		Medium string `json:"user_icon_url_medium,omitempty"`
		Small  string `json:"user_icon_url_small,omitempty"`
	}
	Files []struct {
		FileID                string    `json:"file_id"`
		FileName              string    `json:"file_name"`
		ThumbnailURL          string    `json:"thumbnail_url,omitempty"`
		ThumbnailURLNonCustom string    `json:"thumbnail_url_noncustom,omitempty"`
		FileURL               string    `json:"file_url,omitempty"`
		MimeType              string    `json:"mimetype"`
		SubmissionID          string    `json:"submission_id"`
		UserID                string    `json:"user_id"`
		SubmissionFileOrder   IntString `json:"submission_file_order"`
		FullSizeX             IntString `json:"full_size_x"`
		FullSizeY             IntString `json:"full_size_y"`
		ThumbHugeX            IntString `json:"thumb_huge_x,omitempty"`
		ThumbHugeY            IntString `json:"thumb_huge_y,omitempty"`
		ThumbNonCustomX       IntString `json:"thumb_huge_noncustom_x,omitempty"`
		ThumbNonCustomY       IntString `json:"thumb_huge_noncustom_y,omitempty"`
		InitialFileMD5        string    `json:"initial_file_md5"`
		FullFileMD5           string    `json:"full_file_md5"`
		LargeFileMD5          string    `json:"large_file_md5"`
		SmallFileMD5          string    `json:"small_file_md5"`
		ThumbnailMD5          string    `json:"thumbnail_md5"`
		Deleted               BooleanYN `json:"deleted"`
		CreateDateTime        string    `json:"create_datetime"`
		CreateDateTimeUser    string    `json:"create_datetime_usertime"`
	} `json:"files"`
	Pools []struct {
		PoolID                     string    `json:"pool_id"`
		Name                       string    `json:"name"`
		Description                string    `json:"description"`
		Count                      IntString `json:"count"`
		LeftSubmissionID           string    `json:"submission_left_submission_id"`
		RightSubmissionID          string    `json:"submission_right_submission_id"`
		LeftSubmissionFileName     string    `json:"submission_left_file_name"`
		RightSubmissionFileName    string    `json:"submission_right_file_name"`
		LeftThumbnailURL           string    `json:"submission_left_thumbnail_url,omitempty"`
		RightThumbnailURL          string    `json:"submission_right_thumbnail_url,omitempty"`
		LeftThumbnailURLNonCustom  string    `json:"submission_left_thumbnail_url_noncustom,omitempty"`
		RightThumbnailURLNonCustom string    `json:"submission_right_thumbnail_url_noncustom,omitempty"`
		LeftThumbX                 IntString `json:"submission_left_thumb_huge_x,omitempty"`
		LeftThumbY                 IntString `json:"submission_left_thumb_huge_y,omitempty"`
		RightThumbX                IntString `json:"submission_right_thumb_huge_x,omitempty"`
		RightThumbY                IntString `json:"submission_right_thumb_huge_y,omitempty"`
		LeftThumbNonCustomX        IntString `json:"submission_left_thumb_huge_noncustom_x,omitempty"`
		LeftThumbNonCustomY        IntString `json:"submission_left_thumb_huge_noncustom_y,omitempty"`
		RightThumbNonCustomX       IntString `json:"submission_right_thumb_huge_noncustom_x,omitempty"`
		RightThumbNonCustomY       IntString `json:"submission_right_thumb_huge_noncustom_y,omitempty"`
	} `json:"pools"`
	Description             string    `json:"description"`
	DescriptionBBCodeParsed string    `json:"description_bbcode_parsed"`
	Writing                 string    `json:"writing"`
	WritingBBCodeParsed     string    `json:"writing_bbcode_parsed"`
	PoolsCount              IntString `json:"pools_count"`
	Ratings                 []struct {
		ContentTagID IntString `json:"content_tag_id"`
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		RatingID     IntString `json:"rating_id"`
	} `json:"ratings"`
	CommentsCount    IntString `json:"comments_count"`
	Views            IntString `json:"views"`
	SalesDescription string    `json:"sales_description"`
	ForSale          BooleanYN `json:"forsale"`
	DigitalPrice     IntString `json:"digital_price"`
	Prints           []struct {
		PrintSizeID        IntString   `json:"print_size_id"`
		Name               string      `json:"name"`
		Price              PriceString `json:"price"`
		PriceOwnerDiscount PriceString `json:"price_owner_discount,omitempty"`
	} `json:"prints"`
}

type ThumbnailDimensions struct {
	ThumbMediumX                IntString `json:"thumb_medium_x,omitempty"`
	ThumbLargeX                 IntString `json:"thumb_large_x,omitempty"`
	ThumbHugeX                  IntString `json:"thumb_huge_x,omitempty"`
	ThumbMediumY                IntString `json:"thumb_medium_y,omitempty"`
	ThumbLargeY                 IntString `json:"thumb_large_y,omitempty"`
	ThumbHugeY                  IntString `json:"thumb_huge_y,omitempty"`
	ThumbMediumNonCustomX       IntString `json:"thumb_medium_noncustom_x,omitempty"`
	ThumbLargeNonCustomX        IntString `json:"thumb_large_noncustom_x,omitempty"`
	ThumbHugeNonCustomX         IntString `json:"thumb_huge_noncustom_x,omitempty"`
	ThumbMediumNonCustomY       IntString `json:"thumb_medium_noncustom_y,omitempty"`
	ThumbLargeNonCustomY        IntString `json:"thumb_large_noncustom_y,omitempty"`
	ThumbHugeNonCustomY         IntString `json:"thumb_huge_noncustom_y,omitempty"`
	LatestThumbMediumX          IntString `json:"latest_thumb_medium_x,omitempty"`
	LatestThumbLargeX           IntString `json:"latest_thumb_large_x,omitempty"`
	LatestThumbHugeX            IntString `json:"latest_thumb_huge_x,omitempty"`
	LatestThumbMediumY          IntString `json:"latest_thumb_medium_y,omitempty"`
	LatestThumbLargeY           IntString `json:"latest_thumb_large_y,omitempty"`
	LatestThumbHugeY            IntString `json:"latest_thumb_huge_y,omitempty"`
	LatestThumbMediumNonCustomX IntString `json:"latest_thumb_medium_noncustom_x,omitempty"`
	LatestThumbLargeNonCustomX  IntString `json:"latest_thumb_large_noncustom_x,omitempty"`
	LatestThumbHugeNonCustomX   IntString `json:"latest_thumb_huge_noncustom_x,omitempty"`
	LatestThumbMediumNonCustomY IntString `json:"latest_thumb_medium_noncustom_y,omitempty"`
	LatestThumbLargeNonCustomY  IntString `json:"latest_thumb_large_noncustom_y,omitempty"`
	LatestThumbHugeNonCustomY   IntString `json:"latest_thumb_huge_noncustom_y,omitempty"`
}

type SubmissionDetailsResponse struct {
	Sid          string       `json:"sid"`
	ResultsCount IntString    `json:"results_count"`
	UserLocation string       `json:"user_location"`
	Submissions  []Submission `json:"submissions"`
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

	resp, err := user.Get(apiURL("submissions", urlValues))
	if err != nil {
		return SubmissionDetailsResponse{}, fmt.Errorf("failed to get submission details: %w", err)
	}
	defer resp.Body.Close()

	var submission SubmissionDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&submission); err != nil {
		return SubmissionDetailsResponse{}, fmt.Errorf("failed to unmarshal response body: %w", err)
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
