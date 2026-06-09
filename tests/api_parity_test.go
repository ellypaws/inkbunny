package tests

import (
	"encoding/json"
	"fmt"
	"iter"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ellypaws/inkbunny"
)

type rewriteTransport struct {
	baseURL *url.URL
	next    http.RoundTripper
}

func (r rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	clone.URL.Scheme = r.baseURL.Scheme
	clone.URL.Host = r.baseURL.Host
	return r.next.RoundTrip(clone)
}

func newTestClient(t *testing.T, handler http.HandlerFunc) *inkbunny.Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse server URL: %v", err)
	}

	return inkbunny.NewClient(inkbunny.WithClient(&http.Client{Transport: rewriteTransport{
		baseURL: serverURL,
		next:    server.Client().Transport,
	}}))
}

func parseAPIForm(t *testing.T, r *http.Request) url.Values {
	t.Helper()
	if err := r.ParseMultipartForm(1 << 20); err != nil && !strings.Contains(err.Error(), "multipart") {
		t.Fatalf("parse multipart form: %v", err)
	}
	if err := r.ParseForm(); err != nil {
		t.Fatalf("parse form: %v", err)
	}
	return r.Form
}

func cloneValues(values url.Values) url.Values {
	clone := make(url.Values, len(values))
	for key, items := range values {
		clone[key] = append([]string(nil), items...)
	}
	return clone
}

func writeJSON(t *testing.T, w http.ResponseWriter, value any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		t.Fatalf("encode response: %v", err)
	}
}

func searchPage(sid, rid string, page, pages, perPage int, ids ...int) map[string]any {
	submissions := make([]map[string]string, 0, len(ids))
	for _, id := range ids {
		submissions = append(submissions, map[string]string{"submission_id": fmt.Sprint(id)})
	}
	return map[string]any{
		"sid":                    sid,
		"rid":                    rid,
		"rid_ttl":                "2 hours",
		"results_count_all":      fmt.Sprint(pages * perPage),
		"results_count_thispage": fmt.Sprint(len(submissions)),
		"pages_count":            fmt.Sprint(pages),
		"page":                   fmt.Sprint(page),
		"user_location":          "",
		"search_params":          []map[string]any{},
		"submissions":            submissions,
	}
}

func collectSearchPages(t *testing.T, seq iter.Seq2[inkbunny.SubmissionSearchResponse, error], limit int) []inkbunny.SubmissionSearchResponse {
	t.Helper()
	pages := []inkbunny.SubmissionSearchResponse{}
	for page, err := range seq {
		if err != nil {
			t.Fatalf("collect page: %v", err)
		}
		pages = append(pages, page)
		if limit > 0 && len(pages) == limit {
			break
		}
	}
	return pages
}

func submissionSignature(submissions []inkbunny.SubmissionSearch) string {
	ids := make([]string, 0, len(submissions))
	for _, submission := range submissions {
		ids = append(ids, submission.SubmissionID.String())
	}
	return strings.Join(ids, ",")
}

func TestSearchRequestAllPagesUsesReturnedRID(t *testing.T) {
	var requests []url.Values
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		form := parseAPIForm(t, r)
		requests = append(requests, cloneValues(form))

		switch form.Get("page") {
		case "", "1":
			if form.Get("rid") != "" {
				t.Fatalf("initial request sent rid=%q", form.Get("rid"))
			}
			if form.Get("get_rid") != "yes" {
				t.Fatalf("initial get_rid=%q, want yes", form.Get("get_rid"))
			}
			writeJSON(t, w, searchPage("sid", "rid-1", 1, 3, 3, 1, 2, 3))
		case "2":
			if form.Get("rid") != "rid-1" {
				t.Fatalf("page 2 rid=%q, want rid-1", form.Get("rid"))
			}
			if form.Get("submissions_per_page") != "3" {
				t.Fatalf("page 2 submissions_per_page=%q, want 3", form.Get("submissions_per_page"))
			}
			writeJSON(t, w, searchPage("sid", "rid-1", 2, 3, 3, 4, 5, 6))
		case "3":
			if form.Get("rid") != "rid-1" {
				t.Fatalf("page 3 rid=%q, want rid-1", form.Get("rid"))
			}
			writeJSON(t, w, searchPage("sid", "rid-1", 3, 3, 3, 7, 8, 9))
		default:
			t.Fatalf("unexpected page %q", form.Get("page"))
		}
	})

	oldDefault := inkbunny.DefaultClient
	inkbunny.DefaultClient = client
	t.Cleanup(func() { inkbunny.DefaultClient = oldDefault })

	pages := collectSearchPages(t, inkbunny.SubmissionSearchRequest{
		SID:                "sid",
		Text:               "cub",
		OrderBy:            inkbunny.OrderByViews,
		GetRID:             inkbunny.Yes,
		Page:               1,
		SubmissionsPerPage: 3,
	}.AllPages(), 0)

	if len(pages) != 3 {
		t.Fatalf("got %d pages, want 3", len(pages))
	}
	if got := submissionSignature(pages[0].Submissions); got != "1,2,3" {
		t.Fatalf("page 1 ids=%q", got)
	}
	if got := submissionSignature(pages[1].Submissions); got != "4,5,6" {
		t.Fatalf("page 2 ids=%q", got)
	}
	if got := submissionSignature(pages[2].Submissions); got != "7,8,9" {
		t.Fatalf("page 3 ids=%q", got)
	}
	if len(requests) != 3 {
		t.Fatalf("got %d requests, want 3", len(requests))
	}
}

func TestSearchResponseAllPagesUsesReturnedRIDAndStopsEarly(t *testing.T) {
	requestCount := 0
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		form := parseAPIForm(t, r)
		switch form.Get("page") {
		case "", "1":
			writeJSON(t, w, searchPage("sid", "rid-2", 1, 4, 2, 10, 11))
		case "2":
			if form.Get("rid") != "rid-2" {
				t.Fatalf("page 2 rid=%q, want rid-2", form.Get("rid"))
			}
			writeJSON(t, w, searchPage("sid", "rid-2", 2, 4, 2, 12, 13))
		default:
			t.Fatalf("iterator did not stop early, requested page %q", form.Get("page"))
		}
	})

	first, err := client.SearchSubmissions(inkbunny.SubmissionSearchRequest{
		SID:                "sid",
		GetRID:             inkbunny.Yes,
		SubmissionsPerPage: 2,
	})
	if err != nil {
		t.Fatalf("initial search: %v", err)
	}

	pages := collectSearchPages(t, first.AllPages(), 2)
	if len(pages) != 2 {
		t.Fatalf("got %d pages, want 2", len(pages))
	}
	if requestCount != 2 {
		t.Fatalf("got %d requests, want 2", requestCount)
	}
	if got := submissionSignature(pages[1].Submissions); got != "12,13" {
		t.Fatalf("page 2 ids=%q", got)
	}
}

func TestSearchResponseAllPagesErrorAndEmptyPage(t *testing.T) {
	t.Run("error propagation", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			form := parseAPIForm(t, r)
			if form.Get("page") == "2" {
				writeJSON(t, w, inkbunny.ErrorResponse{Code: newInt(4), Message: "rid expired"})
				return
			}
			writeJSON(t, w, searchPage("sid", "rid-error", 1, 2, 1, 1))
		})

		first, err := client.SearchSubmissions(inkbunny.SubmissionSearchRequest{SID: "sid", GetRID: inkbunny.Yes, SubmissionsPerPage: 1})
		if err != nil {
			t.Fatalf("initial search: %v", err)
		}

		var errs []error
		for _, err := range first.AllPages() {
			if err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) != 1 {
			t.Fatalf("got %d errors, want 1", len(errs))
		}
	})

	t.Run("empty page", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			form := parseAPIForm(t, r)
			if form.Get("page") == "2" {
				writeJSON(t, w, searchPage("sid", "rid-empty", 2, 2, 1))
				return
			}
			writeJSON(t, w, searchPage("sid", "rid-empty", 1, 2, 1, 1))
		})

		first, err := client.SearchSubmissions(inkbunny.SubmissionSearchRequest{SID: "sid", GetRID: inkbunny.Yes, SubmissionsPerPage: 1})
		if err != nil {
			t.Fatalf("initial search: %v", err)
		}
		pages := collectSearchPages(t, first.AllPages(), 0)
		if len(pages) != 2 {
			t.Fatalf("got %d pages, want 2", len(pages))
		}
		if len(pages[1].Submissions) != 0 {
			t.Fatalf("empty page had %d submissions", len(pages[1].Submissions))
		}
	})
}

func TestAllSubmissionsAndAllDetails(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		form := parseAPIForm(t, r)
		switch r.URL.Path {
		case "/api_search.php":
			if form.Get("page") == "2" {
				writeJSON(t, w, searchPage("sid", "rid-details", 2, 2, 2, 3))
				return
			}
			writeJSON(t, w, searchPage("sid", "rid-details", 1, 2, 2, 1, 2))
		case "/api_submissions.php":
			ids := strings.Split(form.Get("submission_ids"), ",")
			submissions := make([]map[string]string, 0, len(ids))
			for _, id := range ids {
				submissions = append(submissions, map[string]string{"submission_id": id})
			}
			writeJSON(t, w, map[string]any{
				"sid":           "sid",
				"results_count": fmt.Sprint(len(submissions)),
				"user_location": "",
				"submissions":   submissions,
			})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	})

	oldDefault := inkbunny.DefaultClient
	inkbunny.DefaultClient = client
	t.Cleanup(func() { inkbunny.DefaultClient = oldDefault })

	req := inkbunny.SubmissionSearchRequest{SID: "sid", GetRID: inkbunny.Yes, SubmissionsPerPage: 2}
	var batches [][]inkbunny.SubmissionSearch
	for submissions, err := range req.AllSubmissions() {
		if err != nil {
			t.Fatalf("AllSubmissions: %v", err)
		}
		batches = append(batches, submissions)
	}
	if len(batches) != 2 {
		t.Fatalf("got %d submission batches, want 2", len(batches))
	}
	if got := submissionSignature(batches[1]); got != "3" {
		t.Fatalf("second batch ids=%q", got)
	}

	var details []inkbunny.SubmissionDetailsResponse
	for detail, err := range req.AllDetails(inkbunny.SubmissionDetailsRequest{ShowPools: inkbunny.Yes}) {
		if err != nil {
			t.Fatalf("AllDetails: %v", err)
		}
		details = append(details, detail)
	}
	if len(details) != 2 {
		t.Fatalf("got %d detail batches, want 2", len(details))
	}
	if got := details[0].Submissions[0].SubmissionID.String(); got != "1" {
		t.Fatalf("first detail id=%q", got)
	}
	if got := details[1].Submissions[0].SubmissionID.String(); got != "3" {
		t.Fatalf("second detail id=%q", got)
	}
}

func TestEditSubmissionKeywordMarshalling(t *testing.T) {
	tests := []struct {
		name       string
		keywords   []string
		wantKey    bool
		wantValue  string
		wantPublic *inkbunny.BooleanYN
		wantNotify *inkbunny.BooleanYN
		wantVis    string
	}{
		{name: "nil omitted", keywords: nil, wantKey: false},
		{name: "empty clears", keywords: []string{}, wantKey: true, wantValue: ""},
		{name: "single keyword", keywords: []string{"fox"}, wantKey: true, wantValue: "fox"},
		{name: "multiple keywords", keywords: []string{"Roger_Rabbit", "space"}, wantKey: true, wantValue: "Roger_Rabbit,space"},
		{name: "public without notify", keywords: []string{"fox"}, wantKey: true, wantValue: "fox", wantPublic: &inkbunny.Yes, wantNotify: &inkbunny.No, wantVis: "yes_nowatch"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				form := parseAPIForm(t, r)
				values, ok := form["keywords"]
				if ok != tt.wantKey {
					t.Fatalf("keywords present=%v, want %v; form=%v", ok, tt.wantKey, form)
				}
				if tt.wantKey && len(values) > 0 && values[0] != tt.wantValue {
					t.Fatalf("keywords=%q, want %q", values[0], tt.wantValue)
				}
				if tt.wantVis != "" && form.Get("visibility") != tt.wantVis {
					t.Fatalf("visibility=%q, want %q", form.Get("visibility"), tt.wantVis)
				}
				writeJSON(t, w, map[string]any{"submission_id": "123", "twitter_authentication_success": "f"})
			})

			_, err := client.EditSubmission(inkbunny.SubmissionEditRequest{
				SID:          "sid",
				SubmissionID: 123,
				Keywords:     tt.keywords,
				Public:       tt.wantPublic,
				Notify:       tt.wantNotify,
			})
			if err != nil {
				t.Fatalf("EditSubmission: %v", err)
			}
		})
	}
}

func TestUploadNotifyMarshalsBooleanYN(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   inkbunny.BooleanYN
		want string
	}{
		{name: "yes", in: inkbunny.Yes, want: "yes"},
		{name: "no", in: inkbunny.No, want: "no"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			body, contentType, err := inkbunny.structToMultipartForm(inkbunny.UploadRequest{SID: "sid", Notify: tt.in})
			if err != nil {
				t.Fatalf("multipart form: %v", err)
			}
			_, params, err := mime.ParseMediaType(contentType)
			if err != nil {
				t.Fatalf("parse media type: %v", err)
			}
			reader := multipart.NewReader(body, params["boundary"])
			form, err := reader.ReadForm(1 << 20)
			if err != nil {
				t.Fatalf("read form: %v", err)
			}
			if got := form.Value["notify"][0]; got != tt.want {
				t.Fatalf("notify=%q, want %q", got, tt.want)
			}
		})
	}
}

func TestMembersAndWatchlistOptions(t *testing.T) {
	t.Run("search members options", func(t *testing.T) {
		var seen []url.Values
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			form := parseAPIForm(t, r)
			seen = append(seen, cloneValues(form))
			writeJSON(t, w, map[string]any{"results": []map[string]string{{"id": "1", "value": "tester"}}})
		})

		if _, err := client.SearchMembersWithOptions(inkbunny.SearchMembersRequest{Username: "est", SearchType: inkbunny.UsernameSearchAny}); err != nil {
			t.Fatalf("SearchMembersWithOptions: %v", err)
		}
		if _, err := client.SearchMembers("legacy"); err != nil {
			t.Fatalf("SearchMembers: %v", err)
		}
		if seen[0].Get("username") != "est" || seen[0].Get("searchtype") != "any" {
			t.Fatalf("options form=%v", seen[0])
		}
		if seen[1].Get("username") != "legacy" || seen[1].Get("searchtype") != "" {
			t.Fatalf("legacy form=%v", seen[1])
		}
	})

	t.Run("watchlist response", func(t *testing.T) {
		var seen []url.Values
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			form := parseAPIForm(t, r)
			seen = append(seen, cloneValues(form))
			writeJSON(t, w, map[string]any{
				"sid":           form.Get("sid"),
				"results_count": "1",
				"watches":       []map[string]string{{"user_id": "7", "username": "watched"}},
			})
		})
		user := &inkbunny.User{SID: "sid", client: client}

		response, err := user.GetWatchingResponse(inkbunny.WatchlistRequest{OrderBy: inkbunny.WatchlistOrderByAlphabetical, Limit: 5})
		if err != nil {
			t.Fatalf("GetWatchingResponse: %v", err)
		}
		if response.SID != "sid" || response.ResultsCount != 1 || response.Watches[0].Username != "watched" {
			t.Fatalf("unexpected response: %#v", response)
		}
		watches, err := user.GetWatching()
		if err != nil {
			t.Fatalf("GetWatching: %v", err)
		}
		if len(watches) != 1 || watches[0].UserID != "7" {
			t.Fatalf("unexpected watches: %#v", watches)
		}
		if seen[0].Get("orderby") != "alphabetical" || seen[0].Get("limit") != "5" {
			t.Fatalf("watchlist options form=%v", seen[0])
		}
		if seen[1].Get("orderby") != "" || seen[1].Get("limit") != "" {
			t.Fatalf("legacy watchlist form=%v", seen[1])
		}
	})
}

func TestOrderByConstants(t *testing.T) {
	tests := map[string]inkbunny.OrderBy{
		"create_datetime":           inkbunny.OrderByCreateDatetime,
		"last_file_update_datetime": inkbunny.OrderByLastFileUpdateDatetime,
		"unread_datetime":           inkbunny.OrderByUnreadDatetime,
		"unread_datetime_reverse":   inkbunny.OrderByUnreadDatetimeReverse,
		"views":                     inkbunny.OrderByViews,
		"total_print_sales":         inkbunny.OrderByTotalPrint,
		"total_digital_sales":       inkbunny.OrderByTotalDigital,
		"total_sales":               inkbunny.OrderByTotalSales,
		"username":                  inkbunny.OrderByUsername,
		"fav_datetime":              inkbunny.OrderByFavDatetime,
		"fav_stars":                 inkbunny.OrderByFavStars,
		"pool_order":                inkbunny.OrderByPoolOrder,
	}
	for want, got := range tests {
		if string(got) != want {
			t.Fatalf("constant value=%q, want %q", got, want)
		}
	}
}

func TestLiveGuestRIDPaging(t *testing.T) {
	user := liveGuestUser(t)

	first, err := user.SearchSubmissions(inkbunny.SubmissionSearchRequest{
		Text:               "cub",
		OrderBy:            inkbunny.OrderByViews,
		GetRID:             inkbunny.Yes,
		SubmissionsPerPage: 3,
		CountLimit:         30,
	})
	if err != nil {
		t.Fatalf("live search: %v", err)
	}
	if first.RID == "" {
		t.Fatal("live search did not return RID")
	}

	seen := map[string]bool{}
	pageCount := 0
	for page, err := range first.AllPages() {
		if err != nil {
			t.Fatalf("live page: %v", err)
		}
		if page.Page.Int() != pageCount+1 {
			t.Fatalf("page number=%d, want %d", page.Page.Int(), pageCount+1)
		}
		if len(page.Submissions) == 0 {
			t.Fatalf("page %d returned no submissions", page.Page.Int())
		}
		signature := submissionSignature(page.Submissions)
		if seen[signature] {
			t.Fatalf("duplicate live page signature %q", signature)
		}
		seen[signature] = true
		pageCount++
		if pageCount == 3 || pageCount == first.PagesCount.Int() {
			break
		}
	}
	if pageCount < 1 {
		t.Fatal("no live pages collected")
	}
}

func TestLiveGuestOnlineAPIs(t *testing.T) {
	user := liveGuestUser(t)

	memberMatches, err := user.SearchMembers("ell")
	if err != nil {
		t.Fatalf("SearchMembers: %v", err)
	}
	if memberMatches == nil {
		t.Fatal("SearchMembers returned nil slice")
	}

	memberOptions, err := user.SearchMembersWithOptions(inkbunny.SearchMembersRequest{
		Username:   "ell",
		SearchType: inkbunny.UsernameSearchAny,
	})
	if err != nil {
		t.Fatalf("SearchMembersWithOptions: %v", err)
	}
	if memberOptions == nil {
		t.Fatal("SearchMembersWithOptions returned nil slice")
	}

	keywords, err := user.KeywordSuggestion("cub", inkbunny.ParseMaskU(inkbunny.General|inkbunny.Nudity|inkbunny.MildViolence|inkbunny.Sexual|inkbunny.StrongViolence), false)
	if err != nil {
		t.Fatalf("KeywordSuggestion: %v", err)
	}
	if len(keywords) == 0 {
		t.Fatal("KeywordSuggestion returned no results")
	}

	search, err := user.SearchSubmissions(inkbunny.SubmissionSearchRequest{
		Text:               "cub",
		OrderBy:            inkbunny.OrderByViews,
		GetRID:             inkbunny.Yes,
		SubmissionsPerPage: 5,
		CountLimit:         25,
	})
	if err != nil {
		t.Fatalf("SearchSubmissions: %v", err)
	}
	if len(search.Submissions) == 0 {
		t.Fatal("SearchSubmissions returned no submissions")
	}

	details, err := search.Details()
	if err != nil {
		t.Fatalf("SubmissionSearchResponse.Details: %v", err)
	}
	if len(details.Submissions) == 0 {
		t.Fatal("Details returned no submissions")
	}
	if details.Submissions[0].SubmissionID == 0 {
		t.Fatalf("Details returned empty submission id: %#v", details.Submissions[0])
	}

	watchlist, err := user.GetWatchingResponse(inkbunny.WatchlistRequest{
		OrderBy: inkbunny.WatchlistOrderByAlphabetical,
		Limit:   5,
	})
	if err != nil {
		t.Fatalf("GetWatchingResponse: %v", err)
	}
	if watchlist.SID == "" {
		t.Fatal("GetWatchingResponse returned empty sid")
	}
	legacyWatchlist, err := user.GetWatching()
	if err != nil {
		t.Fatalf("GetWatching: %v", err)
	}
	if legacyWatchlist == nil {
		t.Fatal("GetWatching returned nil slice")
	}
}

func TestLiveAuthenticatedUpload(t *testing.T) {
	if os.Getenv("INKBUNNY_LIVE_UPLOAD_TESTS") != "true" {
		t.Skip("set INKBUNNY_LIVE_UPLOAD_TESTS=true to run live upload checks")
	}

	username := os.Getenv("INKBUNNY_USER")
	password := os.Getenv("INKBUNNY_PASSWORD")
	if username == "" || password == "" {
		t.Skip("set INKBUNNY_USER and INKBUNNY_PASSWORD for live upload checks")
	}

	paths := strings.Split(os.Getenv("INKBUNNY_UPLOAD_FILES"), string(os.PathListSeparator))
	if len(paths) == 0 || strings.TrimSpace(paths[0]) == "" {
		t.Skip("set INKBUNNY_UPLOAD_FILES to one or more local files separated by the OS path-list separator")
	}

	user, err := inkbunny.NewClient().Login(username, password)
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	t.Cleanup(func() {
		if err := user.Logout(); err != nil {
			t.Logf("logout: %v", err)
		}
	})

	files := make([]inkbunny.FileUpload, 0, len(paths))
	opened := make([]*os.File, 0, len(paths))
	for _, rawPath := range paths {
		path := strings.TrimSpace(rawPath)
		if path == "" {
			continue
		}
		file, err := os.Open(path)
		if err != nil {
			t.Fatalf("open upload file %q: %v", path, err)
		}
		opened = append(opened, file)
		files = append(files, inkbunny.FileUpload{MainFile: &inkbunny.FileContent{
			Name: filepath.Base(path),
			File: file,
		}})
	}
	t.Cleanup(func() {
		for _, file := range opened {
			if err := file.Close(); err != nil {
				t.Logf("close upload file: %v", err)
			}
		}
	})
	if len(files) == 0 {
		t.Skip("INKBUNNY_UPLOAD_FILES did not contain any usable paths")
	}

	response, err := user.Upload(inkbunny.UploadRequest{
		Notify: inkbunny.No,
		Files:  files,
	})
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if response.SubmissionID == "" {
		t.Fatal("Upload returned empty submission id")
	}
	t.Cleanup(func() {
		if err := response.Delete(); err != nil {
			t.Logf("delete uploaded submission %s: %v", response.SubmissionID, err)
		}
	})
}

func liveGuestUser(t *testing.T) *inkbunny.User {
	t.Helper()
	if os.Getenv("INKBUNNY_LIVE_TESTS") != "true" {
		t.Skip("set INKBUNNY_LIVE_TESTS=true to run live Inkbunny API checks")
	}

	user, err := inkbunny.NewClient().Login("guest", "")
	if err != nil {
		t.Fatalf("guest login: %v", err)
	}
	t.Cleanup(func() {
		if err := user.Logout(); err != nil {
			t.Logf("guest logout: %v", err)
		}
	})

	if err := user.ChangeRatings(inkbunny.Ratings{
		Nudity:         &inkbunny.Yes,
		MildViolence:   &inkbunny.Yes,
		Sexual:         &inkbunny.Yes,
		StrongViolence: &inkbunny.Yes,
	}); err != nil {
		t.Fatalf("change ratings: %v", err)
	}
	return user
}

func newInt(v int) *int {
	return &v
}
