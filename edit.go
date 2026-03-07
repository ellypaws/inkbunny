package inkbunny

import (
	"context"
	"io"
)

// SubmissionEditRequest configures Client.EditSubmission.
//
// Optional string and BooleanYN fields use pointers so callers can distinguish
// between "leave unchanged" and "set an explicit empty or no value". The helper
// new is convenient for building these pointers in calling code.
type SubmissionEditRequest struct {
	// SID overrides the current user session for this call.
	SID string `json:"sid"`
	// SubmissionID identifies the submission to edit.
	SubmissionID IntString `json:"submission_id"`
	// Title replaces the submission title when set.
	Title *string `json:"title,omitempty"`
	// Description replaces the submission description when set.
	Description *string `json:"desc,omitempty"`
	// Story replaces the submission writing or story text when set.
	Story io.Reader `json:"story,omitempty"`
	// ConvertHTMLEntities enables HTML-entity decoding on text fields.
	// Should html entities (eg: &nbsp; &gt; &#1234;) in uploaded text (title, desc, story) be converted to normal characters before being saved? Boolean.
	// Note: By default, html entities will be treated as plain text and will not be converted back to regular characters for display on the Inkbunny website.
	// Eg: If you upload the text &nbsp; as part of a title, it will display literally as &nbsp; on the web page.
	// If your uploaded text is likely to contain html entities then always set this option to Yes.
	// Use Yes or No.
	ConvertHTMLEntities BooleanYN `json:"convert_html_entities,omitempty"`
	// SubmissionType replaces the submission type when non-zero.
	// Use SubmissionType* constants.
	SubmissionType SubmissionType `json:"type,omitempty"`
	// Scraps moves the submission into or out of scraps when set. Use Yes or No.
	Scraps *BooleanYN `json:"scraps,omitempty"`
	// UseTwitter enables tweet-on-publish behavior when the account supports it.
	// Use Yes, No, or nil to leave the current setting unchanged.
	// Inkbunny has no default here; omitted means no update is made.
	UseTwitter *BooleanYN `json:"use_twitter,omitempty"`
	// TwitterImagePref controls the tweet image mode.
	// Inkbunny currently expects 0 for text only, 1 for thumbnail, and 2 for full image.
	// When omitted, Inkbunny leaves the current value unchanged.
	TwitterImagePref *int `json:"twitter_image_pref,omitempty"`
	// Public maps to the visibility field using BooleanYN semantics.
	// Inkbunny has no default here; omitted means no update is made.
	// When set to Yes, the package defaults Notify to Yes unless overridden.
	Public *BooleanYN `json:"visibility,omitempty"`
	// Notify overrides the package's default watcher notification behavior for Public.
	// This package defaults Notify to Yes when Public is Yes. The raw notify parameter
	// defaults to No.
	Notify *BooleanYN `json:"-"`
	// Keywords replaces the entire keyword list when provided.
	// Leave it nil to preserve current keywords, or pass an empty slice to clear them.
	Keywords []string `json:"-"`

	// Nudity toggles rating tag 2 using BooleanYN semantics.
	Nudity *BooleanYN `json:"tag[2],omitempty"`
	// MildViolence toggles rating tag 3 using BooleanYN semantics.
	MildViolence *BooleanYN `json:"tag[3],omitempty"`
	// Sexual toggles rating tag 4 using BooleanYN semantics.
	Sexual *BooleanYN `json:"tag[4],omitempty"`
	// StrongViolence toggles rating tag 5 using BooleanYN semantics.
	StrongViolence *BooleanYN `json:"tag[5],omitempty"`

	// GuestBlock toggles guest access using BooleanYN semantics.
	GuestBlock *BooleanYN `json:"guest_block,omitempty"`
	// FriendsOnly limits visibility to friends using BooleanYN semantics.
	FriendsOnly *BooleanYN `json:"friends_only,omitempty"`
}

// EditSubmissionResponse is returned by Client.EditSubmission.
type EditSubmissionResponse struct {
	SubmissionID IntString `json:"submission_id"` // Submission ID of the submission that was edited.
	// TwitterAuthSuccess reports whether Twitter authentication succeeded when tweeting was requested.
	TwitterAuthSuccess BooleanYN `json:"twitter_authentication_success"`
}

// EditSubmission edits an existing submission on Inkbunny based on the provided request parameters.
// This method requires a valid session ID (SID) and submission ID.
func (u *User) EditSubmission(req SubmissionEditRequest) (EditSubmissionResponse, error) {
	return u.EditSubmissionContext(context.Background(), req)
}

// EditSubmissionContext is like EditSubmission but accepts a context.Context
// for per-call cancellation and timeout control.
func (u *User) EditSubmissionContext(ctx context.Context, req SubmissionEditRequest) (EditSubmissionResponse, error) {
	if req.SID == "" {
		req.SID = u.SID
	}

	return u.Client().EditSubmissionContext(ctx, req)
}

// EditSubmission edits an existing submission on Inkbunny based on the provided request parameters.
// This method requires a valid session ID (SID) and submission ID.
func (c *Client) EditSubmission(req SubmissionEditRequest) (EditSubmissionResponse, error) {
	return c.EditSubmissionContext(c.ctx, req)
}

// EditSubmissionContext is like EditSubmission but accepts a context.Context
// for per-call cancellation and timeout control.
func (c *Client) EditSubmissionContext(ctx context.Context, req SubmissionEditRequest) (EditSubmissionResponse, error) {
	if req.SID == "" {
		return EditSubmissionResponse{}, ErrEmptySID
	}
	if req.SubmissionID == 0 {
		return EditSubmissionResponse{}, ErrEmptySubID
	}

	values := structToUrlValues(req)
	if req.Notify != nil && !req.Notify.Bool() && req.Public != nil && req.Public.Bool() {
		values.Set("visibility", "yes_nowatch")
	}

	return PostDecode[EditSubmissionResponse](c.withContext(ctx), ApiUrl("editsubmission"), values)
}

// EditSubmission edits an existing submission on Inkbunny based on the provided request parameters.
// This method requires a valid session ID (SID) and submission ID.
func EditSubmission(req SubmissionEditRequest) (EditSubmissionResponse, error) {
	return DefaultClient.EditSubmission(req)
}
