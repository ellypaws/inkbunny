package inkbunny

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SubmissionSearchRequest configures Client.SearchSubmissions.
//
// Most filter fields are simple strings or IntString values because Inkbunny is
// permissive about string-encoded numbers. Boolean toggles use BooleanYN or
// *BooleanYN so callers can distinguish between sending No and leaving Inkbunny's
// default unchanged. Common documented defaults include JSON, 30, 1, No, No,
// No, FieldJoinTypeOr, JoinTypeAnd, Yes, No, No, No,
// OrderByCreateDatetime, No, ScrapsBoth, and 50000.
type SubmissionSearchRequest struct {
	// SID overrides the client's current session for this call.
	SID string `json:"sid" query:"sid"`
	// OutputMode selects the response format. Use JSON or XML.
	// Defaults to JSON.
	OutputMode OutputMode `json:"output_mode,omitempty" query:"output_mode"`
	// RID pages through an existing result set instead of running a fresh search.
	// When set, Inkbunny ignores the mode-1 search filters below.
	RID string `json:"rid,omitempty" query:"rid"`
	// SubmissionIDsOnly requests only submission IDs instead of full objects.
	// Defaults to No.
	SubmissionIDsOnly BooleanYN `json:"submission_ids_only,omitempty" query:"submission_ids_only"`
	// SubmissionsPerPage limits a page to at most 100 results. Defaults to 30.
	SubmissionsPerPage IntString `json:"submissions_per_page,omitempty" query:"submissions_per_page"`
	// Page selects the page number. Defaults to 1.
	Page IntString `json:"page,omitempty" query:"page"`
	// KeywordsList asks Inkbunny to populate KeywordList in the response. Use Yes or No.
	// Defaults to No.
	KeywordsList BooleanYN `json:"keywords_list,omitempty" query:"keywords_list"`
	// NoSubmissions suppresses full submission payloads when only counts or keywords are needed.
	// Defaults to No.
	NoSubmissions BooleanYN `json:"no_submissions,omitempty" query:"no_submissions"`
	// GetRID asks Inkbunny to return a reusable RID for later paging.
	// Defaults to No.
	GetRID BooleanYN `json:"get_rid,omitempty" query:"get_rid"`

	// FieldJoinType combines keyword, title, and description-style field matches.
	// Use FieldJoinTypeOr or FieldJoinTypeAnd. Defaults to FieldJoinTypeOr.
	FieldJoinType FieldJoinType `json:"field_join_type,omitempty" query:"field_join_type"`
	// Text is the free-text query used by the selected search fields.
	Text string `json:"text,omitempty" query:"text"`
	// StringJoinType controls how words inside Text are matched.
	// Use JoinTypeAnd, JoinTypeOr, or JoinTypeExact. Defaults to JoinTypeAnd.
	StringJoinType JoinType `json:"string_join_type,omitempty" query:"string_join_type"`
	// SearchInKeywords toggles keyword matching for Text.
	// Nil keeps the default, which is Yes.
	SearchInKeywords *BooleanYN `json:"keywords,omitempty" query:"keywords"`
	// Title toggles title matching for Text.
	// Nil keeps the default, which is No.
	Title *BooleanYN `json:"title,omitempty" query:"title"`
	// Description toggles description and writing matching for Text.
	// Nil keeps the default, which is No.
	Description *BooleanYN `json:"description,omitempty" query:"description"`
	// MD5 switches Text into MD5-search mode. Nil keeps the default, which is No.
	MD5 *BooleanYN `json:"md5,omitempty" query:"md5"`
	// KeywordID bypasses text search and filters directly by keyword ID.
	KeywordID IntString `json:"keyword_id,omitempty" query:"keyword_id"`
	// Username restricts results to one owner username.
	Username string `json:"username,omitempty" query:"username"`
	// UserID restricts results to one owner user ID.
	UserID IntString `json:"user_id,omitempty" query:"user_id"`
	// FavsUserID restricts results to a user's favorites.
	FavsUserID IntString `json:"favs_user_id,omitempty" query:"favs_user_id"`
	// UnreadSubmissions restricts results to the current session's unread feed.
	UnreadSubmissions BooleanYN `json:"unread_submissions,omitempty" query:"unread_submissions"`
	// Type restricts results to one or more SubmissionType values.
	// Use SubmissionType* constants.
	Type SubmissionTypes `json:"type,omitempty" query:"type"`
	// Sales filters on legacy sales state. Use SalesFilter* constants when needed.
	// Deprecated: Sales are no longer part of Inkbunny.
	Sales SalesFilter `json:"sales,omitempty" query:"sales"`
	// PoolID restricts results to one pool.
	PoolID IntString `json:"pool_id,omitempty" query:"pool_id"`
	// OrderBy controls sorting. Use OrderBy* constants.
	// Defaults to OrderByCreateDatetime.
	OrderBy OrderBy `json:"orderby,omitempty" query:"orderby"`
	// DaysLimit restricts results to submissions from the last N days.
	DaysLimit IntString `json:"dayslimit,omitempty" query:"dayslimit"`
	// Random shuffles the final result set after filtering and ordering. Use Yes or No.
	// Defaults to No.
	Random BooleanYN `json:"random,omitempty" query:"random"`
	// Scraps controls scraps filtering. Use ScrapsBoth, ScrapsNo, or ScrapsOnly.
	// Defaults to ScrapsBoth.
	Scraps Scraps `json:"scraps,omitempty" query:"scraps"`
	// CountLimit caps the total number of matches considered by Inkbunny.
	// Defaults to 50000.
	CountLimit IntString `json:"count_limit,omitempty" query:"count_limit"`
}

type Scraps = string

const (
	ScrapsBoth Scraps = "both"
	ScrapsNo   Scraps = "no"
	ScrapsOnly Scraps = "only"
)

type SubmissionSearchResponse struct {
	// SID is the current session ID.
	SID string `json:"sid"`
	// UserLocation identifies the location used for user-time timestamps.
	UserLocation string `json:"user_location"`
	// ResultsCountAll is the total number of results across all pages.
	ResultsCountAll IntString `json:"results_count_all"`
	// ResultsCountThisPage is the number of results on the current page.
	ResultsCountThisPage IntString `json:"results_count_thispage"`
	// PagesCount is the total number of result pages.
	PagesCount IntString `json:"pages_count"`
	// Page is the current page number.
	Page IntString `json:"page"`
	// RID is the reusable result-set ID when GetRID was enabled.
	RID string `json:"rid,omitempty"`
	// RIDTTL is the server's human-readable lifetime string for RID.
	RIDTTL string `json:"rid_ttl,omitempty"`
	// RIDTTLDuration is the package-parsed duration derived from RIDTTL.
	RIDTTLDuration time.Duration `json:"-"`
	// RIDExpiry is the package-computed expiry time when RIDTTLDuration is available.
	RIDExpiry time.Time `json:"-"`
	// SearchParams echoes the search parameters used to build the result set.
	SearchParams []SearchParam `json:"search_params"`
	// KeywordList contains top keywords for the current result page when requested.
	KeywordList []KeywordList `json:"keyword_list,omitempty"`
	// Submissions contains the returned search results unless suppressed.
	Submissions []SubmissionSearch `json:"submissions,omitempty"`
	client      *Client
}

// KeywordList describes one keyword aggregate from a search response.
type KeywordList struct {
	// KeywordID is the keyword ID.
	KeywordID IntString `json:"keyword_id"`
	// KeywordName is the keyword text.
	KeywordName string `json:"keyword_name"`
	// SubmissionsCount is the systemwide count of submissions tagged with the keyword.
	SubmissionsCount IntString `json:"submissions_count"`
}

// SubmissionSearch is the per-result submission model returned by search.
type SubmissionSearch struct {
	// SubmissionBasic contains the fields shared with submission-details responses.
	SubmissionBasic
	// UnreadDateSystem is when the submission entered the unread list in system time.
	UnreadDateSystem string `json:"unread_datetime,omitempty"`
	// UnreadDateUser is when the submission entered the unread list in the user's local time.
	UnreadDateUser string `json:"unread_datetime_usertime,omitempty"`
	// Updated reports whether the unread submission was updated since it was added.
	Updated BooleanYN `json:"updated,omitempty"`
	// Stars is the favorite-star count when the search mode includes favorites.
	Stars IntString `json:"stars,omitempty"`
}

// SearchParam is the search parameters that were used to find these search results.
type SearchParam struct {
	Name  string `json:"param_name"`
	Value string `json:"param_value"`
	// Type is kept as a compatibility alias for older code that read the search
	// parameter value from this field before the Inkbunny tag mapping was corrected.
	Type string `json:"-"`
}

func (s *SearchParam) UnmarshalJSON(data []byte) error {
	type rawSearchParam struct {
		Name  string `json:"param_name"`
		Value string `json:"param_value"`
	}

	var raw rawSearchParam
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	s.Name = raw.Name
	s.Value = raw.Value
	s.Type = raw.Value
	return nil
}

func (s SearchParam) MarshalJSON() ([]byte, error) {
	value := s.Value
	if value == "" {
		value = s.Type
	}

	type rawSearchParam struct {
		Name  string `json:"param_name"`
		Value string `json:"param_value"`
	}

	return json.Marshal(rawSearchParam{Name: s.Name, Value: value})
}

type SubmissionType int

const (
	SubmissionTypeAny                       SubmissionType = iota
	SubmissionTypePicturePinup                             // 1 - Picture/Pinup
	SubmissionTypeSketch                                   // 2 - Sketch
	SubmissionTypePictureSeries                            // 3 - Picture Series
	SubmissionTypeComic                                    // 4 - Comic
	SubmissionTypePortfolio                                // 5 - Portfolio
	SubmissionTypeShockwaveFlashAnimation                  // 6 - Shockwave/Flash - Animation
	SubmissionTypeShockwaveFlashInteractive                // 7 - Shockwave/Flash - Interactive
	SubmissionTypeVideoFeatureLength                       // 8 - Video - Feature Length
	SubmissionTypeVideoAnimation3DCGI                      // 9 - Video - Animation/3D/CGI
	SubmissionTypeMusicSingleTrack                         // 10 - Music - Single Track
	SubmissionTypeMusicAlbum                               // 11 - Music - Album
	SubmissionTypeWritingDocument                          // 12 - Writing - Document
	SubmissionTypeCharacterSheet                           // 13 - Character Sheet
	SubmissionTypePhotography                              // 14 - Photography - Fursuit/Sculpture/Jewelry/etc
)

type SubmissionTypes []SubmissionType

func (s SubmissionTypes) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteRune('"')
	for i, t := range s {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(strconv.Itoa(int(t)))
	}
	buffer.WriteRune('"')
	return buffer.Bytes(), nil
}

func (s *SubmissionTypes) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if bytes.Equal(data, []byte(`null`)) {
		return nil
	}
	var submissionTypes []SubmissionType
	if bytes.HasPrefix(data, []byte(`[`)) && bytes.HasSuffix(data, []byte(`]`)) {
		split := strings.Split(strings.Trim(string(data), `["]`), ",")
		submissionTypes = make([]SubmissionType, len(split))
		for i, t := range split {
			if len(t) == 0 {
				continue
			}
			atoi, err := strconv.Atoi(t)
			if err != nil {
				return fmt.Errorf("failed to parse submission type: %w", err)
			}
			submissionTypes[i] = SubmissionType(atoi)
		}
		*s = submissionTypes
	}
	return nil
}

// WithClient sets the client used by AllPages, AllSubmissions, AllDetails, and Details.
// This is useful when a SubmissionSearchResponse has been deserialized or passed
// across a package boundary and the unexported client field is nil.
func (s *SubmissionSearchResponse) WithClient(c *Client) *SubmissionSearchResponse {
	s.client = c
	return s
}

// AllPages returns a sequence of all the pages in a submission search response, repeatedly calling Client.SearchSubmissions.
// Make sure you set SubmissionSearchRequest.GetRID to types.Yes prior or the other pages might not have the correct results.
// Additionally, one should also check SubmissionSearchResponse.RIDTTLDuration or SubmissionSearchResponse.RIDExpiry.
func (s SubmissionSearchResponse) AllPages() iter.Seq2[SubmissionSearchResponse, error] {
	return func(yield func(SubmissionSearchResponse, error) bool) {
		for i := range s.PagesCount.Iter() {
			if i == 0 {
				if !yield(s, nil) {
					return
				}
				continue
			}
			request := SubmissionSearchRequest{
				SID:  s.SID,
				RID:  s.RID,
				Page: i + 1,
			}
			if !yield(s.client.Get().SearchSubmissions(request)) {
				return
			}
		}
	}
}

// AllSubmissions returns a sequence of all submission lists across all pages of the search results, repeatedly calling Client.SearchSubmissions.
// Make sure you set SubmissionSearchRequest.GetRID to types.Yes prior or the other pages might not have the correct results.
// Additionally, one should also check SubmissionSearchResponse.RIDTTLDuration or SubmissionSearchResponse.RIDExpiry.
func (s SubmissionSearchResponse) AllSubmissions() iter.Seq2[[]SubmissionSearch, error] {
	return func(yield func([]SubmissionSearch, error) bool) {
		for i := range s.PagesCount.Iter() {
			if i == 0 {
				if !yield(s.Submissions, nil) {
					return
				}
				continue
			}
			request := SubmissionSearchRequest{
				SID:  s.SID,
				RID:  s.RID,
				Page: i + 1,
			}
			response, err := s.client.Get().SearchSubmissions(request)
			if !yield(response.Submissions, err) {
				return
			}
		}
	}
}

// Details returns the SubmissionDetails of the current page
func (s SubmissionSearchResponse) Details() (SubmissionDetailsResponse, error) {
	ids := make([]string, len(s.Submissions))
	for i, v := range s.Submissions {
		ids[i] = v.SubmissionID.String()
	}
	return s.client.Get().SubmissionDetails(SubmissionDetailsRequest{
		SID:               s.SID,
		SubmissionIDSlice: ids,
	})
}

// AllDetails returns an iterator that yields a SubmissionDetailsResponse for each page of
// search results. The provided SubmissionDetailsRequest is used as a template; SID and
// SubmissionIDSlice are overwritten per page. Fields like ShowDescription, ShowPools, etc.
// are forwarded as-is.
func (s SubmissionSearchResponse) AllDetails(req SubmissionDetailsRequest) iter.Seq2[SubmissionDetailsResponse, error] {
	return func(yield func(SubmissionDetailsResponse, error) bool) {
		for page, err := range s.AllPages() {
			if err != nil {
				yield(SubmissionDetailsResponse{}, err)
				return
			}
			ids := make([]string, len(page.Submissions))
			for i, sub := range page.Submissions {
				ids[i] = sub.SubmissionID.String()
			}
			req.SID = page.SID
			req.SubmissionIDSlice = ids
			if !yield(s.client.Get().SubmissionDetails(req)) {
				return
			}
		}
	}
}

func (u *User) SearchSubmissions(req SubmissionSearchRequest) (SubmissionSearchResponse, error) {
	return u.SearchSubmissionsContext(context.Background(), req)
}

// SearchSubmissionsContext is like SearchSubmissions but accepts a context.Context
// for per-call cancellation and timeout control.
func (u *User) SearchSubmissionsContext(ctx context.Context, req SubmissionSearchRequest) (SubmissionSearchResponse, error) {
	if req.SID == "" {
		if u.SID == "" {
			return SubmissionSearchResponse{}, ErrNotLoggedIn
		}
		req.SID = u.SID
	}

	return u.Client().SearchSubmissionsContext(ctx, req)
}

func (c *Client) SearchSubmissions(req SubmissionSearchRequest) (SubmissionSearchResponse, error) {
	return c.SearchSubmissionsContext(c.ctx, req)
}

// SearchSubmissionsContext is like SearchSubmissions but accepts a context.Context
// for per-call cancellation and timeout control.
func (c *Client) SearchSubmissionsContext(ctx context.Context, req SubmissionSearchRequest) (SubmissionSearchResponse, error) {
	if req.SID == "" {
		return SubmissionSearchResponse{}, ErrEmptySID
	}
	response, err := PostDecode[SubmissionSearchResponse](c.withContext(ctx), ApiUrl("search"), req)
	if err != nil {
		return response, err
	}

	response.client = c

	if response.RIDTTL != "" {
		response.RIDTTLDuration = TTLToDuration(response.RIDTTL)
		response.RIDExpiry = time.Now().Add(response.RIDTTLDuration)
	}

	return response, err
}

func SearchSubmissions(req SubmissionSearchRequest) (SubmissionSearchResponse, error) {
	return DefaultClient.SearchSubmissions(req)
}

var shortDuration = regexp.MustCompile(`\d+[smhdwy]`)

func TTLToDuration(ttl string) time.Duration {
	var d time.Duration
	matches := shortDuration.FindAllString(strings.ReplaceAll(ttl, " ", ""), -1)
	for _, match := range matches {
		i, err := strconv.Atoi(match[:len(match)-1])
		if err != nil {
			continue
		}
		switch match[len(match)-1] {
		case 's':
			d += time.Second * time.Duration(i)
		case 'm':
			d += time.Minute * time.Duration(i)
		case 'h':
			d += time.Hour * time.Duration(i)
		case 'd':
			d += time.Hour * 24 * time.Duration(i)
		case 'w':
			d += time.Hour * 24 * 7 * time.Duration(i)
		case 'y':
			d += time.Hour * 24 * 365 * time.Duration(i)
		}
	}
	return d
}
