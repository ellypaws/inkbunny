package api

import (
	"encoding/json"
	"fmt"
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

type Rating struct {
	General        bool `json:"1,omitempty"` // Show images with Rating tag: "General - Suitable for all ages".
	Nudity         bool `json:"2,omitempty"` // Show images with Rating tag: "Nudity - Nonsexual nudity exposing breasts or genitals (must not show arousal)".
	MildViolence   bool `json:"3,omitempty"` // Show images with Rating tag: "Violence - Mild violence".
	Sexual         bool `json:"4,omitempty"` // Show images with Rating tag: "Sexual Themes - Erotic imagery, sexual activity or arousal".
	StrongViolence bool `json:"5,omitempty"` // Show images with Rating tag: "Strong Violence - Strong violence, blood, serious injury or death".
}

func optIn(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func (user Credentials) changeRating(rating Rating) error {
	resp, err := user.PostForm(inkbunnyURL("userrating"), url.Values{
		"sid":    {user.Sid},
		"tag[2]": {optIn(rating.Nudity)},
		"tag[3]": {optIn(rating.MildViolence)},
		"tag[4]": {optIn(rating.Sexual)},
		"tag[5]": {optIn(rating.StrongViolence)},
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

	if loginResp.Sid != user.Sid {
		return fmt.Errorf("session ID changed after rating change, expected: [%s], got: [%s]", user.Sid, loginResp.Sid)
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
