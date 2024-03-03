package api

import (
	"encoding/json"
	"fmt"
	"github.com/ellypaws/inkbunny/utils"
	"io"
)

type SearchRequest struct {
	SID                string `json:"sid"`
	OutputMode         string `json:"output_mode,omitempty"`
	RID                string `json:"rid,omitempty"`
	SubmissionIDsOnly  string `json:"submission_ids_only,omitempty"`
	SubmissionsPerPage int    `json:"submissions_per_page,omitempty"`
	Page               int    `json:"page,omitempty"`
	KeywordsList       string `json:"keywords_list,omitempty"`
	NoSubmissions      string `json:"no_submissions,omitempty"`
	GetRID             string `json:"get_rid,omitempty"`
	FieldJoinType      string `json:"field_join_type,omitempty"` // "or" or "and"
	Text               string `json:"text,omitempty"`
	StringJoinType     string `json:"string_join_type,omitempty"`
	Keywords           string `json:"keywords,omitempty"`
	Title              string `json:"title,omitempty"`
	Description        string `json:"description,omitempty"`
	MD5                string `json:"md5,omitempty"`
	KeywordID          string `json:"keyword_id,omitempty"`
	Username           string `json:"username,omitempty"`
	UserID             string `json:"user_id,omitempty"`
	FavsUserID         string `json:"favs_user_id,omitempty"`
	UnreadSubmissions  string `json:"unread_submissions,omitempty"`
	Type               string `json:"type,omitempty"`
	Sales              string `json:"sales,omitempty"`
	PoolID             string `json:"pool_id,omitempty"`
	OrderBy            string `json:"orderby,omitempty"` // Values: create_datetime, unread_datetime, views, total_print_sales, total_digital_sales, total_sales, username, fav_datetime, fav_stars, pool_order. Default: create_datetime.
	DaysLimit          int    `json:"dayslimit,omitempty"`
	Random             string `json:"random,omitempty"`
	Scraps             string `json:"scraps,omitempty"`
	CountLimit         int    `json:"count_limit,omitempty"`
}

type SearchResponse struct {
	Sid                  string            `json:"sid"`
	UserLocation         string            `json:"user_location"`
	ResultsCountAll      int               `json:"results_count_all"`
	ResultsCountThisPage int               `json:"results_count_thispage"`
	PagesCount           int               `json:"pages_count"`
	Page                 int               `json:"page"`
	RID                  string            `json:"rid,omitempty"`
	RIDTTL               string            `json:"rid_ttl,omitempty"`
	SearchParams         map[string]string `json:"search_params"`
	KeywordList          []struct {
		KeywordID        string `json:"keyword_id"`
		KeywordName      string `json:"keyword_name"`
		SubmissionsCount int    `json:"submissions_count"`
	} `json:"keyword_list,omitempty"`
	Submissions []struct {
		SubmissionBasic
		UnreadDateSystem string    `json:"unread_datetime_system"`
		UnreadDateUser   string    `json:"unread_datetime"`
		Updated          BooleanYN `json:"updated"`
		Stars            string    `json:"stars"`
	} `json:"submissions,omitempty"`
}

func (user Credentials) SearchSubmissions(req SearchRequest) (SearchResponse, error) {
	if !user.LoggedIn() {
		return SearchResponse{}, ErrNotLoggedIn
	}
	if req.SID == "" {
		req.SID = user.Sid
	}

	resp, err := user.Get(apiURL("search", utils.StructToUrlValues(req)))

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
	return user.SearchSubmissions(SearchRequest{SID: user.Sid, Username: user.Username})
}

func (user Credentials) UserSubmissions(username string) (SearchResponse, error) {
	return user.SearchSubmissions(SearchRequest{SID: user.Sid, Username: username})
}
