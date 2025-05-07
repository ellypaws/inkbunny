package inkbunny

import (
	"io"

	"github.com/ellypaws/inkbunny/types"
)

type SubmissionEditRequest struct {
	SID          string          `json:"sid"`
	SubmissionID types.IntString `json:"submission_id"`
	Title        *string         `json:"title,omitempty"`
	Description  *string         `json:"desc,omitempty"`
	Story        io.Reader       `json:"story,omitempty"`
	// Should html entities (eg: &nbsp; &gt; &#1234;) in uploaded text (title, desc,
	// story) be converted to normal characters before being saved? Boolean. Note:
	// By default, html entities will be treated as plain text and will not be
	// converted back to regular characters for display on the Inkbunny website. Eg:
	// If you upload the text &nbsp; as part of a title, it will display literally
	// as &nbsp; on the web page. If your uploaded text is likely to contain html
	// entities then always set this option to types.Yes
	ConvertHTMLEntities types.BooleanYN  `json:"convert_html_entities,omitempty"`
	SubmissionType      SubmissionType   `json:"submission_type,omitempty"`
	Scraps              *types.BooleanYN `json:"scraps,omitempty"`
	// Announce this submission via the owner's Twitter account, if configured and
	// enabled. If this property is set to "yes", then the announcement occurs the
	// first time the submission is set Public. Note that announcement via Twitter
	// only occurs if the owner's Twitter account is authenticated via their account
	// settings at https://inkbunny.net/account.php, they have "Tweet Submissions"
	// turned on in their account settings, and "use_twitter" is enabled for this
	// submission. Boolean.
	UseTwitter *types.BooleanYN `json:"use_twitter,omitempty"`
	// Lets you choose if you want to send an image along in the tweet announcing the submission. Available Options are:
	// 0 - Send only Text
	// 1 - Send the Thumbnail (if a custom thumbnail was added, it will send that one. If not, the generated one)
	// 2 - Send the Full Picture (it will send a proportional, 920px-wide version of the full picture)
	// Values: (Twitter Image Preference ID). Default: set on upload to current user preference, which defaults to 1. Required: No
	TwitterImagePref *int `json:"twitter_image_pref,omitempty"`
	// Change Public/Non-Public status of the submission
	Public *types.BooleanYN `json:"visibility,omitempty"`
	// Will only announce the first time a submission is set to Public, and if Public is set to true.
	Announce *types.BooleanYN `json:"-"`
	// Keywords for this submission. Keyword entries must be separated by commas or
	// spaces. When adding new keywords, all the old keywords must be specified here
	// too. The entry here entirely REPLACES the existing keywords list for this
	// submission. Sending the keywords variable but leaving its value blank will
	// REMOVE all keywords from this submission. To avoid clearing the keywords when
	// updating a submission, simply do not send the keywords property and the
	// existing keywords will be preserved.
	Keywords []string `json:"-"`

	Nudity         *types.BooleanYN `json:"tag[2],omitempty"`
	MildViolence   *types.BooleanYN `json:"tag[3],omitempty"`
	Sexual         *types.BooleanYN `json:"tag[4],omitempty"`
	StrongViolence *types.BooleanYN `json:"tag[5],omitempty"`

	GuestBlock  *types.BooleanYN `json:"guest_block,omitempty"`
	FriendsOnly *types.BooleanYN `json:"friends_only,omitempty"`
}
