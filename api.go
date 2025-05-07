package inkbunny

import (
	"fmt"
	"net/url"
	"strings"
)

// InkbunnyUrl is a helper function to generate Inkbunny URLs with a given path and optional query parameters
func InkbunnyUrl(path string, values ...url.Values) *url.URL {
	request := &url.URL{
		Scheme: "https",
		Host:   "inkbunny.net",
		Path:   path,
	}

	var valueStrings []string
	for _, value := range values {
		valueStrings = append(valueStrings, value.Encode())
	}
	request.RawQuery = strings.Join(valueStrings, "&")

	return request
}

// ApiUrl is a helper function to generate Inkbunny API URLs.
// path is the name of the API endpoint, without the "api_" prefix or ".php" suffix
// example: "login" for "https://inkbunny.net/api_login.php"
//
//	url := ApiUrl("login", url.Values{"username": {"guest"}, "password": {""}})
func ApiUrl(path string, values ...url.Values) *url.URL {
	return InkbunnyUrl(fmt.Sprintf("api_%v.php", path), values...)
}
