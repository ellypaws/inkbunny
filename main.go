package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
)

const (
	baseURL = "https://inkbunny.net/api_"
)

type LoginResponse struct {
	Sid string `json:"sid"`
}

type WatchlistResponse struct {
	Watches []struct {
		Username string `json:"username"`
	} `json:"watches"`
}

// Function to login and get session ID
func login(username, password string) (string, error) {
	resp, err := http.PostForm(baseURL+"login.php", url.Values{"username": {username}, "password": {password}})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return "", err
	}

	return loginResp.Sid, nil
}

// Function to get watchlist for a given user
func getWatchlist(sid string) ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("%swatchlist.php?sid=%s&user_id=%s", baseURL, sid))
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

func changeRating(sid string) error {
	resp, err := http.PostForm(baseURL+"userrating.php", url.Values{
		"sid":    {sid},
		"tag[2]": {"yes"},
		"tag[3]": {"yes"},
		"tag[4]": {"yes"},
		"tag[5]": {"yes"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return err
	}

	if loginResp.Sid != sid {
		return fmt.Errorf("session ID changed after rating change, expected: [%s], got: [%s]", sid, loginResp.Sid)
	}

	return nil
}

type LogoutResponse struct {
	Sid    string `json:"sid"`
	Logout string `json:"logout"`
}

func logout(sid string) error {
	resp, err := http.PostForm(baseURL+"logout.php", url.Values{"sid": {sid}})
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

type User struct {
	ID         string `json:"id"`
	Value      string `json:"value"`
	Icon       string `json:"icon"`
	Info       string `json:"info"`
	SingleWord string `json:"singleword"`
	SearchTerm string `json:"searchterm"`
}

func getUserID(username string) (User, error) {
	resp, err := http.Get(fmt.Sprintf("%susername_autosuggest.php?username=%s", baseURL, username))
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return User{}, err
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return User{}, err
	}

	return user, nil
}

func main() {
	// prompt for username and password
	var username, password string
	fmt.Print("Enter username: ")
	fmt.Scanln(&username)
	if username == "" {
		username = "guest"
	}
	fmt.Print("Enter password: ")
	fmt.Scanln(&password)
	sid, err := login(username, password)
	if err != nil {
		fmt.Println("Login error:", err)
		return
	}

	log.Printf("Logged in as [%s] with session ID [%v] ", username, sid)

	watchlist1, err := getWatchlist(sid)
	if err != nil {
		log.Println("Error getting watchlist for user 1:", err)
		return
	}

	fmt.Println("Watch list:", watchlist1)

	if err := logout(sid); err != nil {
		log.Println("Logout error:", err)
		return
	}
	log.Println("Logged out")
}
