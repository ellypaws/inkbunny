package inkbunny

import (
	"context"
	"net/url"
	"strings"
)

// SubmissionDetailsRequest configures Client.SubmissionDetails.
//
// This package uses OutputMode for format selection and BooleanYN for yes or no
// flags. Leaving optional fields at their zero value lets Inkbunny apply its
// documented defaults: JSON, "alphabetical", No, No, No, No, and No.
type SubmissionDetailsRequest struct {
	// SID overrides the client's current session for this call.
	SID string `json:"sid,omitempty" query:"sid"`
	// SubmissionIDs is the raw comma-separated form accepted by Inkbunny.
	SubmissionIDs string `json:"submission_ids" query:"submission_ids"` // SubmissionIDs is a comma-separated list of submission IDs
	// SubmissionIDSlice is the package-friendly alternative to SubmissionIDs.
	// Any values here are joined and appended before the request is sent.
	SubmissionIDSlice []string `json:"-"` // SubmissionIDSlice will be joined as a comma-separated into SubmissionIDs
	// OutputMode selects the response format. Use JSON or XML.
	// Defaults to JSON.
	OutputMode OutputMode `json:"output_mode,omitempty" query:"output_mode"`
	// SortKeywordsBy controls keyword ordering. Common values are
	// "alphabetical" and "submissions_count". Defaults to "alphabetical".
	SortKeywordsBy string `json:"sort_keywords_by,omitempty" query:"sort_keywords_by"`
	// ShowDescription enables the raw description field. Use Yes or No.
	// Defaults to No.
	ShowDescription BooleanYN `json:"show_description,omitempty" query:"show_description"`
	// ShowDescriptionBbcodeParsed enables the parsed HTML description. Use Yes or No.
	// Defaults to No.
	ShowDescriptionBbcodeParsed BooleanYN `json:"show_description_bbcode_parsed,omitempty" query:"show_description_bbcode_parsed"`
	// ShowWriting enables the raw writing or story field. Use Yes or No.
	// Defaults to No.
	ShowWriting BooleanYN `json:"show_writing,omitempty" query:"show_writing"`
	// ShowWritingBbcodeParsed enables the parsed HTML writing field. Use Yes or No.
	// Defaults to No.
	ShowWritingBbcodeParsed BooleanYN `json:"show_writing_bbcode_parsed,omitempty" query:"show_writing_bbcode_parsed"`
	// ShowPools includes pool membership details. Use Yes or No.
	// Defaults to No.
	ShowPools BooleanYN `json:"show_pools,omitempty" query:"show_pools"`
}

// SubmissionBasic contains the shared submission fields returned by search and
// submission-details calls.
//
// Numeric values are exposed as IntString because Inkbunny often encodes numbers
// as strings. Boolean values use BooleanYN so the same type works for both
// request and response payloads.
type SubmissionBasic struct {
	// SubmissionID is the submission ID.
	SubmissionID IntString `json:"submission_id"`
	// Hidden reports whether the submission is blocked by the caller's filters.
	Hidden BooleanYN `json:"hidden,omitempty"`
	// Username is the submission owner's username.
	Username string `json:"username,omitempty"`
	// UserID is the submission owner's user ID.
	UserID IntString `json:"user_id,omitempty"`
	// CreateDateSystem is the upload or publish time in system time.
	CreateDateSystem string `json:"create_datetime,omitempty"`
	// CreateDateUser is the upload or publish time in the user's local time.
	CreateDateUser string `json:"create_datetime_usertime,omitempty"`
	// UpdateDateSystem is the latest file update time in system time.
	UpdateDateSystem string `json:"last_file_update_datetime,omitempty"`
	// UpdateDateUser is the latest file update time in the user's local time.
	UpdateDateUser string `json:"last_file_update_datetime_usertime,omitempty"`
	// FileName is the primary file name.
	FileName FalsyString `json:"file_name,omitempty"`
	// LatestFileName is the latest file name added to the submission.
	LatestFileName FalsyString `json:"latest_file_name,omitempty"`
	// Title is the submission title.
	Title string `json:"title,omitempty"`
	// Deleted reports whether the submission is deleted.
	Deleted BooleanYN `json:"deleted,omitempty"`
	// Public reports whether the submission is public.
	Public BooleanYN `json:"public,omitempty"`
	// MimeType is the primary file MIME type.
	MimeType string `json:"mimetype,omitempty"`
	// LatestMimeType is the latest file MIME type.
	LatestMimeType string `json:"latest_mimetype,omitempty"`
	// PageCount is the number of files attached to the submission.
	PageCount IntString `json:"pagecount,omitempty"`
	// RatingID is the overall submission rating ID.
	RatingID IntString `json:"rating_id,omitempty"`
	// RatingName is the overall submission rating name.
	RatingName string `json:"rating_name,omitempty"`
	// FileURL contains URLs for the primary file assets.
	FileURL // FileURL is the Full URL of the (SIZE) asset for the PRIMARY file of this submission. SIZE can be one of "full, screen, preview".
	// Thumbs contains dimensions and URLs for the primary file thumbnails.
	Thumbs
	// LatestThumbs contains dimensions and URLs for the latest file thumbnails.
	LatestThumbs
	// SubmissionTypeID is the numeric submission type.
	SubmissionTypeID IntString `json:"submission_type_id,omitempty"`
	// TypeName is the submission type name.
	TypeName string `json:"type_name,omitempty"`
	// Digitalsales reports whether digital sales are enabled.
	Digitalsales BooleanYN `json:"digitalsales,omitempty"`
	// Printsales reports whether print sales are enabled.
	Printsales BooleanYN `json:"printsales,omitempty"`
	// FriendsOnly reports whether the submission is visible only to friends.
	FriendsOnly BooleanYN `json:"friends_only,omitempty"`
	// GuestBlock reports whether guest access is blocked.
	GuestBlock BooleanYN `json:"guest_block,omitempty"`
	// Scraps reports whether the submission is in the scraps gallery.
	Scraps BooleanYN `json:"scraps,omitempty"`
}

// UserIconURLs groups user icon URLs by size.
type UserIconURLs struct {
	// Large is the large user icon URL.
	Large string `json:"user_icon_url_large,omitempty"`
	// Medium is the medium user icon URL.
	Medium string `json:"user_icon_url_medium,omitempty"`
	// Small is the small user icon URL.
	Small string `json:"user_icon_url_small,omitempty"`
}

// SubmissionDetails is the expanded submission model returned by
// Client.SubmissionDetails.
type SubmissionDetails struct {
	// SubmissionBasic contains the fields shared with search results.
	SubmissionBasic
	// Keywords lists the keywords assigned to the submission.
	Keywords []Keyword `json:"keywords"`
	// Favorite reports whether the current session has favorited the submission.
	Favorite BooleanYN `json:"favorite"`
	// FavoritesCount is the total number of favorites on the submission.
	FavoritesCount IntString `json:"favorites_count"`
	// UserIconFileName is the owning user's user-icon file name.
	UserIconFileName string `json:"user_icon_file_name"`
	// UserIconURLs contains the owning user's icon URLs.
	UserIconURLs
	// LatestFileURL contains URLs for the latest file assets.
	LatestFileURL
	// Files lists all files attached to the submission.
	Files []File `json:"files"`
	// Pools lists the pools that contain the submission.
	Pools []Pool `json:"pools"`
	// Description is the raw submission description.
	Description string `json:"description"`
	// DescriptionBBCodeParsed is the description rendered to HTML.
	DescriptionBBCodeParsed string `json:"description_bbcode_parsed"`
	// Writing is the raw writing or story text.
	Writing string `json:"writing"`
	// WritingBBCodeParsed is the writing rendered to HTML.
	WritingBBCodeParsed string `json:"writing_bbcode_parsed"`
	// PoolsCount is the number of pools that contain the submission.
	PoolsCount IntString `json:"pools_count"`
	// Ratings lists the fine-grained content tags that produced RatingID.
	Ratings []SubmissionRating `json:"ratings"`
	// CommentsCount is the total number of comments on the submission.
	CommentsCount IntString `json:"comments_count"`
	// Views is the total number of recorded views.
	Views IntString `json:"views"`
	// SalesDescription is the short digital-sales description.
	SalesDescription string `json:"sales_description"`
	// ForSale reports whether any sale mode is active.
	ForSale BooleanYN `json:"forsale"`
	// DigitalPrice is the digital-sale price in USD.
	DigitalPrice string `json:"digital_price"`
	// Prints lists available print-sale sizes.
	Prints []Print `json:"prints"`
}

// Keyword describes one submission keyword entry.
type Keyword struct {
	// KeywordID is the keyword ID.
	KeywordID IntString `json:"keyword_id"`
	// KeywordName is the keyword text.
	KeywordName string `json:"keyword_name"`
	// Suggested reports whether the keyword is suggested rather than owner-assigned.
	Suggested BooleanYN `json:"contributed"`
	// Count is the number of submissions that match the keyword.
	Count IntString `json:"submissions_count"`
}

// File describes one attached submission file.
type File struct {
	// FileID is the file ID.
	FileID IntString `json:"file_id"`
	// FileName is the base file name.
	FileName string `json:"file_name"`
	// Thumbs contains dimensions and URLs for this file's thumbnails.
	Thumbs
	// FileURL contains URLs for this file's asset variants.
	FileURL // Full URL of the (SIZE) asset for this file. SIZE can be one of "full, screen, preview".
	// MimeType is the file MIME type.
	MimeType string `json:"mimetype"`
	// SubmissionID is the submission that owns the file.
	SubmissionID IntString `json:"submission_id"`
	// UserID is the user who owns the file.
	UserID IntString `json:"user_id"`
	// SubmissionFileOrder is the display order of the file within the submission.
	SubmissionFileOrder IntString `json:"submission_file_order"` // An integer showing the order in which the files attached to this submission should be displayed. Starts counting at 0 for the first file/page in the submission.
	// FileDimensions contains image dimensions for the file variants.
	FileDimensions
	// FileMD5 contains MD5 checksums for the file variants.
	FileMD5
	// Deleted reports whether the file is deleted.
	Deleted BooleanYN `json:"deleted"`
	// CreateDateTime is the upload time in system time.
	CreateDateTime string `json:"create_datetime"`
	// CreateDateTimeUser is the upload time in the user's local time.
	CreateDateTimeUser string `json:"create_datetime_usertime"`
}

// FileDimensions contains pixel dimensions for each file variant returned by Inkbunny.
type FileDimensions struct {
	// FullSizeX is the full-size width in pixels.
	FullSizeX IntString `json:"full_size_x"`
	// FullSizeY is the full-size height in pixels.
	FullSizeY IntString `json:"full_size_y"`
	// ScreenSizeX is the screen-size width in pixels.
	ScreenSizeX IntString `json:"screen_size_x"`
	// ScreenSizeY is the screen-size height in pixels.
	ScreenSizeY IntString `json:"screen_size_y"`
	// PreviewSizeX is the preview-size width in pixels.
	PreviewSizeX IntString `json:"preview_size_x"`
	// PreviewSizeY is the preview-size height in pixels.
	PreviewSizeY IntString `json:"preview_size_y"`
}

// FileMD5 contains MD5 checksums for the returned file variants.
type FileMD5 struct {
	// InitialFileMD5 is the MD5 of the originally uploaded file.
	InitialFileMD5 string `json:"initial_file_md5"`
	// FullFileMD5 is the MD5 of the full-size file.
	FullFileMD5 string `json:"full_file_md5"`
	// LargeFileMD5 is the MD5 of the screen-size file.
	LargeFileMD5 string `json:"large_file_md5"`
	// SmallFileMD5 is the MD5 of the preview-size file.
	SmallFileMD5 string `json:"small_file_md5"`
	// ThumbnailMD5 is the MD5 of the custom thumbnail, when present.
	ThumbnailMD5 string `json:"thumbnail_md5"`
}

// FileURL contains URLs for the primary file variants used throughout the package.
type FileURL struct {
	// FileURLFull is the full-size asset URL.
	FileURLFull FalsyString `json:"file_url_full,omitempty"`
	// FileURLScreen is the screen-size asset URL.
	FileURLScreen FalsyString `json:"file_url_screen,omitempty"`
	// FileURLPreview is the preview-size asset URL.
	FileURLPreview FalsyString `json:"file_url_preview,omitempty"`
}

// Pool describes one related pool entry returned by SubmissionDetails.
type Pool struct {
	// PoolID is the pool ID.
	PoolID IntString `json:"pool_id"`
	// Name is the pool name.
	Name string `json:"name"`
	// Description is the pool description.
	Description string `json:"description"`
	// Count is the total number of submissions in the pool.
	Count IntString `json:"count"`
	// LeftSubmissionID is the submission to the left of the current one in the pool.
	LeftSubmissionID IntString `json:"submission_left_submission_id"`
	// RightSubmissionID is the submission to the right of the current one in the pool.
	RightSubmissionID IntString `json:"submission_right_submission_id"`
	// LeftSubmissionFileName is the left submission's primary file name.
	LeftSubmissionFileName string `json:"submission_left_file_name"`
	// RightSubmissionFileName is the right submission's primary file name.
	RightSubmissionFileName string `json:"submission_right_file_name"`
	// LeftThumbnailURLMedium is the left submission's medium thumbnail URL.
	LeftThumbnailURLMedium string `json:"submission_left_thumbnail_url_medium,omitempty"`
	// LeftThumbnailURLLarge is the left submission's large thumbnail URL.
	LeftThumbnailURLLarge string `json:"submission_left_thumbnail_url_large,omitempty"`
	// LeftThumbnailURL is the left submission's huge thumbnail URL.
	LeftThumbnailURL string `json:"submission_left_thumbnail_url_huge,omitempty"`
	// RightThumbnailURLMedium is the right submission's medium thumbnail URL.
	RightThumbnailURLMedium string `json:"submission_right_thumbnail_url_medium,omitempty"`
	// RightThumbnailURLLarge is the right submission's large thumbnail URL.
	RightThumbnailURLLarge string `json:"submission_right_thumbnail_url_large,omitempty"`
	// RightThumbnailURL is the right submission's huge thumbnail URL.
	RightThumbnailURL string `json:"submission_right_thumbnail_url_huge,omitempty"`
	// LeftThumbnailURLMediumNonCustom is the left submission's medium non-custom thumbnail URL.
	LeftThumbnailURLMediumNonCustom string `json:"submission_left_thumbnail_url_medium_noncustom,omitempty"`
	// LeftThumbnailURLLargeNonCustom is the left submission's large non-custom thumbnail URL.
	LeftThumbnailURLLargeNonCustom string `json:"submission_left_thumbnail_url_large_noncustom,omitempty"`
	// LeftThumbnailURLNonCustom is the left submission's huge non-custom thumbnail URL.
	LeftThumbnailURLNonCustom string `json:"submission_left_thumbnail_url_huge_noncustom,omitempty"`
	// RightThumbnailURLMediumNonCustom is the right submission's medium non-custom thumbnail URL.
	RightThumbnailURLMediumNonCustom string `json:"submission_right_thumbnail_url_medium_noncustom,omitempty"`
	// RightThumbnailURLLargeNonCustom is the right submission's large non-custom thumbnail URL.
	RightThumbnailURLLargeNonCustom string `json:"submission_right_thumbnail_url_large_noncustom,omitempty"`
	// RightThumbnailURLNonCustom is the right submission's huge non-custom thumbnail URL.
	RightThumbnailURLNonCustom string `json:"submission_right_thumbnail_url_huge_noncustom,omitempty"`
	// LeftThumbMediumX is the left submission's medium thumbnail width.
	LeftThumbMediumX IntString `json:"submission_left_thumb_medium_x,omitempty"`
	// LeftThumbMediumY is the left submission's medium thumbnail height.
	LeftThumbMediumY IntString `json:"submission_left_thumb_medium_y,omitempty"`
	// LeftThumbLargeX is the left submission's large thumbnail width.
	LeftThumbLargeX IntString `json:"submission_left_thumb_large_x,omitempty"`
	// LeftThumbLargeY is the left submission's large thumbnail height.
	LeftThumbLargeY IntString `json:"submission_left_thumb_large_y,omitempty"`
	// LeftThumbX is the left submission's huge thumbnail width.
	LeftThumbX IntString `json:"submission_left_thumb_huge_x,omitempty"`
	// LeftThumbY is the left submission's huge thumbnail height.
	LeftThumbY IntString `json:"submission_left_thumb_huge_y,omitempty"`
	// RightThumbMediumX is the right submission's medium thumbnail width.
	RightThumbMediumX IntString `json:"submission_right_thumb_medium_x,omitempty"`
	// RightThumbMediumY is the right submission's medium thumbnail height.
	RightThumbMediumY IntString `json:"submission_right_thumb_medium_y,omitempty"`
	// RightThumbLargeX is the right submission's large thumbnail width.
	RightThumbLargeX IntString `json:"submission_right_thumb_large_x,omitempty"`
	// RightThumbLargeY is the right submission's large thumbnail height.
	RightThumbLargeY IntString `json:"submission_right_thumb_large_y,omitempty"`
	// RightThumbX is the right submission's huge thumbnail width.
	RightThumbX IntString `json:"submission_right_thumb_huge_x,omitempty"`
	// RightThumbY is the right submission's huge thumbnail height.
	RightThumbY IntString `json:"submission_right_thumb_huge_y,omitempty"`
	// LeftThumbMediumNonCustomX is the left submission's medium non-custom thumbnail width.
	LeftThumbMediumNonCustomX IntString `json:"submission_left_thumb_medium_noncustom_x,omitempty"`
	// LeftThumbMediumNonCustomY is the left submission's medium non-custom thumbnail height.
	LeftThumbMediumNonCustomY IntString `json:"submission_left_thumb_medium_noncustom_y,omitempty"`
	// LeftThumbLargeNonCustomX is the left submission's large non-custom thumbnail width.
	LeftThumbLargeNonCustomX IntString `json:"submission_left_thumb_large_noncustom_x,omitempty"`
	// LeftThumbLargeNonCustomY is the left submission's large non-custom thumbnail height.
	LeftThumbLargeNonCustomY IntString `json:"submission_left_thumb_large_noncustom_y,omitempty"`
	// LeftThumbNonCustomX is the left submission's huge non-custom thumbnail width.
	LeftThumbNonCustomX IntString `json:"submission_left_thumb_huge_noncustom_x,omitempty"`
	// LeftThumbNonCustomY is the left submission's huge non-custom thumbnail height.
	LeftThumbNonCustomY IntString `json:"submission_left_thumb_huge_noncustom_y,omitempty"`
	// RightThumbMediumNonCustomX is the right submission's medium non-custom thumbnail width.
	RightThumbMediumNonCustomX IntString `json:"submission_right_thumb_medium_noncustom_x,omitempty"`
	// RightThumbMediumNonCustomY is the right submission's medium non-custom thumbnail height.
	RightThumbMediumNonCustomY IntString `json:"submission_right_thumb_medium_noncustom_y,omitempty"`
	// RightThumbLargeNonCustomX is the right submission's large non-custom thumbnail width.
	RightThumbLargeNonCustomX IntString `json:"submission_right_thumb_large_noncustom_x,omitempty"`
	// RightThumbLargeNonCustomY is the right submission's large non-custom thumbnail height.
	RightThumbLargeNonCustomY IntString `json:"submission_right_thumb_large_noncustom_y,omitempty"`
	// RightThumbNonCustomX is the right submission's huge non-custom thumbnail width.
	RightThumbNonCustomX IntString `json:"submission_right_thumb_huge_noncustom_x,omitempty"`
	// RightThumbNonCustomY is the right submission's huge non-custom thumbnail height.
	RightThumbNonCustomY IntString `json:"submission_right_thumb_huge_noncustom_y,omitempty"`
}

// Print describes one print option returned with a submission.
type Print struct {
	// PrintSizeID is the numeric print size ID.
	PrintSizeID IntString `json:"print_size_id"`
	// Name is the print size label.
	Name string `json:"name"`
	// Price is the retail price in USD.
	Price PriceString `json:"price"`
	// PriceOwnerDiscount is the discounted price visible to the submission owner.
	PriceOwnerDiscount PriceString `json:"price_owner_discount,omitempty"`
}

// SubmissionRating describes one fine-grained rating tag assigned to a submission.
type SubmissionRating struct {
	// ContentTagID is the assigned content-tag ID.
	ContentTagID IntString `json:"content_tag_id"`
	// Name is the rating-tag name.
	Name string `json:"name"`
	// Description is the rating-tag description.
	Description string `json:"description"`
	// RatingID is the coarse rating bucket for the tag.
	RatingID IntString `json:"rating_id"`
}

// LatestFileURL contains URLs for the latest file variants.
type LatestFileURL struct {
	// LatestFileURLFull is the full-size asset URL for the latest file.
	LatestFileURLFull string `json:"latest_file_url_full"`
	// LatestFileURLScreen is the screen-size asset URL for the latest file.
	LatestFileURLScreen string `json:"latest_file_url_screen"`
	// LatestFileURLPreview is the preview-size asset URL for the latest file.
	LatestFileURLPreview string `json:"latest_file_url_preview"`
}

// LatestThumbs contains thumbnail URLs and dimensions for the latest file.
type LatestThumbs struct {
	// LatestThumbnailURLMedium is the medium thumbnail URL for the latest file.
	LatestThumbnailURLMedium string `json:"latest_thumbnail_url_medium,omitempty"`
	// LatestThumbnailURLLarge is the large thumbnail URL for the latest file.
	LatestThumbnailURLLarge string `json:"latest_thumbnail_url_large,omitempty"`
	// LatestThumbnailURLHuge is the huge thumbnail URL for the latest file.
	LatestThumbnailURLHuge string `json:"latest_thumbnail_url_huge,omitempty"`
	// LatestThumbnailURLMediumNonCustom is the medium non-custom thumbnail URL for the latest file.
	LatestThumbnailURLMediumNonCustom string `json:"latest_thumbnail_url_medium_noncustom,omitempty"`
	// LatestThumbnailURLLargeNonCustom is the large non-custom thumbnail URL for the latest file.
	LatestThumbnailURLLargeNonCustom string `json:"latest_thumbnail_url_large_noncustom,omitempty"`
	// LatestThumbnailURLHugeNonCustom is the huge non-custom thumbnail URL for the latest file.
	LatestThumbnailURLHugeNonCustom string `json:"latest_thumbnail_url_huge_noncustom,omitempty"`

	// LatestThumbMediumX is the medium thumbnail width for the latest file.
	LatestThumbMediumX IntString `json:"latest_thumb_medium_x,omitempty"`
	// LatestThumbMediumY is the medium thumbnail height for the latest file.
	LatestThumbMediumY IntString `json:"latest_thumb_medium_y,omitempty"`
	// LatestThumbLargeX is the large thumbnail width for the latest file.
	LatestThumbLargeX IntString `json:"latest_thumb_large_x,omitempty"`
	// LatestThumbLargeY is the large thumbnail height for the latest file.
	LatestThumbLargeY IntString `json:"latest_thumb_large_y,omitempty"`
	// LatestThumbHugeX is the huge thumbnail width for the latest file.
	LatestThumbHugeX IntString `json:"latest_thumb_huge_x,omitempty"`
	// LatestThumbHugeY is the huge thumbnail height for the latest file.
	LatestThumbHugeY IntString `json:"latest_thumb_huge_y,omitempty"`
	// LatestThumbMediumNonCustomX is the medium non-custom thumbnail width for the latest file.
	LatestThumbMediumNonCustomX IntString `json:"latest_thumb_medium_noncustom_x,omitempty"`
	// LatestThumbMediumNonCustomY is the medium non-custom thumbnail height for the latest file.
	LatestThumbMediumNonCustomY IntString `json:"latest_thumb_medium_noncustom_y,omitempty"`
	// LatestThumbLargeNonCustomX is the large non-custom thumbnail width for the latest file.
	LatestThumbLargeNonCustomX IntString `json:"latest_thumb_large_noncustom_x,omitempty"`
	// LatestThumbLargeNonCustomY is the large non-custom thumbnail height for the latest file.
	LatestThumbLargeNonCustomY IntString `json:"latest_thumb_large_noncustom_y,omitempty"`
	// LatestThumbHugeNonCustomX is the huge non-custom thumbnail width for the latest file.
	LatestThumbHugeNonCustomX IntString `json:"latest_thumb_huge_noncustom_x,omitempty"`
	// LatestThumbHugeNonCustomY is the huge non-custom thumbnail height for the latest file.
	LatestThumbHugeNonCustomY IntString `json:"latest_thumb_huge_noncustom_y,omitempty"`
}

// Thumbs contains thumbnail URLs and dimensions for the primary file.
type Thumbs struct {
	// ThumbnailURLMedium is the medium thumbnail URL.
	ThumbnailURLMedium string `json:"thumbnail_url_medium,omitempty"`
	// ThumbnailURLLarge is the large thumbnail URL.
	ThumbnailURLLarge string `json:"thumbnail_url_large,omitempty"`
	// ThumbnailURLHuge is the huge thumbnail URL.
	ThumbnailURLHuge string `json:"thumbnail_url_huge,omitempty"`
	// ThumbnailURLMediumNonCustom is the medium non-custom thumbnail URL.
	ThumbnailURLMediumNonCustom string `json:"thumbnail_url_medium_noncustom,omitempty"`
	// ThumbnailURLLargeNonCustom is the large non-custom thumbnail URL.
	ThumbnailURLLargeNonCustom string `json:"thumbnail_url_large_noncustom,omitempty"`
	// ThumbnailURLHugeNonCustom is the huge non-custom thumbnail URL.
	ThumbnailURLHugeNonCustom string `json:"thumbnail_url_huge_noncustom,omitempty"`

	// ThumbMediumX is the medium thumbnail width.
	ThumbMediumX IntString `json:"thumb_medium_x,omitempty"`
	// ThumbMediumY is the medium thumbnail height.
	ThumbMediumY IntString `json:"thumb_medium_y,omitempty"`
	// ThumbLargeX is the large thumbnail width.
	ThumbLargeX IntString `json:"thumb_large_x,omitempty"`
	// ThumbLargeY is the large thumbnail height.
	ThumbLargeY IntString `json:"thumb_large_y,omitempty"`
	// ThumbHugeX is the huge thumbnail width.
	ThumbHugeX IntString `json:"thumb_huge_x,omitempty"`
	// ThumbHugeY is the huge thumbnail height.
	ThumbHugeY IntString `json:"thumb_huge_y,omitempty"`
	// ThumbMediumNonCustomX is the medium non-custom thumbnail width.
	ThumbMediumNonCustomX IntString `json:"thumb_medium_noncustom_x,omitempty"`
	// ThumbMediumNonCustomY is the medium non-custom thumbnail height.
	ThumbMediumNonCustomY IntString `json:"thumb_medium_noncustom_y,omitempty"`
	// ThumbLargeNonCustomX is the large non-custom thumbnail width.
	ThumbLargeNonCustomX IntString `json:"thumb_large_noncustom_x,omitempty"`
	// ThumbLargeNonCustomY is the large non-custom thumbnail height.
	ThumbLargeNonCustomY IntString `json:"thumb_large_noncustom_y,omitempty"`
	// ThumbHugeNonCustomX is the huge non-custom thumbnail width.
	ThumbHugeNonCustomX IntString `json:"thumb_huge_noncustom_x,omitempty"`
	// ThumbHugeNonCustomY is the huge non-custom thumbnail height.
	ThumbHugeNonCustomY IntString `json:"thumb_huge_noncustom_y,omitempty"`
}

// SubmissionDetailsResponse is the top-level result returned by
// Client.SubmissionDetails.
type SubmissionDetailsResponse struct {
	// SID is the current session ID.
	SID string `json:"sid"`
	// ResultsCount is the number of submissions returned.
	ResultsCount IntString `json:"results_count"`
	// UserLocation identifies the timezone context used for user-time timestamps.
	UserLocation string `json:"user_location"`
	// Submissions is the list of returned submissions.
	Submissions []SubmissionDetails `json:"submissions"`
}

// SubmissionFavoritesResponse is returned by User.SubmissionFavorites.
type SubmissionFavoritesResponse struct {
	// Sid is the current session ID.
	Sid string `json:"sid"`
	// Users is the list of users who favorited the submission.
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
