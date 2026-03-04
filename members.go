package inkbunny

import (
	"net/url"
)

func (c *Client) SearchMembers(username string) ([]Autocomplete, error) {
	type results struct {
		Results []Autocomplete `json:"results" query:"results"`
	}
	response, err := PostDecode[results](c, ApiUrl("username_autosuggest"), url.Values{"username": {username}})
	return response.Results, err
}

func (u *User) SearchMembers(username string) ([]Autocomplete, error) {
	return u.Client().SearchMembers(username)
}

func SearchMembers(username string) ([]Autocomplete, error) {
	return DefaultClient.SearchMembers(username)
}

// GetWatching gets the watchlist of a logged-in user
func (u *User) GetWatching() ([]UsernameID, error) {
	if u.SID == "" {
		return nil, ErrNotLoggedIn
	}
	type results struct {
		Watches []UsernameID `json:"watches"`
	}
	response, err := PostDecode[results](u.Client(), ApiUrl("watchlist"), url.Values{"sid": {u.SID}})
	return response.Watches, err
}
