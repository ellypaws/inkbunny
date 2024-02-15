package api

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"net/http"
	"net/url"
)

func Login(user *Credentials) tea.Cmd {
	if user.Username == "" {
		user.Username = "guest"
	} else if user.Password == "" {
		return Wrap(fmt.Errorf("username is set but password is empty"))
	}
	resp, err := http.PostForm(inkbunnyURL("login"), url.Values{"username": {user.Username}, "password": {user.Password}})
	if err != nil {
		return Wrap(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Wrap(err)
	}

	if err = json.Unmarshal(body, user); err != nil {
		return Wrap(err)
	}

	return Wrap(user)
}

func logout(sid string) error {
	resp, err := http.Get(inkbunnyURL("logout", url.Values{"sid": {sid}}))
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

	return nil
}
