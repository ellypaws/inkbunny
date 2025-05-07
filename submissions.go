package inkbunny

import (
	"net/url"
	"strings"

	"github.com/ellypaws/inkbunny/types"
)

// SubmissionDetailsRequest is modified to use BooleanYN for fields requiring "yes" or "no" representation.
type SubmissionDetailsRequest struct {
	SID                         string           `json:"sid" query:"sid"`
	SubmissionIDs               string           `json:"submission_ids" query:"submission_ids"` // SubmissionIDs is a comma-separated list of submission IDs
	SubmissionIDSlice           []string         `json:"-"`                                     // SubmissionIDSlice will be joined as a comma-separated into SubmissionIDs
	OutputMode                  types.OutputMode `json:"output_mode,omitempty" query:"output_mode"`
	SortKeywordsBy              string           `json:"sort_keywords_by,omitempty" query:"sort_keywords_by"`
	ShowDescription             types.BooleanYN  `json:"show_description,omitempty" query:"show_description"`
	ShowDescriptionBbcodeParsed types.BooleanYN  `json:"show_description_bbcode_parsed,omitempty" query:"show_description_bbcode_parsed"`
	ShowWriting                 types.BooleanYN  `json:"show_writing,omitempty" query:"show_writing"`
	ShowWritingBbcodeParsed     types.BooleanYN  `json:"show_writing_bbcode_parsed,omitempty" query:"show_writing_bbcode_parsed"`
	ShowPools                   types.BooleanYN  `json:"show_pools,omitempty" query:"show_pools"`
}

// SubmissionBasic combines elements common in SubmissionSearch and SubmissionDetails
type SubmissionBasic struct {
	SubmissionID     types.IntString `json:"submission_id"`
	Hidden           types.BooleanYN `json:"hidden,omitempty"`
	Username         string          `json:"username,omitempty"`
	UserID           types.IntString `json:"user_id,omitempty"`
	CreateDateSystem string          `json:"create_datetime,omitempty"`
	CreateDateUser   string          `json:"create_datetime_usertime,omitempty"`
	UpdateDateSystem string          `json:"last_file_update_datetime,omitempty"`
	UpdateDateUser   string          `json:"last_file_update_datetime_usertime,omitempty"`
	FileName         string          `json:"file_name,omitempty"`
	LatestFileName   string          `json:"latest_file_name,omitempty"`
	Title            string          `json:"title,omitempty"`
	Deleted          types.BooleanYN `json:"deleted,omitempty"`
	Public           types.BooleanYN `json:"public,omitempty"`
	MimeType         string          `json:"mimetype,omitempty"`
	LatestMimeType   string          `json:"latest_mimetype,omitempty"`
	PageCount        types.IntString `json:"pagecount,omitempty"`
	RatingID         types.IntString `json:"rating_id,omitempty"`
	RatingName       string          `json:"rating_name,omitempty"`
	FileURL                          // FileURL is the Full URL of the (SIZE) asset for the PRIMARY file of this submission. SIZE can be one of "full, screen, preview".
	Thumbs
	LatestThumbs
	SubmissionTypeID types.IntString `json:"submission_type_id,omitempty"`
	TypeName         string          `json:"type_name,omitempty"`
	Digitalsales     types.BooleanYN `json:"digitalsales,omitempty"`
	Printsales       types.BooleanYN `json:"printsales,omitempty"`
	FriendsOnly      types.BooleanYN `json:"friends_only,omitempty"`
	GuestBlock       types.BooleanYN `json:"guest_block,omitempty"`
	Scraps           types.BooleanYN `json:"scraps,omitempty"`
}

type UserIconURLs struct {
	Large  string `json:"user_icon_url_large,omitempty"`
	Medium string `json:"user_icon_url_medium,omitempty"`
	Small  string `json:"user_icon_url_small,omitempty"`
}

type SubmissionDetails struct {
	SubmissionBasic
	Keywords         []Keyword       `json:"keywords"`
	Favorite         types.BooleanYN `json:"favorite"`
	FavoritesCount   types.IntString `json:"favorites_count"`
	UserIconFileName string          `json:"user_icon_file_name"`
	UserIconURLs
	LatestFileURL
	Files                   []File             `json:"files"`
	Pools                   []Pool             `json:"pools"`
	Description             string             `json:"description"`
	DescriptionBBCodeParsed string             `json:"description_bbcode_parsed"`
	Writing                 string             `json:"writing"`
	WritingBBCodeParsed     string             `json:"writing_bbcode_parsed"`
	PoolsCount              int                `json:"pools_count"`
	Ratings                 []SubmissionRating `json:"ratings"`
	CommentsCount           types.IntString    `json:"comments_count"`
	Views                   types.IntString    `json:"views"`
	SalesDescription        string             `json:"sales_description"`
	ForSale                 types.BooleanYN    `json:"forsale"`
	DigitalPrice            string             `json:"digital_price"`
	Prints                  []Print            `json:"prints"`
}

type Keyword struct {
	KeywordID   types.IntString `json:"keyword_id"`
	KeywordName string          `json:"keyword_name"`
	Suggested   types.BooleanYN `json:"contributed"`
	Count       types.IntString `json:"submissions_count"`
}

type File struct {
	FileID   types.IntString `json:"file_id"`
	FileName string          `json:"file_name"`
	Thumbs
	FileURL                             // Full URL of the (SIZE) asset for this file. SIZE can be one of "full, screen, preview".
	MimeType            string          `json:"mimetype"`
	SubmissionID        types.IntString `json:"submission_id"`
	UserID              types.IntString `json:"user_id"`
	SubmissionFileOrder types.IntString `json:"submission_file_order"` // An integer showing the order in which the files attached to this submission should be displayed. Starts counting at 0 for the first file/page in the submission.
	FileDimensions
	FileMD5
	Deleted            types.BooleanYN `json:"deleted"`
	CreateDateTime     string          `json:"create_datetime"`
	CreateDateTimeUser string          `json:"create_datetime_usertime"`
}

type FileDimensions struct {
	FullSizeX    types.IntString `json:"full_size_x"`
	FullSizeY    types.IntString `json:"full_size_y"`
	ScreenSizeX  types.IntString `json:"screen_size_x"`
	ScreenSizeY  types.IntString `json:"screen_size_y"`
	PreviewSizeX types.IntString `json:"preview_size_x"`
	PreviewSizeY types.IntString `json:"preview_size_y"`
}

type FileMD5 struct {
	InitialFileMD5 string `json:"initial_file_md5"`
	FullFileMD5    string `json:"full_file_md5"`
	LargeFileMD5   string `json:"large_file_md5"`
	SmallFileMD5   string `json:"small_file_md5"`
	ThumbnailMD5   string `json:"thumbnail_md5"`
}

type FileURL struct {
	FileURLFull    string `json:"file_url_full,omitempty"`
	FileURLScreen  string `json:"file_url_screen,omitempty"`
	FileURLPreview string `json:"file_url_preview,omitempty"`
}

type Pool struct {
	PoolID                     types.IntString `json:"pool_id"`
	Name                       string          `json:"name"`
	Description                string          `json:"description"`
	Count                      types.IntString `json:"count"`
	LeftSubmissionID           types.IntString `json:"submission_left_submission_id"`
	RightSubmissionID          types.IntString `json:"submission_right_submission_id"`
	LeftSubmissionFileName     string          `json:"submission_left_file_name"`
	RightSubmissionFileName    string          `json:"submission_right_file_name"`
	LeftThumbnailURL           string          `json:"submission_left_thumbnail_url,omitempty"`
	RightThumbnailURL          string          `json:"submission_right_thumbnail_url,omitempty"`
	LeftThumbnailURLNonCustom  string          `json:"submission_left_thumbnail_url_noncustom,omitempty"`
	RightThumbnailURLNonCustom string          `json:"submission_right_thumbnail_url_noncustom,omitempty"`
	LeftThumbX                 types.IntString `json:"submission_left_thumb_huge_x,omitempty"`
	LeftThumbY                 types.IntString `json:"submission_left_thumb_huge_y,omitempty"`
	RightThumbX                types.IntString `json:"submission_right_thumb_huge_x,omitempty"`
	RightThumbY                types.IntString `json:"submission_right_thumb_huge_y,omitempty"`
	LeftThumbNonCustomX        types.IntString `json:"submission_left_thumb_huge_noncustom_x,omitempty"`
	LeftThumbNonCustomY        types.IntString `json:"submission_left_thumb_huge_noncustom_y,omitempty"`
	RightThumbNonCustomX       types.IntString `json:"submission_right_thumb_huge_noncustom_x,omitempty"`
	RightThumbNonCustomY       types.IntString `json:"submission_right_thumb_huge_noncustom_y,omitempty"`
}

type Print struct {
	PrintSizeID        types.IntString   `json:"print_size_id"`
	Name               string            `json:"name"`
	Price              types.PriceString `json:"price"`
	PriceOwnerDiscount types.PriceString `json:"price_owner_discount,omitempty"`
}

type SubmissionRating struct {
	ContentTagID types.IntString `json:"content_tag_id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	RatingID     types.IntString `json:"rating_id"`
}

// LatestFileURL Full URL of the (SIZE) asset for the LATEST added file of this submission. SIZE can be one of "full, screen, preview".
type LatestFileURL struct {
	LatestFileURLFull    string `json:"latest_file_url_full"`
	LatestFileURLScreen  string `json:"latest_file_url_screen"`
	LatestFileURLPreview string `json:"latest_file_url_preview"`
}

type LatestThumbs struct {
	LatestThumbnailURLMedium          string `json:"latest_thumbnail_url_medium,omitempty"`
	LatestThumbnailURLLarge           string `json:"latest_thumbnail_url_large,omitempty"`
	LatestThumbnailURLHuge            string `json:"latest_thumbnail_url_huge,omitempty"`
	LatestThumbnailURLMediumNonCustom string `json:"latest_thumbnail_url_medium_noncustom,omitempty"`
	LatestThumbnailURLLargeNonCustom  string `json:"latest_thumbnail_url_large_noncustom,omitempty"`
	LatestThumbnailURLHugeNonCustom   string `json:"latest_thumbnail_url_huge_noncustom,omitempty"`

	LatestThumbMediumX          types.IntString `json:"latest_thumb_medium_x,omitempty"`
	LatestThumbMediumY          types.IntString `json:"latest_thumb_medium_y,omitempty"`
	LatestThumbLargeX           types.IntString `json:"latest_thumb_large_x,omitempty"`
	LatestThumbLargeY           types.IntString `json:"latest_thumb_large_y,omitempty"`
	LatestThumbHugeX            types.IntString `json:"latest_thumb_huge_x,omitempty"`
	LatestThumbHugeY            types.IntString `json:"latest_thumb_huge_y,omitempty"`
	LatestThumbMediumNonCustomX types.IntString `json:"latest_thumb_medium_noncustom_x,omitempty"`
	LatestThumbMediumNonCustomY types.IntString `json:"latest_thumb_medium_noncustom_y,omitempty"`
	LatestThumbLargeNonCustomX  types.IntString `json:"latest_thumb_large_noncustom_x,omitempty"`
	LatestThumbLargeNonCustomY  types.IntString `json:"latest_thumb_large_noncustom_y,omitempty"`
	LatestThumbHugeNonCustomX   types.IntString `json:"latest_thumb_huge_noncustom_x,omitempty"`
	LatestThumbHugeNonCustomY   types.IntString `json:"latest_thumb_huge_noncustom_y,omitempty"`
}

type Thumbs struct {
	ThumbnailURLMedium          string `json:"thumbnail_url_medium,omitempty"`
	ThumbnailURLLarge           string `json:"thumbnail_url_large,omitempty"`
	ThumbnailURLHuge            string `json:"thumbnail_url_huge,omitempty"`
	ThumbnailURLMediumNonCustom string `json:"thumbnail_url_medium_noncustom,omitempty"`
	ThumbnailURLLargeNonCustom  string `json:"thumbnail_url_large_noncustom,omitempty"`
	ThumbnailURLHugeNonCustom   string `json:"thumbnail_url_huge_noncustom,omitempty"`

	ThumbMediumX          types.IntString `json:"thumb_medium_x,omitempty"`
	ThumbMediumY          types.IntString `json:"thumb_medium_y,omitempty"`
	ThumbLargeX           types.IntString `json:"thumb_large_x,omitempty"`
	ThumbLargeY           types.IntString `json:"thumb_large_y,omitempty"`
	ThumbHugeX            types.IntString `json:"thumb_huge_x,omitempty"`
	ThumbHugeY            types.IntString `json:"thumb_huge_y,omitempty"`
	ThumbMediumNonCustomX types.IntString `json:"thumb_medium_noncustom_x,omitempty"`
	ThumbMediumNonCustomY types.IntString `json:"thumb_medium_noncustom_y,omitempty"`
	ThumbLargeNonCustomX  types.IntString `json:"thumb_large_noncustom_x,omitempty"`
	ThumbLargeNonCustomY  types.IntString `json:"thumb_large_noncustom_y,omitempty"`
	ThumbHugeNonCustomX   types.IntString `json:"thumb_huge_noncustom_x,omitempty"`
	ThumbHugeNonCustomY   types.IntString `json:"thumb_huge_noncustom_y,omitempty"`
}
type SubmissionDetailsResponse struct {
	SID          string              `json:"sid"`
	ResultsCount types.IntString     `json:"results_count"`
	UserLocation string              `json:"user_location"`
	Submissions  []SubmissionDetails `json:"submissions"`
}

type SubmissionFavoritesResponse struct {
	Sid   string             `json:"sid"`
	Users []types.UsernameID `json:"favingusers"`
}

func (u *User) SubmissionDetails(req SubmissionDetailsRequest) (SubmissionDetailsResponse, error) {
	if req.SID == "" {
		if u.SID == "" {
			return SubmissionDetailsResponse{}, ErrNotLoggedIn
		}
		req.SID = u.SID
	}
	return u.Client().SubmissionDetails(req)
}

func (c *Client) SubmissionDetails(req SubmissionDetailsRequest) (SubmissionDetailsResponse, error) {
	if req.SID == "" {
		return SubmissionDetailsResponse{}, ErrNotLoggedIn
	}
	if len(req.SubmissionIDSlice) > 0 {
		if req.SubmissionIDs != "" {
			req.SubmissionIDs += ","
		}
		req.SubmissionIDs += strings.Join(req.SubmissionIDSlice, ",")
		req.SubmissionIDSlice = nil
	}
	return PostDecode[SubmissionDetailsResponse](c, ApiUrl("submissions"), req)
}

func GetSubmissionDetails(req SubmissionDetailsRequest) (SubmissionDetailsResponse, error) {
	return DefaultClient.SubmissionDetails(req)
}

func (u *User) SubmissionFavorites(id types.IntString) (SubmissionFavoritesResponse, error) {
	if u.SID == "" {
		return SubmissionFavoritesResponse{}, ErrNotLoggedIn
	}
	val := url.Values{"sid": {u.SID}, "submission_id": {id.String()}}
	return PostDecode[SubmissionFavoritesResponse](u.Client(), ApiUrl("submissionfavingusers"), val)
}
