package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
)

var (
	ErrNilUser       = errors.New("user is nil")
	ErrEmptyPassword = errors.New("username is set but password is empty")
	ErrNotLoggedIn   = errors.New("not logged in")
)

func (user *Credentials) Login() (*Credentials, error) {
	if user == nil {
		return nil, ErrNilUser
	}
	if user.Username == "" || user.Username == "guest" {
		user.Username = "guest"
	} else if user.Password == "" {
		return nil, ErrEmptyPassword
	}
	resp, err := user.PostForm(apiURL("login"), url.Values{"username": {user.Username}, "password": {user.Password}})
	user.Password = ""
	if err != nil {
		return nil, fmt.Errorf("error logging in: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if err = json.Unmarshal(body, user); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	user.Ratings.parseMask()

	return user, nil
}

func (user Credentials) LoggedIn() bool {
	return user.Sid != ""
}

func (ratings *Ratings) parseMask() {
	// RatingsMask - Binary string representation of the users Allowed Ratings choice. The bits are in this order left-to-right:
	// Eg: A string 11100 means only items rated General, Nudity and Violence are allowed, but Sex and Strong Violence are blocked.
	// A string 11111 means items of any rating would be shown. Only 'left-most significant bits' are returned. So 11010 and 1101 are the same, and 10000 and 1 are the same.
	set := func(r int32) bool {
		return r == '1'
	}
	for i, rating := range ratings.RatingsMask {
		switch i {
		case 0:
			ratings.General = set(rating)
		case 1:
			ratings.Nudity = set(rating)
		case 2:
			ratings.MildViolence = set(rating)
		case 3:
			ratings.Sexual = set(rating)
		case 4:
			ratings.StrongViolence = set(rating)
		}
	}
}

func (user *Credentials) Logout() error {
	if user == nil {
		return ErrNotLoggedIn
	}
	resp, err := user.Get(apiURL("logout", url.Values{"sid": {user.Sid}}))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var logoutResp LogoutResponse
	if err = json.Unmarshal(body, &logoutResp); err != nil {
		return err
	}

	if logoutResp.Logout != "success" {
		return fmt.Errorf("logout failed, response: %s", logoutResp.Logout)
	}

	*user = Credentials{}
	return nil
}
