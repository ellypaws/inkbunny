package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

func (user *Credentials) Login() (error, *Credentials) {
	if user.Username == "" {
		user.Username = "guest"
	} else if user.Password == "" {
		return fmt.Errorf("username is set but password is empty"), nil
	}
	resp, err := user.PostForm(inkbunnyURL("login"), url.Values{"username": {user.Username}, "password": {user.Password}})
	user.Password = ""
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}

	if err = json.Unmarshal(body, user); err != nil {
		return err, nil
	}

	return nil, user
}

func (user *Credentials) Logout() error {
	resp, err := user.Get(inkbunnyURL("logout", url.Values{"sid": {user.Sid}}))
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
