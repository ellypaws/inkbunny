package api

import (
	"github.com/ellypaws/inkbunny/api/utils"
	"testing"
)

func TestGuest(t *testing.T) {
	guest := Guest()
	t.Logf("Logging in as guest: %v", ApiUrl("login", utils.StructToUrlValues(guest)))
	user, err := guest.Login()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if user.Username != "guest" {
		t.Errorf("Expected username to be guest, got %s", user.Username)
	}

	if user.Sid == "" {
		t.Errorf("Expected sid to be non-empty, got empty")
	}

	if user.Ratings.String() != "101" {
		t.Errorf("Expected 101, got %s", user.Ratings.String())
	}
}

func TestCredentials_Login(t *testing.T) {
	user := &Credentials{Username: "guest"}
	t.Logf("Logging in as guest: %v", ApiUrl("login", utils.StructToUrlValues(user)))
	user, err := user.Login()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if user.Username != "guest" {
		t.Errorf("Expected username to be guest, got %s", user.Username)
	}

	if user.Sid == "" {
		t.Errorf("Expected sid to be non-empty, got empty")
	}

	if user.Ratings.String() != "101" {
		t.Errorf("Expected 101, got %s", user.Ratings.String())
	}
}

func TestCredentials_Logout(t *testing.T) {
	guest := Guest()
	t.Logf("Logging in as guest: %v", ApiUrl("login", utils.StructToUrlValues(guest)))
	user, err := guest.Login()
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}
	t.Logf("Logged in as %s, sid: %s\n", user.Username, user.Sid)

	t.Logf("Logging out: %v", ApiUrl("logout", utils.StructToUrlValues(user)))
	err = user.Logout()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if user.Sid != "" {
		t.Errorf("Expected sid to be empty, got %s", user.Sid)
	}
}
