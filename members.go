package inkbunny

import (
	"net/url"

	"github.com/ellypaws/inkbunny/types"
)

func (c *Client) SearchMembers(username string) ([]types.Autocomplete, error) {
	type results struct {
		Results []types.Autocomplete `json:"results" query:"results"`
	}
	response, err := PostDecode[results](c, ApiUrl("username_autosuggest"), url.Values{"username": {username}})
	return response.Results, err
}

func (u *User) SearchMembers(username string) ([]types.Autocomplete, error) {
	return u.Client().SearchMembers(username)
}

func SearchMembers(username string) ([]types.Autocomplete, error) {
	return DefaultClient.SearchMembers(username)
}

// GetWatching gets the watchlist of a logged-in user
func (u *User) GetWatching() ([]types.UsernameID, error) {
	if u.SID == "" {
		return nil, ErrNotLoggedIn
	}
	type results struct {
		Watches []types.UsernameID `json:"watches"`
	}
	response, err := PostDecode[results](u.Client(), ApiUrl("watchlist"), url.Values{"sid": {u.SID}})
	return response.Watches, err
}
