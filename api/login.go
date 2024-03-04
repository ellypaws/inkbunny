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

func Guest() *Credentials {
	return &Credentials{Username: "guest"}
}

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
