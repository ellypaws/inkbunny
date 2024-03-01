package api

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

func (user Credentials) GetWatchers(username string) ([]string, error) {
	var watchers []string
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
		foundWatchers := false
		tableSelection.Find("tbody > tr").Each(func(i int, s *goquery.Selection) {
			log.Printf("Processing tr #%d", i+1)
			// Ignore rows that don't have a link (e.g., date rows)
			link := s.Find("td > span > a")
			href, exists := link.Attr("href")
			if exists {
				foundWatchers = true
				// Extract username from URL
				watcherUsername := strings.TrimPrefix(href, "/")
				watchers = append(watchers, watcherUsername)
				log.Printf("Found watcher: %s", watcherUsername)
			}
		})

		if !foundWatchers {
			log.Println("No watchers found on this page, assuming end of list.")
			break
		}

		time.Sleep(1 * time.Second)

		// Increment the page number for the next iteration
		page++
		retries = 0
	}

	return watchers, nil
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
