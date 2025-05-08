package inkbunny

import (
	"bytes"
	"fmt"
	"iter"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ellypaws/inkbunny/types"
)

type SubmissionSearchRequest struct {
	SID        string           `json:"sid" query:"sid"`
	OutputMode types.OutputMode `json:"output_mode,omitempty" query:"output_mode"`
	// Setting an RID uses Mode 2: Page through results.
	//
	// Using a Results ID (RID), which can be returned by any search in Mode 1, you can specify a results set to page through.
	// This means you can run the search once in Mode 1, and then return the results page by page without running the search again.
	// This is much faster than using Mode 1 over and over again to get each subsequent page of results.
	//
	// The disadvantage is that the results are not updated between requests.
	// If submissions are added, changed or deleted in a way that would alter the current search results, you will not see the change until the search is run again with Mode 1.
	// Note: If you specify an RID, the search script will ignore all "Mode 1" search parameters.
	//
	// Results ID of results set to page through. To get an RID to use here, run a search (Mode 1, as described above) with GetRID: Yes first.
	// Then use the returned RID here to page through those results without needing to run the same search again. Only used for Mode 2, as described above.
	// Note: Results sets will be automatically removed after not being accessed for a certain amount of time, or if an excess number of results sets are created by a user.
	// Attempting to access a results set that has been removed will throw an error. See the [Error Codes] section in this document.
	//
	// [Error Codes]: https://wiki.inkbunny.net/wiki/API#Error_Codes
	RID                string          `json:"rid,omitempty" query:"rid"`
	SubmissionIDsOnly  types.BooleanYN `json:"submission_ids_only,omitempty" query:"submission_ids_only"`
	SubmissionsPerPage types.IntString `json:"submissions_per_page,omitempty" query:"submissions_per_page"`
	Page               types.IntString `json:"page,omitempty" query:"page"` // Results page number to return. Default: 1.
	// Not to be confused with Text. This is a boolean value to return list of Top 100 Keywords.
	// Return list of Top 100 Keywords associated with all submissions on current results page.
	// Note that this list includes both officially assigned keywords and also keywords
	// suggested for this submission by other users.
	KeywordsList types.BooleanYN `json:"keywords_list,omitempty" query:"keywords_list"`
	// Skip returning submission info.
	// Useful when you are just returning Top Keywords or Submission Counts for searches, and you don't want all the other submission data.
	NoSubmissions types.BooleanYN `json:"no_submissions,omitempty" query:"no_submissions"`
	// Return a Results ID for this search, which can then be used in Mode 2 (By setting the RID) to page through the results without running the search again for each page.
	GetRID        types.BooleanYN     `json:"get_rid,omitempty" query:"get_rid"`
	FieldJoinType types.FieldJoinType `json:"field_join_type,omitempty" query:"field_join_type"` // "or" or "and"
	// Text to search chosen fields for. eg "dragon", "wolf", etc.
	// A Full Text search is performed using this string (see the meaning of Full Text searches in the Postgresql Documentation).
	// The characters "_" and "," are converted to spaces automatically.
	// Characters which have special meanings for Full Text searches in Postgresql (such as |, &, :, ! and ~) are ignored.
	//
	// Note: At least one of the Search Field parameters Keywords, Title, Description or MD5 must be set to Yes for text search to work.
	// By default, Keywords is set to Yes, so all searches with no Search Field specified will search in keywords.
	// Values: (Any text string).
	//
	// Default: n/a. Required: No
	Text           string         `json:"text,omitempty" query:"text"`
	StringJoinType types.JoinType `json:"string_join_type,omitempty" query:"string_join_type"`
	// Search Keywords for the chosen text.
	// Note: This is ON (Yes) by default, and is the standard field that text searches look in, unless specified otherwise.
	// Note: At least one of keywords, title or description must be set to Yes for text search to work.
	Keywords types.BooleanYN `json:"keywords,omitempty" query:"keywords"`
	// Search Title for the chosen text.
	// Note: At least one of keywords, title or description must be set to Yes for text search to work.
	Title types.BooleanYN `json:"title,omitempty" query:"title"`
	// Search the Description AND Story fields for the chosen text.
	// Note: At least one of keywords, title or description must be set to Yes for text search to work.
	Description types.BooleanYN `json:"description,omitempty" query:"description"`
	// Search for the chosen text in the MD5 Checksum/hash of the Initial.
	// (as uploaded and before any conversion), Full (may have metadata removed and
	// optimised for lossless compression), Large (also known as Screen), Small, or
	// HiRes/Sales versions of a file.
	// This is useful for finding files based on their content, and for finding identical files.
	//
	//	* Although the MD5 Hash for the HiRes/Sales file is only shown to the Submission owner, they are still found when anyone runs an MD5 search.
	//	* Deleted Files - This search will also find submissions based the MD5 of free and sales files that are marked deleted (that have been removed from a submission). This is to assist with finding submissions even if their files are updated later.
	//	* The property "string_join_type" has no effect on searching for MD5 strings. It always assumes Or when multiple MD5 Hashes are given.
	//	* See [MD5 Checksums] for more information on how MD5 is used in Inkbunny.
	//
	// [MD5 Checksums]: https://wiki.inkbunny.net/wiki/MD5
	MD5               types.BooleanYN   `json:"md5,omitempty" query:"md5"`
	KeywordID         types.IntString   `json:"keyword_id,omitempty" query:"keyword_id"`
	Username          string            `json:"username,omitempty" query:"username"`
	UserID            types.IntString   `json:"user_id,omitempty" query:"user_id"`
	FavsUserID        types.IntString   `json:"favs_user_id,omitempty" query:"favs_user_id"`
	UnreadSubmissions types.BooleanYN   `json:"unread_submissions,omitempty" query:"unread_submissions"`
	Type              SubmissionTypes   `json:"type,omitempty" query:"type"`
	Sales             types.SalesFilter `json:"sales,omitempty" query:"sales"` // Values: forsale, digital, prints
	PoolID            types.IntString   `json:"pool_id,omitempty" query:"pool_id"`
	OrderBy           types.OrderBy     `json:"orderby,omitempty" query:"orderby"` // Values: create_datetime, unread_datetime, views, total_print_sales, total_digital_sales, total_sales, username, fav_datetime, fav_stars, pool_order. Default: create_datetime.
	DaysLimit         types.IntString   `json:"dayslimit,omitempty" query:"dayslimit"`
	// Sort results randomly. This is done after all other filters and sort orders
	// are applied. This can be used in conjunction with "orderby". You can order
	// results with OrderBy, limit the number returned with other filters like
	// CountLimit, and then if Random: Yes it will sort those results randomly.
	// Eg: Set OrderBy: OrderByViews and CountLimit: 100 to get the top 100 submissions,
	// then with "random=yes" those top 100 are sorted randomly AFTER the other
	// limits and conditions are used. Does your head hurt? Mine does.
	Random types.BooleanYN `json:"random,omitempty" query:"random"`
	// Scraps Set how submissions marked as “Scraps” are returned.
	// Possible values are:
	// 	both – show submissions from Scraps and Main galleries.
	// 	no – Do not show Scraps. Shows only submissions from Main galleries.
	// 	only – Show only submissions from Scraps galleries, not Main galleries.
	Scraps     string          `json:"scraps,omitempty" query:"scraps"`
	CountLimit types.IntString `json:"count_limit,omitempty" query:"count_limit"`
}

type SubmissionSearchResponse struct {
	SID                  string             `json:"sid"`
	UserLocation         string             `json:"user_location"`
	ResultsCountAll      types.IntString    `json:"results_count_all"`
	ResultsCountThisPage types.IntString    `json:"results_count_thispage"`
	PagesCount           types.IntString    `json:"pages_count"`
	Page                 types.IntString    `json:"page"`
	RID                  string             `json:"rid,omitempty"`
	RIDTTL               string             `json:"rid_ttl,omitempty"`
	RIDTTLDuration       time.Duration      `json:"-"`
	RIDExpiry            time.Time          `json:"-"`
	SearchParams         []SearchParam      `json:"search_params"`
	KeywordList          []KeywordList      `json:"keyword_list,omitempty"`
	Submissions          []SubmissionSearch `json:"submissions,omitempty"`
	client               *Client
}

type KeywordList struct {
	KeywordID        types.IntString `json:"keyword_id"`
	KeywordName      string          `json:"keyword_name"`
	SubmissionsCount types.IntString `json:"submissions_count"`
}

type SubmissionSearch struct {
	SubmissionBasic
	UnreadDateSystem string          `json:"unread_datetime_system,omitempty"`
	UnreadDateUser   string          `json:"unread_datetime,omitempty"`
	Updated          types.BooleanYN `json:"updated,omitempty"`
	Stars            types.IntString `json:"stars,omitempty"`
}

// SearchParam is the search parameters that were used to find these search results.
type SearchParam struct {
	Name string `json:"param_name"`
	Type string `json:"param_type"`
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
				Page: i,
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

func (u *User) SearchSubmissions(req SubmissionSearchRequest) (SubmissionSearchResponse, error) {
	if req.SID == "" {
		if u.SID == "" {
			return SubmissionSearchResponse{}, ErrEmptySID
		}
		req.SID = u.SID
	}

	return u.Client().SearchSubmissions(req)
}

func (c *Client) SearchSubmissions(req SubmissionSearchRequest) (SubmissionSearchResponse, error) {
	response, err := PostDecode[SubmissionSearchResponse](c, ApiUrl("search"), req)
	if err != nil {
		return response, err
	}

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
