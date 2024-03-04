package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ellypaws/inkbunny/api"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/term"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
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
		watchedBy, err := GetWatchedBy(*u.Credentials, username)
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

// GetWatchedBy gets the list of users watching a given username
// It returns a slice of WatchInfo structs
func GetWatchedBy(user api.Credentials, username string) ([]api.WatchInfo, error) {
	// Assuming GetFirstUser correctly retrieves a user ID for a given username
	userID, err := api.GetFirstUser(username)
	if err != nil {
		return nil, fmt.Errorf("error getting user ID: %w", err)
	}

	// Create a cookie jar to manage session cookies
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	// Add the PHPSESSID cookie to the jar if needed
	//urlObj, _ := url.Parse("https://inkbunny.net")
	//phpSessionCookies := []*http.Cookie{
	//	{Name: "PHPSESSID", Value: user.Sid},
	//}
	//jar.SetCookies(urlObj, phpSessionCookies)

	client := &http.Client{
		Jar: jar,
	}

	// Initialize pagination
	page := 1
	backoff := 5 * time.Second
	var retries int

	var watchlist []api.WatchInfo
	for {
		// Build the URL with query parameters for the current page
		apiUrl := api.InkbunnyUrl("watchlist_process.php", url.Values{
			"page":      {strconv.Itoa(page)},
			"user_id":   {userID.ID},
			"orderby":   {"added"},
			"namesonly": {"yes"},
			"mode":      {"watchedby"},
			//"sid":       {user.Sid},
		}).String()

		log.Printf("Fetching page %d: %s", page, apiUrl)

		// Make the HTTP request using the custom client with cookie jar
		resp, err := client.Get(apiUrl)
		if err != nil {
			return nil, fmt.Errorf("error fetching page %d: %w", page, err)
		}
		defer resp.Body.Close()

		// Check the response status code, retry on 429
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := resp.Header.Get("Retry-After")
			if retryAfter != "" {
				waitDuration, err := parseRetryAfter(retryAfter)
				if err != nil {
					log.Printf("Failed to parse Retry-After header: %v, waiting for 5 seconds", err)
					time.Sleep(5 * time.Second)
					continue
				} else {
					log.Printf("Rate limited. Waiting for %v before retrying...", waitDuration)
					time.Sleep(waitDuration)
					continue
				}
			} else {
				log.Printf("Rate limited but no Retry-After header found, waiting for %v before retrying...", backoff)
				time.Sleep(backoff)
				if retries > 0 {
					backoff *= 2
				}
				retries++
				continue
			}
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code %d for page %d", resp.StatusCode, page)
		}

		// Parse the HTML response
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error parsing HTML for page %d: %w", page, err)
		}

		// Check if the table is found
		tableSelection := doc.Find("table.warthog")
		if tableSelection.Length() == 0 {
			log.Println("No table found on the page")
			break
		} else {
			log.Println("Table found")
		}

		// Select the table and iterate over each row
		var foundWatchers bool
		var numWatchers int
		tableSelection.Find("tbody > tr").Each(func(i int, s *goquery.Selection) {
			log.Printf("Processing tr #%d", i+1)
			// Ignore rows that don't have a link (e.g., date rows)
			link := s.Find("td > span > a")
			if link.Length() == 0 {
				return
			}
			foundWatchers = true
			watchlist = append(watchlist, api.WatchInfo{
				Username:  link.Text(),
				Watching:  link.HasClass("watching"),
				WatchedBy: true,
			})
			numWatchers++
		})

		if !foundWatchers {
			log.Println("No watchers found on this page, assuming end of list.")
			break
		}

		log.Printf("Found %d watchers on page %d", numWatchers, page)

		time.Sleep(1 * time.Second)

		// Increment the page number for the next iteration
		page++
		retries = 0
	}

	return watchlist, nil
}

// parseRetryAfter attempts to parse the Retry-After header value and return the corresponding duration.
// The Retry-After header can be in the form of a number (of seconds) or a date.
func parseRetryAfter(retryAfter string) (time.Duration, error) {
	// First, try to parse as a number of seconds.
	if seconds, err := strconv.Atoi(retryAfter); err == nil {
		return time.Duration(seconds) * time.Second, nil
	}

	// If not a number, try to parse as a date.
	if retryTime, err := http.ParseTime(retryAfter); err == nil {
		return time.Until(retryTime), nil
	}

	return 0, fmt.Errorf("invalid Retry-After format")
}
