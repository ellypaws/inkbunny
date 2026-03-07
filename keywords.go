package inkbunny

// KeywordAutocomplete is one keyword suggestion returned by KeywordSuggestion.
type KeywordAutocomplete struct {
	// ID is the keyword ID.
	ID         IntString `json:"id"`
	Value      string    `json:"value"`      // Value is the completed input string returned by Inkbunny.
	Icon       string    `json:"icon"`       // Icon is unused by this endpoint and is usually empty.
	Info       string    `json:"info"`       // Info is unused by this endpoint and is usually empty.
	Keyword    string    `json:"singleword"` // Keyword is the suggested keyword token.
	SearchTerm string    `json:"searchterm"` // SearchTerm is the fragment matched from the caller input.
	// SubmissionsCount is the number of submissions assigned to the keyword.
	SubmissionsCount IntString `json:"submissions_count"`
}

// KeywordSuggestion returns keyword autocomplete results for the active user or client.
//
// Pass Ratings to control content filtering, typically using ParseMaskU or a
// Ratings literal. Set underscore to true when the caller is working with
// underscore-joined multiword keywords.
func (u *User) KeywordSuggestion(keyword string, ratings Ratings, underscore bool) ([]KeywordAutocomplete, error) {
	return u.Client().KeywordSuggestion(keyword, ratings, underscore)
}

// KeywordSuggestion returns keyword autocomplete results for a client.
//
// Pass Ratings to control content filtering, typically using ParseMaskU or a
// Ratings literal. Set underscore to true when the caller is working with
// underscore-joined multiword keywords.
func (c *Client) KeywordSuggestion(keyword string, ratings Ratings, underscore bool) ([]KeywordAutocomplete, error) {
	type params struct {
		Keyword          string    `json:"keyword"`
		Ratings          Ratings   `json:"ratingsmask"`
		UnderscoreSpaces BooleanYN `json:"underscorespaces"`
	}
	param := params{
		Keyword:          keyword,
		Ratings:          ratings,
		UnderscoreSpaces: BooleanYN(underscore),
	}
	type results struct {
		Results []KeywordAutocomplete `json:"results"`
	}
	response, err := PostDecode[results](c, ApiUrl("search_autosuggest"), structToUrlValues(param))
	return response.Results, err
}

// KeywordSuggestion suggests keywords based on partial keyword names entered by the user. It searches the start of keywords for the matching strings. It returns multiple matching suggestions as well as a count of the number of submissions each suggestion would find.
// The HTML response header will contain a directive for your client to cache the result data for 1 day, if it supports caching.
//   - All results are returned with HTML entities encoded. Eg: & will appear as &amp;, > will appear as &gt;, etc.
//   - All suggestions are returned ordered first by how many submissions they are assigned to (most commonly used first) and then alphabetically when that number is the same.
//   - Entering "hu" will search for keywords starting with "hu" and return an array of suggestions like "husky, human, hug, huge, humor, hunter" etc. Eg: https://inkbunny.net/api_search_autosuggest.php?keyword=hu&ratingsmask=11111
//   - Entering multiple words will search for words matching the last word, then matching the last two words, and so on. So entering "my li" will first search for all matches for "li" (lion, little, lizard...) and then all matches for "my li" (my little pony). This is similar to the keyword suggestion method Google uses and is most suitable for search boxes where users are likely to be entering multiple words. Eg: https://inkbunny.net/api_search_autosuggest.php?keyword=my+li&ratingsmask=11111&output_mode=xml
//   - To force the search to treat a set of words as one keyword only, they must be joined with underscores and the "underscorespaces" parameter must be set to yes. Eg: "my_li" will then only return a result for "my little pony". https://inkbunny.net/api_search_autosuggest.php?keyword=my_li&underscorespaces=yes&ratingsmask=11111
//   - Keywords are filtered roughly by user ratings settings. So keywords that generally appear on mature/adult rated images will be hidden from users who have those higher ratings turned off. You must send the user ratings selection (see "ratingsmask" parameter below) or the default "G rated only" will be used. Eg: https://inkbunny.net/api_search_autosuggest.php?keyword=hu&ratingsmask=11100&output_mode=xml
func KeywordSuggestion(keyword string, ratings Ratings, underscore bool) ([]KeywordAutocomplete, error) {
	return DefaultClient.KeywordSuggestion(keyword, ratings, underscore)
}
