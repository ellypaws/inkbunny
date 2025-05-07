package inkbunny

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ellypaws/inkbunny/types"
	"github.com/ellypaws/inkbunny/utils"
)

type User struct {
	client   *Client
	SID      string          `json:"sid" query:"sid"`
	Username string          `json:"username,omitempty" query:"username"`
	UserID   types.IntString `json:"user_id,omitempty" query:"user_id"`
	Ratings  types.Ratings   `json:"ratingsmask,omitempty" query:"ratingsmask"`
}

func (u *User) Client() *Client {
	if u.client != nil {
		return u.client
	}
	return DefaultClient
}

var (
	ErrNilUser       = errors.New("user is nil")
	ErrEmptyPassword = errors.New("username is set but password is empty")
	ErrNotLoggedIn   = errors.New("not logged in")
	ErrEmptySID      = errors.New("SID is empty")
	ErrUnexpectedSID = errors.New("unexpected session ID from response")
)

func Login(username, password string) (*User, error) {
	return DefaultClient.Login(username, password)
}

func (c *Client) Login(username, password string) (*User, error) {
	if username == "" {
		return nil, ErrNilUser
	}
	if username != "guest" && password == "" {
		return nil, ErrEmptyPassword
	}
	response, err := PostDecode[User](c, ApiUrl("login"), url.Values{"username": {username}, "password": {password}})
	if err != nil {
		return nil, fmt.Errorf("error logging in: %w", err)
	}
	if response.SID == "" {
		return nil, ErrUnexpectedSID
	}
	response.Username = username
	response.client = c
	return &response, nil
}

func (u *User) Logout() error {
	if u == nil {
		return ErrNilUser
	}
	if u.SID == "" {
		return ErrEmptySID
	}
	response, err := PostDecode[types.LogoutResponse](u.Client(), ApiUrl("logout"), url.Values{"sid": {u.SID}})
	if err != nil {
		return fmt.Errorf("error logging out: %w", err)
	}
	if response.Logout != "success" {
		return fmt.Errorf("logout failed, unexpected response: %s", response.Logout)
	}
	if response.SID != response.SID {
		return ErrUnexpectedSID
	}
	return nil
}

// ChangeRatings allows guest users to change their rating settings
//   - If you use this script to change rating settings for a logged-in registered member,
//     it will affect the current session only.
//     The changes to their allowed ratings will not be saved to their account.
//   - Members can still choose to block their work from Guest users, regardless of the Guests' rating choice, so some work may still not appear for Guests even with all rating options turned on.
//   - New Guest sessions and newly created accounts have the tag “Violence - Mild violence” enabled by default, so images tagged with this will be visible.
//     However, when calling this script, that tag will be set to “off”
//     unless you explicitly keep it activated with the parameter Ratings{MildViolence: true}.
//
// You can also call types.ParseMaskU if you want to use a bitmask.
func (u *User) ChangeRatings(ratings types.Ratings) error {
	if u == nil {
		return ErrNilUser
	}
	if u.SID == "" {
		return ErrEmptySID
	}
	values := utils.StructToUrlValues(ratings)
	values.Set("sid", u.SID)
	response, err := PostDecode[User](u.Client(), ApiUrl("userrating"), values)
	if err != nil {
		return fmt.Errorf("error changing ratings: %w", err)
	}
	if response.SID != u.SID {
		return ErrUnexpectedSID
	}
	u.Ratings = ratings
	return nil
}
