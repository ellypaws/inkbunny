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

func main() {
	u, err := login()
	if err != nil {
		log.Fatalf("error logging in: %v", err)
	}
	u.store()

	log.Printf("logged in as %s, session id: %s", u.Credentials.Username, u.Credentials.Sid)

	watchlistStrings := u.getWatchlist()

	log.Printf("watchers for %s: %v\n", u.Credentials.Username, watchlistStrings)

	u.Watchers = updateFollowers(u, watchlistStrings)
	u.store()

	var logout string
	fmt.Print("Logout? (y/[n]): ")
	fmt.Scanln(&logout)
	if logout == "y" {
		u.logout()
	}
}

func login() (user, error) {
	var u user = read()

	var err error
	if u.Credentials == nil || u.Credentials.Sid == "" {
		u.Credentials, err = loginPrompt().Login()
		if err != nil {
			return u, err
		}
	}
	return u, err
}

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

func (u user) logout() {
	err := u.Credentials.Logout()
	if err != nil {
		log.Fatalf("error logging out: %v", err)
	}
	u.store()
}

func updateFollowers(u user, watchlistStrings []string) map[string][]state {
	currentTime := time.Now().UTC()
	newWatchlist, newFollows := u.newFollowers(watchlistStrings, currentTime)
	if len(newFollows) > 0 {
		log.Printf("new follows (%d): %v", len(newFollows), newFollows)
	}

	if len(u.Watchers) == 0 {
		// First run, no previous states
		return newWatchlist
	}

	unfollows := u.checkMissing(newWatchlist, currentTime)
	if len(unfollows) > 0 {
		log.Printf("unfollows (%d): %v", len(unfollows), unfollows)
	}

	return newWatchlist
}

// newFollowers checks for new followers by checking if the watchlistStrings doesn't exist in user.Watchers or if the last state is unfollowing.
// If a username is new or previously unfollowed, it's added to the new watchlist with the last state set to following.
// Otherwise, the existing states are copied to the new watchlist.
func (u user) newFollowers(watchlistStrings []string, currentTime time.Time) (map[string][]state, []string) {
	newWatchlist := make(map[string][]state)
	var newFollows []string
	// Process new follows
	for _, watcher := range watchlistStrings {
		states, exists := u.Watchers[watcher]
		if !exists || !states[len(states)-1].Watching {
			newWatchlist[watcher] = append(states, state{Date: currentTime, Watching: true})
			newFollows = append(newFollows, watcher)
		} else {
			newWatchlist[watcher] = states // Copy existing states if no change
		}
	}
	return newWatchlist, newFollows
}

// checkMissing checks for unfollowing by looking at the missing usernames the new watchlist doesn't have that user.Watchers has.
// If a username is missing from the new watchlist, it's added to the new watchlist with the last state set to unfollowing.
func (u user) checkMissing(newWatchlist map[string][]state, currentTime time.Time) []string {
	var unfollows []string
	for username, states := range u.Watchers {
		if _, found := newWatchlist[username]; !found && states[len(states)-1].Watching {
			// Last state was following, now it's not in the new watchers list -> unfollow
			newWatchlist[username] = append(states, state{Date: currentTime, Watching: false})
			unfollows = append(unfollows, username)
		}
	}
	return unfollows
}

func (u user) getWatchlist() []string {
	var username string
	var watchlistStrings []string
	fmt.Print("Enter username to get watchers for: ")
	fmt.Scanln(&username)
	if username == "" || username == u.Credentials.Username || username == "self" {
		watching, err := u.Credentials.GetWatching()
		if err != nil {
			u.logout()
		}
		for i := range watching {
			watchlistStrings = append(watchlistStrings, watching[i].Username)
		}
	} else {
		watchedBy, err := u.Credentials.GetWatchedBy(username)
		if err != nil {
			u.logout()
		}
		for i := range watchedBy {
			watchlistStrings = append(watchlistStrings, watchedBy[i].Username)
		}
	}

	return watchlistStrings
}

type user struct {
	Credentials *api.Credentials   `json:"credentials"`
	Watchers    map[string][]state `json:"watchers"`
}

type state struct {
	Date     time.Time `json:"date"`
	Watching bool      `json:"following"`
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

func (u user) store() {
	tempFile, err := os.CreateTemp(".", "temp")
	if err != nil {
		log.Fatalf("error creating temp file: %v", err)
	}
	// It's important to close the file after creating it to ensure no locks are held on it.
	defer tempFile.Close()

	byte, err := json.Marshal(u)
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
