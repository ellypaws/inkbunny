package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sahilm/fuzzy"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

// inkbunnyURL is a helper function to generate Inkbunny URLs with a given path and optional query parameters
func inkbunnyURL(path string, values ...url.Values) *url.URL {
	request := &url.URL{
		Scheme: "https",
		Host:   "inkbunny.net",
		Path:   path,
	}
	for i := range values {
		request.RawQuery = values[i].Encode()
	}

	return request
}

// apiURL is a helper function to generate Inkbunny API URLs.
// path is the name of the API endpoint, without the "api_" prefix or ".php" suffix
// example: "login" for "https://inkbunny.net/api_login.php"
//
//	url := apiURL("login", url.Values{"username": {"guest"}, "password": {""}})
func apiURL(path string, values ...url.Values) *url.URL {
	return inkbunnyURL(fmt.Sprintf("api_%v.php", path), values...)
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

func optIn(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// ChangeRating allows guest users to change their rating settings
//   - If you use this script to change rating settings for a logged in registered member,
//     it will affect the current session only.
//     The changes to their allowed ratings will not be saved to their account.
//   - Members can still choose to block their work from Guest users, regardless of the Guests' rating choice, so some work may still not appear for Guests even with all rating options turned on.
//   - New Guest sessions and newly created accounts have the tag “Violence - Mild violence” enabled by default, so images tagged with this will be visible.
//     However, when calling this script, that tag will be set to “off”
//     unless you explicitly keep it activated with the parameter Ratings{MildViolence: true}.
func (user Credentials) ChangeRating(rating Ratings) error {
	resp, err := user.PostForm(apiURL("userrating"), url.Values{
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

func GetUserID(username string) (UsernameAutocomplete, error) {
	resp, err := Credentials{}.Get(apiURL("username_autosuggest", url.Values{"username": {username}}))
	if err != nil {
		return UsernameAutocomplete{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UsernameAutocomplete{}, err
	}

	var users UsernameAutocomplete
	if err := json.Unmarshal(body, &users); err != nil {
		return UsernameAutocomplete{}, err
	}

	return users, nil
}

// GetSingleUser gets a single user by username, returns an error if no user is found
func GetSingleUser(username string) (Autocomplete, error) {
	users, err := GetUserID(username)
	if err != nil {
		return Autocomplete{}, err
	}
	if len(users.Results) == 0 {
		return Autocomplete{}, errors.New("user not found")
	}
	// sort by the closest match using fuzzy
	matches := fuzzy.FindFrom(username, users)
	if len(matches) == 0 {
		return Autocomplete{}, errors.New("user not found")
	}
	return users.Results[matches[0].Index], nil
}
