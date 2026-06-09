package inkbunny

type UsernameSearchType = string

const (
	// UsernameSearchStart matches usernames that start with the query.
	UsernameSearchStart UsernameSearchType = "start"
	// UsernameSearchAny matches usernames containing the query anywhere.
	UsernameSearchAny UsernameSearchType = "any"
)

type SearchMembersRequest struct {
	Username   string             `json:"username" query:"username"`
	SearchType UsernameSearchType `json:"searchtype,omitempty" query:"searchtype"`
}

type WatchlistOrderBy = string

const (
	// WatchlistOrderByCreateDatetime sorts by the date the account was watched, newest first.
	WatchlistOrderByCreateDatetime WatchlistOrderBy = "create_datetime"
	// WatchlistOrderByAlphabetical sorts watched accounts alphabetically.
	WatchlistOrderByAlphabetical WatchlistOrderBy = "alphabetical"
)

type WatchlistRequest struct {
	SID     string           `json:"sid" query:"sid"`
	OrderBy WatchlistOrderBy `json:"orderby,omitempty" query:"orderby"`
	Limit   IntString        `json:"limit,omitempty" query:"limit"`
}

type WatchlistResponse struct {
	SID          string       `json:"sid"`
	ResultsCount IntString    `json:"results_count"`
	Watches      []UsernameID `json:"watches"`
}

func (c *Client) SearchMembers(username string) ([]Autocomplete, error) {
	return c.SearchMembersWithOptions(SearchMembersRequest{Username: username})
}

func (c *Client) SearchMembersWithOptions(req SearchMembersRequest) ([]Autocomplete, error) {
	type results struct {
		Results []Autocomplete `json:"results" query:"results"`
	}
	response, err := PostDecode[results](c, ApiUrl("username_autosuggest"), structToUrlValues(req))
	return response.Results, err
}

func (u *User) SearchMembers(username string) ([]Autocomplete, error) {
	return u.Client().SearchMembers(username)
}

func (u *User) SearchMembersWithOptions(req SearchMembersRequest) ([]Autocomplete, error) {
	return u.Client().SearchMembersWithOptions(req)
}

func SearchMembers(username string) ([]Autocomplete, error) {
	return DefaultClient.SearchMembers(username)
}

func SearchMembersWithOptions(req SearchMembersRequest) ([]Autocomplete, error) {
	return DefaultClient.SearchMembersWithOptions(req)
}

// GetWatching gets the watchlist of a logged-in user
func (u *User) GetWatching() ([]UsernameID, error) {
	response, err := u.GetWatchingResponse(WatchlistRequest{})
	return response.Watches, err
}

func (u *User) GetWatchingResponse(req WatchlistRequest) (WatchlistResponse, error) {
	if req.SID == "" {
		if u.SID == "" {
			return WatchlistResponse{}, ErrNotLoggedIn
		}
		req.SID = u.SID
	}
	return PostDecode[WatchlistResponse](u.Client(), ApiUrl("watchlist"), structToUrlValues(req))
}
