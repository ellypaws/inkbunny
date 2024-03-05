package api

import (
	"encoding/json"
	"fmt"
	"github.com/ellypaws/inkbunny/api/utils"
	"io"
)

type SubmissionSearchRequest struct {
	SID                string     `json:"sid"`
	OutputMode         OutputMode `json:"output_mode,omitempty"`
	RID                string     `json:"rid,omitempty"`
	SubmissionIDsOnly  string     `json:"submission_ids_only,omitempty"`
	SubmissionsPerPage IntString  `json:"submissions_per_page,omitempty"`
	Page               IntString  `json:"page,omitempty"`
	// Not to be confused with Text. This is a boolean value to return list of Top 100 Keywords.
	// Return list of Top 100 Keywords associated with all submissions on current results page.
	// Note that this list includes both officially assigned keywords and also keywords
	// suggested for this submission by other users.
	KeywordsList  BooleanYN `json:"keywords_list,omitempty"`
	NoSubmissions BooleanYN `json:"no_submissions,omitempty"`
	GetRID        BooleanYN `json:"get_rid,omitempty"`
	FieldJoinType JoinType  `json:"field_join_type,omitempty"` // "or" or "and"
	// Text to search chosen fields for. eg "dragon", "wolf", etc.
	// A Full Text search is performed using this string (see the meaning of Full Text searches in the Postgresql Documentation).
	// The characters "_" and "," are converted to spaces automatically.
	// Characters which have special meanings for Full Text searches in Postgresql (such as |, &, :, ! and ~) are ignored.
	//
	// Note: At least one of the Search Field parameters "keywords", "title", "description" or "md5" must be set to “yes” for text search to work.
	// By default, "keywords" is set to "yes", so all searches with no Search Field specified will search in keywords.
	// Values: (Any text string).
	//
	// Default: n/a. Required: No
	Text              string         `json:"text,omitempty"`
	StringJoinType    JoinType       `json:"string_join_type,omitempty"`
	Keywords          BooleanYN      `json:"keywords,omitempty"`
	Title             BooleanYN      `json:"title,omitempty"`
	Description       BooleanYN      `json:"description,omitempty"`
	MD5               BooleanYN      `json:"md5,omitempty"`
	KeywordID         string         `json:"keyword_id,omitempty"`
	Username          string         `json:"username,omitempty"`
	UserID            string         `json:"user_id,omitempty"`
	FavsUserID        string         `json:"favs_user_id,omitempty"`
	UnreadSubmissions BooleanYN      `json:"unread_submissions,omitempty"`
	Type              SubmissionType `json:"type,omitempty"`
	Sales             string         `json:"sales,omitempty"` // Values: forsale, digital, prints
	PoolID            string         `json:"pool_id,omitempty"`
	OrderBy           string         `json:"orderby,omitempty"` // Values: create_datetime, unread_datetime, views, total_print_sales, total_digital_sales, total_sales, username, fav_datetime, fav_stars, pool_order. Default: create_datetime.
	DaysLimit         IntString      `json:"dayslimit,omitempty"`
	Random            BooleanYN      `json:"random,omitempty"`
	// Scraps Set how submissions marked as “Scraps” are returned.
	// Possible values are:
	// 	both – show submissions from Scraps and Main galleries.
	// 	no – Do not show Scraps. Shows only submissions from Main galleries.
	// 	only – Show only submissions from Scraps galleries, not Main galleries.
	Scraps     string    `json:"scraps,omitempty"`
	CountLimit IntString `json:"count_limit,omitempty"`
}

type SearchResponse struct {
	Sid                  string    `json:"sid"`
	UserLocation         string    `json:"user_location"`
	ResultsCountAll      IntString `json:"results_count_all"`
	ResultsCountThisPage IntString `json:"results_count_thispage"`
	PagesCount           IntString `json:"pages_count"`
	Page                 IntString `json:"page"`
	RID                  string    `json:"rid,omitempty"`
	RIDTTL               string    `json:"rid_ttl,omitempty"`
	SearchParams         any       `json:"search_params"`
	KeywordList          []struct {
		KeywordID        string    `json:"keyword_id"`
		KeywordName      string    `json:"keyword_name"`
		SubmissionsCount IntString `json:"submissions_count"`
	} `json:"keyword_list,omitempty"`
	Submissions []struct {
		SubmissionBasic
		UnreadDateSystem string    `json:"unread_datetime_system"`
		UnreadDateUser   string    `json:"unread_datetime"`
		Updated          BooleanYN `json:"updated"`
		Stars            string    `json:"stars"`
	} `json:"submissions,omitempty"`
}

func (user Credentials) SearchSubmissions(req SubmissionSearchRequest) (SearchResponse, error) {
	if !user.LoggedIn() {
		return SearchResponse{}, ErrNotLoggedIn
	}
	if req.SID == "" {
		req.SID = user.Sid
	}

	resp, err := user.Get(ApiUrl("search", utils.StructToUrlValues(req)))

	if err != nil {
		return SearchResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SearchResponse{}, err
	}

	if err := CheckError(body); err != nil {
		return SearchResponse{}, fmt.Errorf("error searching submissions: %w", err)
	}

	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return SearchResponse{}, err
	}

	return searchResp, nil
}

func (user Credentials) OwnSubmissions() (SearchResponse, error) {
	return user.SearchSubmissions(SubmissionSearchRequest{SID: user.Sid, Username: user.Username})
}

func (user Credentials) UserSubmissions(username string) (SearchResponse, error) {
	return user.SearchSubmissions(SubmissionSearchRequest{SID: user.Sid, Username: username})
}
