package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ellypaws/inkbunny/api/utils"
	"io"
	"strconv"
	"strings"
)

type SubmissionSearchRequest struct {
	SID                string     `json:"sid" query:"sid"`
	OutputMode         OutputMode `json:"output_mode,omitempty" query:"output_mode"`
	RID                string     `json:"rid,omitempty" query:"rid"`
	SubmissionIDsOnly  BooleanYN  `json:"submission_ids_only,omitempty" query:"submission_ids_only"`
	SubmissionsPerPage IntString  `json:"submissions_per_page,omitempty" query:"submissions_per_page"`
	Page               IntString  `json:"page,omitempty" query:"page"` // Results page number to return. Default: 1.
	// Not to be confused with Text. This is a boolean value to return list of Top 100 Keywords.
	// Return list of Top 100 Keywords associated with all submissions on current results page.
	// Note that this list includes both officially assigned keywords and also keywords
	// suggested for this submission by other users.
	KeywordsList  BooleanYN `json:"keywords_list,omitempty" query:"keywords_list"`
	NoSubmissions BooleanYN `json:"no_submissions,omitempty" query:"no_submissions"`
	GetRID        BooleanYN `json:"get_rid,omitempty" query:"get_rid"`
	FieldJoinType JoinType  `json:"field_join_type,omitempty" query:"field_join_type"` // "or" or "and"
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
	Text              string          `json:"text,omitempty" query:"text"`
	StringJoinType    JoinType        `json:"string_join_type,omitempty" query:"string_join_type"`
	Keywords          BooleanYN       `json:"keywords,omitempty" query:"keywords"`
	Title             BooleanYN       `json:"title,omitempty" query:"title"`
	Description       BooleanYN       `json:"description,omitempty" query:"description"`
	MD5               BooleanYN       `json:"md5,omitempty" query:"md5"`
	KeywordID         string          `json:"keyword_id,omitempty" query:"keyword_id"`
	Username          string          `json:"username,omitempty" query:"username"`
	UserID            string          `json:"user_id,omitempty" query:"user_id"`
	FavsUserID        string          `json:"favs_user_id,omitempty" query:"favs_user_id"`
	UnreadSubmissions BooleanYN       `json:"unread_submissions,omitempty" query:"unread_submissions"`
	Type              SubmissionTypes `json:"type,omitempty" query:"type"`
	Sales             string          `json:"sales,omitempty" query:"sales"` // Values: forsale, digital, prints
	PoolID            string          `json:"pool_id,omitempty" query:"pool_id"`
	OrderBy           string          `json:"orderby,omitempty" query:"orderby"` // Values: create_datetime, unread_datetime, views, total_print_sales, total_digital_sales, total_sales, username, fav_datetime, fav_stars, pool_order. Default: create_datetime.
	DaysLimit         IntString       `json:"dayslimit,omitempty" query:"dayslimit"`
	Random            BooleanYN       `json:"random,omitempty" query:"random"`
	// Scraps Set how submissions marked as “Scraps” are returned.
	// Possible values are:
	// 	both – show submissions from Scraps and Main galleries.
	// 	no – Do not show Scraps. Shows only submissions from Main galleries.
	// 	only – Show only submissions from Scraps galleries, not Main galleries.
	Scraps     string    `json:"scraps,omitempty" query:"scraps"`
	CountLimit IntString `json:"count_limit,omitempty" query:"count_limit"`
}

type SubmissionSearchResponse struct {
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

type SubmissionType int

type SubmissionTypes []SubmissionType

func (s SubmissionTypes) MarshalJSON() ([]byte, error) {
	var types bytes.Buffer
	types.WriteRune('"')
	for i, t := range s {
		if i > 0 {
			types.WriteString(",")
		}
		types.WriteString(fmt.Sprintf("%d", t))
	}
	types.WriteRune('"')
	return types.Bytes(), nil
}

func (s *SubmissionTypes) UnmarshalJSON(data []byte) error {
	var types []SubmissionType
	for _, t := range strings.Split(string(data), ",") {
		i, err := strconv.Atoi(t)
		if err != nil {
			return fmt.Errorf("failed to parse submission type: %w", err)
		}
		types = append(types, SubmissionType(i))
	}
	*s = types
	return nil
}

const (
	SubmissionTypePicturePinup              SubmissionType = iota + 1 //1 - Picture/Pinup
	SubmissionTypeSketch                                              //2 - Sketch
	SubmissionTypePictureSeries                                       //3 - Picture Series
	SubmissionTypeComic                                               //4 - Comic
	SubmissionTypePortfolio                                           //5 - Portfolio
	SubmissionTypeShockwaveFlashAnimation                             //6 - Shockwave/Flash - Animation
	SubmissionTypeShockwaveFlashInteractive                           //7 - Shockwave/Flash - Interactive
	SubmissionTypeVideoFeatureLength                                  //8 - Video - Feature Length
	SubmissionTypeVideoAnimation3DCGI                                 //9 - Video - Animation/3D/CGI
	SubmissionTypeMusicSingleTrack                                    //10 - Music - Single Track
	SubmissionTypeMusicAlbum                                          //11 - Music - Album
	SubmissionTypeWritingDocument                                     //12 - Writing - Document
	SubmissionTypeCharacterSheet                                      //13 - Character Sheet
	SubmissionTypePhotography                                         //14 - Photography - Fursuit/Sculpture/Jewelry/etc
)

func (user Credentials) SearchSubmissions(req SubmissionSearchRequest) (SubmissionSearchResponse, error) {
	if !user.LoggedIn() {
		return SubmissionSearchResponse{}, ErrNotLoggedIn
	}
	if req.SID == "" {
		req.SID = user.Sid
	}

	resp, err := user.Get(ApiUrl("search", utils.StructToUrlValues(req)))

	if err != nil {
		return SubmissionSearchResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SubmissionSearchResponse{}, err
	}

	if err := CheckError(body); err != nil {
		return SubmissionSearchResponse{}, fmt.Errorf("error searching submissions: %w", err)
	}

	var searchResp SubmissionSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return SubmissionSearchResponse{}, err
	}

	return searchResp, nil
}

func (user Credentials) OwnSubmissions() (SubmissionSearchResponse, error) {
	return user.SearchSubmissions(SubmissionSearchRequest{SID: user.Sid, Username: user.Username})
}

func (user Credentials) UserSubmissions(username string) (SubmissionSearchResponse, error) {
	return user.SearchSubmissions(SubmissionSearchRequest{SID: user.Sid, Username: username})
}
