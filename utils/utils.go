package utils

import (
	"github.com/charmbracelet/bubbletea"
	"log"
	"net/url"
	"reflect"
	"slices"
	"time"
)

// Wrap casts a message into a tea.Cmd
func Wrap(msg any) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

// StructToUrlValues uses reflect to read json struct fields and set them as url.Values
// It also checks if omitempty is set and ignores empty fields
// Example:
//
//	type Example struct {
//		Field1 string `json:"field1,omitempty"`
//		Field2 string `json:"field2"`
//	}
func StructToUrlValues(s any) url.Values {
	var urlValues url.Values
	v := reflect.ValueOf(s)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if field.Tag.Get("json") == "" {
			continue
		}
		if field.Tag.Get("json") == "omitempty" && value.String() == "" {
			continue
		}

		// if boolean, use "yes" or "no"
		switch value.Kind() {
		case reflect.Bool:
			if value.Bool() {
				urlValues.Add(field.Tag.Get("json"), "yes")
			} else {
				urlValues.Add(field.Tag.Get("json"), "no")
			}
		default:
			urlValues.Add(field.Tag.Get("json"), value.String())
		}

	}
	return urlValues
}

type WatchInfo struct {
	Username  string
	Date      time.Time
	Watching  bool // Watching is true if you are watching the Username
	WatchedBy bool // WatchedBy is true if the Username is watching you
}

func UpdateNewMissing(oldList, newList []string) map[string][]WatchInfo {
	currentTime := time.Now().UTC()
	newWatchlist, newFollows := newFollowers(oldList, newList, currentTime)
	if len(newFollows) > 0 {
		log.Printf("new follows (%d): %v", len(newFollows), newFollows)
	}

	if len(newFollows) == 0 {
		// First run, no previous states
		return newWatchlist
	}

	unfollows := checkMissing(oldList, newWatchlist, currentTime)
	if len(unfollows) > 0 {
		log.Printf("unfollows (%d): %v", len(unfollows), unfollows)
	}

	return newWatchlist
}

// newFollowers checks for new followers by checking if the watchlistStrings doesn't exist in user.Watchers or if the last state is unfollowing.
// If a username is new or previously unfollowed, it's added to the new watchlist with the last state set to following.
// Otherwise, the existing states are copied to the new watchlist.
func newFollowers(oldList, newList []string, currentTime time.Time) (map[string][]WatchInfo, []string) {
	newWatchlist := make(map[string][]WatchInfo)
	var newFollows []string
	for _, username := range newList {
		// check if the username is in the old list
		if !slices.Contains(oldList, username) {
			// username is new
			newWatchlist[username] = append(newWatchlist[username], WatchInfo{Date: currentTime, Watching: true})
			newFollows = append(newFollows, username)
		} else {
			// username is not new
			newWatchlist[username] = append(newWatchlist[username], WatchInfo{Date: currentTime, Watching: false})
		}
	}
	return newWatchlist, newFollows
}

// checkMissing checks for unfollowing by looking at the missing usernames the new watchlist doesn't have that user.Watchers has.
// If a username is missing from the new watchlist, it's added to the new watchlist with the last state set to unfollowing.
func checkMissing(oldList []string, newWatchlist map[string][]WatchInfo, currentTime time.Time) []string {
	var unfollows []string
	for _, username := range oldList {
		states, exists := newWatchlist[username]
		if !exists || !states[len(states)-1].Watching {
			// Last state was user was watching you; now it's not in the new watchers list -> unfollow
			newWatchlist[username] = append(states, WatchInfo{Date: currentTime, Watching: false})
			unfollows = append(unfollows, username)
		}
	}
	return unfollows
}
