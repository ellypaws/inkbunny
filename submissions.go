package inkbunny

import (
	"context"
	"net/url"
	"strings"
)

// SubmissionDetailsRequest is modified to use BooleanYN for fields requiring "yes" or "no" representation.
type SubmissionDetailsRequest struct {
	SID                         string     `json:"sid,omitempty" query:"sid"`
	SubmissionIDs               string     `json:"submission_ids" query:"submission_ids"` // SubmissionIDs is a comma-separated list of submission IDs
	SubmissionIDSlice           []string   `json:"-"`                                     // SubmissionIDSlice will be joined as a comma-separated into SubmissionIDs
	OutputMode                  OutputMode `json:"output_mode,omitempty" query:"output_mode"`
	SortKeywordsBy              string     `json:"sort_keywords_by,omitempty" query:"sort_keywords_by"`
	ShowDescription             BooleanYN  `json:"show_description,omitempty" query:"show_description"`
	ShowDescriptionBbcodeParsed BooleanYN  `json:"show_description_bbcode_parsed,omitempty" query:"show_description_bbcode_parsed"`
	ShowWriting                 BooleanYN  `json:"show_writing,omitempty" query:"show_writing"`
	ShowWritingBbcodeParsed     BooleanYN  `json:"show_writing_bbcode_parsed,omitempty" query:"show_writing_bbcode_parsed"`
	ShowPools                   BooleanYN  `json:"show_pools,omitempty" query:"show_pools"`
}

// SubmissionBasic combines elements common in SubmissionSearch and SubmissionDetails
type SubmissionBasic struct {
	SubmissionID     IntString   `json:"submission_id"`
	Hidden           BooleanYN   `json:"hidden,omitempty"`
	Username         string      `json:"username,omitempty"`
	UserID           IntString   `json:"user_id,omitempty"`
	CreateDateSystem string      `json:"create_datetime,omitempty"`
	CreateDateUser   string      `json:"create_datetime_usertime,omitempty"`
	UpdateDateSystem string      `json:"last_file_update_datetime,omitempty"`
	UpdateDateUser   string      `json:"last_file_update_datetime_usertime,omitempty"`
	FileName         FalsyString `json:"file_name,omitempty"`
	LatestFileName   FalsyString `json:"latest_file_name,omitempty"`
	Title            string      `json:"title,omitempty"`
	Deleted          BooleanYN   `json:"deleted,omitempty"`
	Public           BooleanYN   `json:"public,omitempty"`
	MimeType         string      `json:"mimetype,omitempty"`
	LatestMimeType   string      `json:"latest_mimetype,omitempty"`
	PageCount        IntString   `json:"pagecount,omitempty"`
	RatingID         IntString   `json:"rating_id,omitempty"`
	RatingName       string      `json:"rating_name,omitempty"`
	FileURL                      // FileURL is the Full URL of the (SIZE) asset for the PRIMARY file of this submission. SIZE can be one of "full, screen, preview".
	Thumbs
	LatestThumbs
	SubmissionTypeID IntString `json:"submission_type_id,omitempty"`
	TypeName         string    `json:"type_name,omitempty"`
	Digitalsales     BooleanYN `json:"digitalsales,omitempty"`
	Printsales       BooleanYN `json:"printsales,omitempty"`
	FriendsOnly      BooleanYN `json:"friends_only,omitempty"`
	GuestBlock       BooleanYN `json:"guest_block,omitempty"`
	Scraps           BooleanYN `json:"scraps,omitempty"`
}

type UserIconURLs struct {
	Large  string `json:"user_icon_url_large,omitempty"`
	Medium string `json:"user_icon_url_medium,omitempty"`
	Small  string `json:"user_icon_url_small,omitempty"`
}

type SubmissionDetails struct {
	SubmissionBasic
	Keywords         []Keyword `json:"keywords"`
	Favorite         BooleanYN `json:"favorite"`
	FavoritesCount   IntString `json:"favorites_count"`
	UserIconFileName string    `json:"user_icon_file_name"`
	UserIconURLs
	LatestFileURL
	Files                   []File             `json:"files"`
	Pools                   []Pool             `json:"pools"`
	Description             string             `json:"description"`
	DescriptionBBCodeParsed string             `json:"description_bbcode_parsed"`
	Writing                 string             `json:"writing"`
	WritingBBCodeParsed     string             `json:"writing_bbcode_parsed"`
	PoolsCount              IntString          `json:"pools_count"`
	Ratings                 []SubmissionRating `json:"ratings"`
	CommentsCount           IntString          `json:"comments_count"`
	Views                   IntString          `json:"views"`
	SalesDescription        string             `json:"sales_description"`
	ForSale                 BooleanYN          `json:"forsale"`
	DigitalPrice            string             `json:"digital_price"`
	Prints                  []Print            `json:"prints"`
}

type Keyword struct {
	KeywordID   IntString `json:"keyword_id"`
	KeywordName string    `json:"keyword_name"`
	Suggested   BooleanYN `json:"contributed"`
	Count       IntString `json:"submissions_count"`
}

type File struct {
	FileID   IntString `json:"file_id"`
	FileName string    `json:"file_name"`
	Thumbs
	FileURL                       // Full URL of the (SIZE) asset for this file. SIZE can be one of "full, screen, preview".
	MimeType            string    `json:"mimetype"`
	SubmissionID        IntString `json:"submission_id"`
	UserID              IntString `json:"user_id"`
	SubmissionFileOrder IntString `json:"submission_file_order"` // An integer showing the order in which the files attached to this submission should be displayed. Starts counting at 0 for the first file/page in the submission.
	FileDimensions
	FileMD5
	Deleted            BooleanYN `json:"deleted"`
	CreateDateTime     string    `json:"create_datetime"`
	CreateDateTimeUser string    `json:"create_datetime_usertime"`
}

type FileDimensions struct {
	FullSizeX    IntString `json:"full_size_x"`
	FullSizeY    IntString `json:"full_size_y"`
	ScreenSizeX  IntString `json:"screen_size_x"`
	ScreenSizeY  IntString `json:"screen_size_y"`
	PreviewSizeX IntString `json:"preview_size_x"`
	PreviewSizeY IntString `json:"preview_size_y"`
}

type FileMD5 struct {
	InitialFileMD5 string `json:"initial_file_md5"`
	FullFileMD5    string `json:"full_file_md5"`
	LargeFileMD5   string `json:"large_file_md5"`
	SmallFileMD5   string `json:"small_file_md5"`
	ThumbnailMD5   string `json:"thumbnail_md5"`
}

type FileURL struct {
	FileURLFull    FalsyString `json:"file_url_full,omitempty"`
	FileURLScreen  FalsyString `json:"file_url_screen,omitempty"`
	FileURLPreview FalsyString `json:"file_url_preview,omitempty"`
}

type Pool struct {
	PoolID                           IntString `json:"pool_id"`
	Name                             string    `json:"name"`
	Description                      string    `json:"description"`
	Count                            IntString `json:"count"`
	LeftSubmissionID                 IntString `json:"submission_left_submission_id"`
	RightSubmissionID                IntString `json:"submission_right_submission_id"`
	LeftSubmissionFileName           string    `json:"submission_left_file_name"`
	RightSubmissionFileName          string    `json:"submission_right_file_name"`
	LeftThumbnailURLMedium           string    `json:"submission_left_thumbnail_url_medium,omitempty"`
	LeftThumbnailURLLarge            string    `json:"submission_left_thumbnail_url_large,omitempty"`
	LeftThumbnailURL                 string    `json:"submission_left_thumbnail_url_huge,omitempty"`
	RightThumbnailURLMedium          string    `json:"submission_right_thumbnail_url_medium,omitempty"`
	RightThumbnailURLLarge           string    `json:"submission_right_thumbnail_url_large,omitempty"`
	RightThumbnailURL                string    `json:"submission_right_thumbnail_url_huge,omitempty"`
	LeftThumbnailURLMediumNonCustom  string    `json:"submission_left_thumbnail_url_medium_noncustom,omitempty"`
	LeftThumbnailURLLargeNonCustom   string    `json:"submission_left_thumbnail_url_large_noncustom,omitempty"`
	LeftThumbnailURLNonCustom        string    `json:"submission_left_thumbnail_url_huge_noncustom,omitempty"`
	RightThumbnailURLMediumNonCustom string    `json:"submission_right_thumbnail_url_medium_noncustom,omitempty"`
	RightThumbnailURLLargeNonCustom  string    `json:"submission_right_thumbnail_url_large_noncustom,omitempty"`
	RightThumbnailURLNonCustom       string    `json:"submission_right_thumbnail_url_huge_noncustom,omitempty"`
	LeftThumbMediumX                 IntString `json:"submission_left_thumb_medium_x,omitempty"`
	LeftThumbMediumY                 IntString `json:"submission_left_thumb_medium_y,omitempty"`
	LeftThumbLargeX                  IntString `json:"submission_left_thumb_large_x,omitempty"`
	LeftThumbLargeY                  IntString `json:"submission_left_thumb_large_y,omitempty"`
	LeftThumbX                       IntString `json:"submission_left_thumb_huge_x,omitempty"`
	LeftThumbY                       IntString `json:"submission_left_thumb_huge_y,omitempty"`
	RightThumbMediumX                IntString `json:"submission_right_thumb_medium_x,omitempty"`
	RightThumbMediumY                IntString `json:"submission_right_thumb_medium_y,omitempty"`
	RightThumbLargeX                 IntString `json:"submission_right_thumb_large_x,omitempty"`
	RightThumbLargeY                 IntString `json:"submission_right_thumb_large_y,omitempty"`
	RightThumbX                      IntString `json:"submission_right_thumb_huge_x,omitempty"`
	RightThumbY                      IntString `json:"submission_right_thumb_huge_y,omitempty"`
	LeftThumbMediumNonCustomX        IntString `json:"submission_left_thumb_medium_noncustom_x,omitempty"`
	LeftThumbMediumNonCustomY        IntString `json:"submission_left_thumb_medium_noncustom_y,omitempty"`
	LeftThumbLargeNonCustomX         IntString `json:"submission_left_thumb_large_noncustom_x,omitempty"`
	LeftThumbLargeNonCustomY         IntString `json:"submission_left_thumb_large_noncustom_y,omitempty"`
	LeftThumbNonCustomX              IntString `json:"submission_left_thumb_huge_noncustom_x,omitempty"`
	LeftThumbNonCustomY              IntString `json:"submission_left_thumb_huge_noncustom_y,omitempty"`
	RightThumbMediumNonCustomX       IntString `json:"submission_right_thumb_medium_noncustom_x,omitempty"`
	RightThumbMediumNonCustomY       IntString `json:"submission_right_thumb_medium_noncustom_y,omitempty"`
	RightThumbLargeNonCustomX        IntString `json:"submission_right_thumb_large_noncustom_x,omitempty"`
	RightThumbLargeNonCustomY        IntString `json:"submission_right_thumb_large_noncustom_y,omitempty"`
	RightThumbNonCustomX             IntString `json:"submission_right_thumb_huge_noncustom_x,omitempty"`
	RightThumbNonCustomY             IntString `json:"submission_right_thumb_huge_noncustom_y,omitempty"`
}

type Print struct {
	PrintSizeID        IntString   `json:"print_size_id"`
	Name               string      `json:"name"`
	Price              PriceString `json:"price"`
	PriceOwnerDiscount PriceString `json:"price_owner_discount,omitempty"`
}

type SubmissionRating struct {
	ContentTagID IntString `json:"content_tag_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	RatingID     IntString `json:"rating_id"`
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

	LatestThumbMediumX          IntString `json:"latest_thumb_medium_x,omitempty"`
	LatestThumbMediumY          IntString `json:"latest_thumb_medium_y,omitempty"`
	LatestThumbLargeX           IntString `json:"latest_thumb_large_x,omitempty"`
	LatestThumbLargeY           IntString `json:"latest_thumb_large_y,omitempty"`
	LatestThumbHugeX            IntString `json:"latest_thumb_huge_x,omitempty"`
	LatestThumbHugeY            IntString `json:"latest_thumb_huge_y,omitempty"`
	LatestThumbMediumNonCustomX IntString `json:"latest_thumb_medium_noncustom_x,omitempty"`
	LatestThumbMediumNonCustomY IntString `json:"latest_thumb_medium_noncustom_y,omitempty"`
	LatestThumbLargeNonCustomX  IntString `json:"latest_thumb_large_noncustom_x,omitempty"`
	LatestThumbLargeNonCustomY  IntString `json:"latest_thumb_large_noncustom_y,omitempty"`
	LatestThumbHugeNonCustomX   IntString `json:"latest_thumb_huge_noncustom_x,omitempty"`
	LatestThumbHugeNonCustomY   IntString `json:"latest_thumb_huge_noncustom_y,omitempty"`
}

type Thumbs struct {
	ThumbnailURLMedium          string `json:"thumbnail_url_medium,omitempty"`
	ThumbnailURLLarge           string `json:"thumbnail_url_large,omitempty"`
	ThumbnailURLHuge            string `json:"thumbnail_url_huge,omitempty"`
	ThumbnailURLMediumNonCustom string `json:"thumbnail_url_medium_noncustom,omitempty"`
	ThumbnailURLLargeNonCustom  string `json:"thumbnail_url_large_noncustom,omitempty"`
	ThumbnailURLHugeNonCustom   string `json:"thumbnail_url_huge_noncustom,omitempty"`

	ThumbMediumX          IntString `json:"thumb_medium_x,omitempty"`
	ThumbMediumY          IntString `json:"thumb_medium_y,omitempty"`
	ThumbLargeX           IntString `json:"thumb_large_x,omitempty"`
	ThumbLargeY           IntString `json:"thumb_large_y,omitempty"`
	ThumbHugeX            IntString `json:"thumb_huge_x,omitempty"`
	ThumbHugeY            IntString `json:"thumb_huge_y,omitempty"`
	ThumbMediumNonCustomX IntString `json:"thumb_medium_noncustom_x,omitempty"`
	ThumbMediumNonCustomY IntString `json:"thumb_medium_noncustom_y,omitempty"`
	ThumbLargeNonCustomX  IntString `json:"thumb_large_noncustom_x,omitempty"`
	ThumbLargeNonCustomY  IntString `json:"thumb_large_noncustom_y,omitempty"`
	ThumbHugeNonCustomX   IntString `json:"thumb_huge_noncustom_x,omitempty"`
	ThumbHugeNonCustomY   IntString `json:"thumb_huge_noncustom_y,omitempty"`
}
type SubmissionDetailsResponse struct {
	SID          string              `json:"sid"`
	ResultsCount IntString           `json:"results_count"`
	UserLocation string              `json:"user_location"`
	Submissions  []SubmissionDetails `json:"submissions"`
}

type SubmissionFavoritesResponse struct {
	Sid   string       `json:"sid"`
	Users []UsernameID `json:"favingusers"`
}

func (u *User) SubmissionDetails(req SubmissionDetailsRequest) (SubmissionDetailsResponse, error) {
	return u.SubmissionDetailsContext(context.Background(), req)
}

// SubmissionDetailsContext is like SubmissionDetails but accepts a context.Context
// for per-call cancellation and timeout control.
func (u *User) SubmissionDetailsContext(ctx context.Context, req SubmissionDetailsRequest) (SubmissionDetailsResponse, error) {
	if req.SID == "" {
		if u.SID == "" {
			return SubmissionDetailsResponse{}, ErrNotLoggedIn
		}
		req.SID = u.SID
	}
	return u.Client().SubmissionDetailsContext(ctx, req)
}

func (c *Client) SubmissionDetails(req SubmissionDetailsRequest) (SubmissionDetailsResponse, error) {
	return c.SubmissionDetailsContext(c.ctx, req)
}

// SubmissionDetailsContext is like SubmissionDetails but accepts a context.Context
// for per-call cancellation and timeout control.
func (c *Client) SubmissionDetailsContext(ctx context.Context, req SubmissionDetailsRequest) (SubmissionDetailsResponse, error) {
	if req.SID == "" {
		return SubmissionDetailsResponse{}, ErrEmptySID
	}
	if len(req.SubmissionIDSlice) > 0 {
		if req.SubmissionIDs != "" {
			req.SubmissionIDs += ","
		}
		req.SubmissionIDs += strings.Join(req.SubmissionIDSlice, ",")
		req.SubmissionIDSlice = nil
	}
	return PostDecode[SubmissionDetailsResponse](c.withContext(ctx), ApiUrl("submissions"), req)
}

func GetSubmissionDetails(req SubmissionDetailsRequest) (SubmissionDetailsResponse, error) {
	return DefaultClient.SubmissionDetails(req)
}

// SubmissionFavorites retrieves the list of users who have favorited a specific submission.
func (u *User) SubmissionFavorites(id IntString) (SubmissionFavoritesResponse, error) {
	if u.SID == "" {
		return SubmissionFavoritesResponse{}, ErrNotLoggedIn
	}
	val := url.Values{"sid": {u.SID}, "submission_id": {id.String()}}
	return PostDecode[SubmissionFavoritesResponse](u.Client(), ApiUrl("submissionfavingusers"), val)
}
