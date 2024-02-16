package api

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

// inkbunnyURL is a helper function to generate Inkbunny API URLs
func inkbunnyURL(path string, values ...url.Values) *url.URL {
	request := &url.URL{
		Scheme: "https",
		Host:   "inkbunny.net",
		Path:   fmt.Sprintf("api_%v.php", path),
	}
	for i := range values {
		request.RawQuery = values[i].Encode()
	}

	return request
}

func (user Credentials) Request(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if user.Sid != "" {
		req.AddCookie(&http.Cookie{Name: "PHPSESSID", Value: user.Sid})
	}
	return req, nil
}

func (user Credentials) Get(url *url.URL) (*http.Response, error) {
	req, _ := user.Request("GET", url.String(), nil)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (user Credentials) Post(url *url.URL, contentType string, body io.Reader) (*http.Response, error) {
	req, _ := user.Request("POST", url.String(), body)
	req.Header.Set("Content-Type", contentType)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (user Credentials) PostForm(url *url.URL, values url.Values) (*http.Response, error) {
	return user.Post(url, "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
}

// Wrap casts a message into a tea.Cmd
func Wrap(msg any) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

// GetWatchlist gets the watchlist of a logged-in user
func GetWatchlist(user Credentials) ([]string, error) {
	resp, err := user.Get(inkbunnyURL("watchlist", url.Values{"sid": {user.Sid}}))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var watchResp WatchlistResponse
	if err := json.Unmarshal(body, &watchResp); err != nil {
		return nil, err
	}

	var usernames []string
	for _, watch := range watchResp.Watches {
		usernames = append(usernames, watch.Username)
	}

	return usernames, nil
}

// Function to find mutual elements in two slices
func findMutual(a, b []string) []string {
	var mutual []string
	for _, val := range a {
		if slices.Contains(b, val) {
			mutual = append(mutual, val)
		}
	}
	return mutual
}

func changeRating(sid string) error {
	resp, err := http.PostForm(inkbunnyURL("userrating"), url.Values{
		"sid":    {sid},
		"tag[2]": {"yes"},
		"tag[3]": {"yes"},
		"tag[4]": {"yes"},
		"tag[5]": {"yes"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var loginResp Credentials
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return err
	}

	if loginResp.Sid != sid {
		return fmt.Errorf("session ID changed after rating change, expected: [%s], got: [%s]", sid, loginResp.Sid)
	}

	return nil
}

func (user Credentials) getUserID(username string) ([]User, error) {
	resp, err := user.Get(inkbunnyURL("username_autosuggest", url.Values{"username": {username}}))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var users AutocompleteResponse
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, err
	}

	return users.Results, nil
}
