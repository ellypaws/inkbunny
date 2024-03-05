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
	resp, err := user.PostForm(ApiUrl("login"), url.Values{"username": {user.Username}, "password": {user.Password}})
	user.Password = ""
	if err != nil {
		return nil, fmt.Errorf("error logging in: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var respLog struct {
		Credentials
		RatingsMask string `json:"ratingsmask"`
	}
	if err = json.Unmarshal(body, &respLog); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	if respLog.Sid == "" {
		return nil, fmt.Errorf("sid is empty, response: %s", body)
	}

	user.Sid = respLog.Sid
	user.Ratings = parseMask(respLog.RatingsMask)

	return user, nil
}

func (user Credentials) LoggedIn() bool {
	return user.Sid != ""
}

func (user *Credentials) Logout() error {
	if user == nil {
		return ErrNotLoggedIn
	}
	resp, err := user.Get(ApiUrl("logout", url.Values{"sid": {user.Sid}}))
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

	if logoutResp.Sid != user.Sid {
		return fmt.Errorf("session ID changed after logout, expected: [%s], got: [%s]", user.Sid, logoutResp.Sid)
	}

	*user = Credentials{}
	return nil
}
