package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ellypaws/inkbunny/api/utils"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

// InkbunnyUrl is a helper function to generate Inkbunny URLs with a given path and optional query parameters
func InkbunnyUrl(path string, values ...url.Values) *url.URL {
	request := &url.URL{
		Scheme: "https",
		Host:   "inkbunny.net",
		Path:   path,
	}

	var valueStrings []string
	for _, value := range values {
		valueStrings = append(valueStrings, value.Encode())
	}
	request.RawQuery = strings.Join(valueStrings, "&")

	return request
}

// ApiUrl is a helper function to generate Inkbunny API URLs.
// path is the name of the API endpoint, without the "api_" prefix or ".php" suffix
// example: "login" for "https://inkbunny.net/api_login.php"
//
//	url := ApiUrl("login", url.Values{"username": {"guest"}, "password": {""}})
func ApiUrl(path string, values ...url.Values) *url.URL {
	return InkbunnyUrl(fmt.Sprintf("api_%v.php", path), values...)
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

// ChangeRating allows guest users to change their rating settings
//   - If you use this script to change rating settings for a logged in registered member,
//     it will affect the current session only.
//     The changes to their allowed ratings will not be saved to their account.
//   - Members can still choose to block their work from Guest users, regardless of the Guests' rating choice, so some work may still not appear for Guests even with all rating options turned on.
//   - New Guest sessions and newly created accounts have the tag “Violence - Mild violence” enabled by default, so images tagged with this will be visible.
//     However, when calling this script, that tag will be set to “off”
//     unless you explicitly keep it activated with the parameter Ratings{MildViolence: true}.
func (user *Credentials) ChangeRating(ratings Ratings) error {
	values := utils.StructToUrlValues(ratings)
	values.Set("sid", user.Sid)
	resp, err := user.PostForm(ApiUrl("userrating"), values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var loginResp Credentials
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return err
	}

	if loginResp.Sid != user.Sid {
		return fmt.Errorf("session ID changed after rating change, expected: [%+v], got: [%+v]", user.Sid, loginResp.Sid)
	}

	user.Ratings = ratings

	return nil
}

func GetUserID(username string) (UsernameAutocomplete, error) {
	resp, err := Credentials{}.Get(ApiUrl("username_autosuggest", url.Values{"username": {username}}))
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

// GetFirstUser gets a single user by username, returns an error if no user is found
func GetFirstUser(username string) (Autocomplete, error) {
	users, err := GetUserID(username)
	if err != nil {
		return Autocomplete{}, err
	}
	if len(users.Results) == 0 {
		return Autocomplete{}, errors.New("user not found")
	}
	return users.Results[0], nil
}
