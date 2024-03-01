package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/term"
	"inkbunny/api"
	"log"
	"os"
	"time"
)

func loginPrompt() *api.Credentials {
	var user api.Credentials
	fmt.Print("Enter username: ")
	fmt.Scanln(&user.Username)
	fmt.Print("Enter password: ")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(bytePassword)

	return &user
}

func main() {
	var u user = read()

	var err error
	if u.Credentials == nil || u.Credentials.Sid == "" {
		u.Credentials, err = loginPrompt().Login()
		if err != nil {
			log.Fatalf("error logging in: %v", err)
		}
		u.store()
	}

	log.Printf("logged in as %s, session id: %s", u.Credentials.Username, u.Credentials.Sid)

	var username string
	var watchlistStrings []string
	fmt.Print("Enter username to get watchers for: ")
	fmt.Scanln(&username)
	if username == "" || username == u.Credentials.Username || username == "self" {
		watchlistStrings, err = u.Credentials.GetWatchlist()
	} else {
		watchlistStrings, err = u.Credentials.GetWatchers(username)
	}
	if err != nil {
		err := u.Credentials.Logout()
		if err != nil {
			log.Fatalf("error logging out: %v", err)
		}
		u.store()
		log.Fatalf("error getting watchers: %v", err)
	}

	log.Printf("watchers for %s: %v\n", username, watchlistStrings)

	currentTime := time.Now().UTC()
	newWatchersList := make(map[string][]state)

	var newFollows []string
	// Process new follows
	for _, watcher := range watchlistStrings {
		states, exists := u.Watchers[watcher]
		if !exists || !states[len(states)-1].Following {
			newWatchersList[watcher] = append(states, state{Date: currentTime, Following: true})
			newFollows = append(newFollows, watcher)
		} else {
			newWatchersList[watcher] = states // Copy existing states if no change
		}
	}

	if len(newFollows) > 0 {
		log.Printf("new follows: %v", newFollows)
	}

	// Check for unfollows
	var unfollows []string
	for watcher, states := range u.Watchers {
		if _, found := newWatchersList[watcher]; !found && states[len(states)-1].Following {
			// Last state was following, now it's not in the new watchers list -> unfollow
			newWatchersList[watcher] = append(states, state{Date: currentTime, Following: false})
			unfollows = append(unfollows, watcher)
		} else if !states[len(states)-1].Following {
			// If the last state was not following, just copy over the states
			newWatchersList[watcher] = states
		}
	}

	if len(unfollows) > 0 {
		log.Printf("unfollows: %v", unfollows)
	}

	u.Watchers = newWatchersList
	u.store()

	var logout string
	fmt.Print("Logout? (y/[n]): ")
	fmt.Scanln(&logout)
	if logout == "y" {
		err = u.Credentials.Logout()
		if err != nil {
			log.Fatalf("error logging out: %v", err)
		}
		u.store()
	}
}

type user struct {
	Credentials *api.Credentials   `json:"credentials"`
	Watchers    map[string][]state `json:"watchers"`
}

type state struct {
	Date      time.Time `json:"date"`
	Following bool      `json:"following"`
}

func read() user {
	if _, err := os.Stat("user.json"); err != nil {
		return user{}
	}
	byte, err := os.ReadFile("user.json")
	if err != nil {
		log.Fatalf("error reading user file: %v", err)
	}
	var u user
	err = json.Unmarshal(byte, &u)
	if err != nil {
		log.Fatalf("error unmarshalling user file: %v", err)
	}
	return u
}

func (user user) store() {
	tempFile, err := os.CreateTemp(".", "temp")
	if err != nil {
		log.Fatalf("error creating temp file: %v", err)
	}
	// It's important to close the file after creating it to ensure no locks are held on it.
	defer tempFile.Close()

	byte, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("error marshalling user: %v", err)
	}

	_, err = tempFile.Write(byte)
	if err != nil {
		log.Fatalf("error writing to temp file: %v", err)
	}

	// Sync and then close the file explicitly before attempting to rename it.
	// The defer statement will still attempt to close the file, but it's already closed at this point.
	// This redundancy in closing is fine and ensures the file is definitely closed before renaming.
	err = tempFile.Sync()
	if err != nil {
		log.Fatalf("error syncing temp file: %v", err)
	}

	// Close the file before renaming to release any locks.
	tempFile.Close()

	err = os.Rename(tempFile.Name(), "user.json")
	if err != nil {
		log.Fatalf("error renaming temp file: %v", err)
	}

	log.Printf("successfully updated user.json")
}
