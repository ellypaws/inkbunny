package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

type WatchInfo struct {
	Username  string
	Date      time.Time
	Watching  bool // Watching is true if you are watching the Username
	WatchedBy bool // WatchedBy is true if the Username is watching you
}

// GetWatchedBy gets the list of users watching a given username
// It returns a slice of WatchInfo structs
func (user Credentials) GetWatchedBy(username string) ([]WatchInfo, error) {
	// Assuming GetSingleUser correctly retrieves a user ID for a given username
	userID, err := GetSingleUser(username)
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

	var watchlist []WatchInfo
	for {
		// Build the URL with query parameters for the current page
		apiUrl := inkbunnyURL("watchlist_process.php", url.Values{
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
			watchlist = append(watchlist, WatchInfo{
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

// GetWatching gets the watchlist of a logged-in user
func (user Credentials) GetWatching() ([]UsernameID, error) {
	resp, err := user.Get(apiURL("watchlist", url.Values{"sid": {user.Sid}}))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := CheckError(body); err != nil {
		return nil, fmt.Errorf("error getting watchlist: %w", err)
	}

	var watchResp WatchlistResponse
	if err := json.Unmarshal(body, &watchResp); err != nil {
		return nil, err
	}

	return watchResp.Watches, nil
}
