package inkbunny

import (
	"github.com/ellypaws/inkbunny/types"
	"github.com/ellypaws/inkbunny/utils"
)

type KeywordAutocomplete struct {
	ID               types.IntString `json:"id"`
	Value            string          `json:"value"`      // The user input with the partial keyword replaced by the suggested keyword.
	Icon             string          `json:"icon"`       // Deprecated: Not applicable to this usage (empty string).
	Info             string          `json:"info"`       // Deprecated: Not applicable to this usage (empty string).
	Keyword          string          `json:"singleword"` // The single keyword being suggested.
	SearchTerm       string          `json:"searchterm"` // The keyword identified in the user input being used to generate this suggestion.
	SubmissionsCount types.IntString `json:"submissions_count"`
}

// KeywordSuggestion suggests keywords based on partial keyword names entered by the user. It searches the start of keywords for the matching strings. It returns multiple matching suggestions as well as a count of the number of submissions each suggestion would find.
// The HTML response header will contain a directive for your client to cache the result data for 1 day, if it supports caching.
//   - All results are returned with HTML entities encoded. Eg: & will appear as &amp;, > will appear as &gt;, etc.
//   - All suggestions are returned ordered first by how many submissions they are assigned to (most commonly used first) and then alphabetically when that number is the same.
//   - Entering "hu" will search for keywords starting with "hu" and return an array of suggestions like "husky, human, hug, huge, humor, hunter" etc. Eg: https://inkbunny.net/api_search_autosuggest.php?keyword=hu&ratingsmask=11111
//   - Entering multiple words will search for words matching the last word, then matching the last two words, and so on. So entering "my li" will first search for all matches for "li" (lion, little, lizard...) and then all matches for "my li" (my little pony). This is similar to the keyword suggestion method Google uses and is most suitable for search boxes where users are likely to be entering multiple words. Eg: https://inkbunny.net/api_search_autosuggest.php?keyword=my+li&ratingsmask=11111&output_mode=xml
//   - To force the search to treat a set of words as one keyword only, they must be joined with underscores and the "underscorespaces" parameter must be set to yes. Eg: "my_li" will then only return a result for "my little pony". https://inkbunny.net/api_search_autosuggest.php?keyword=my_li&underscorespaces=yes&ratingsmask=11111
//   - Keywords are filtered roughly by user ratings settings. So keywords that generally appear on mature/adult rated images will be hidden from users who have those higher ratings turned off. You must send the user ratings selection (see "ratingsmask" parameter below) or the default "G rated only" will be used. Eg: https://inkbunny.net/api_search_autosuggest.php?keyword=hu&ratingsmask=11100&output_mode=xml
func (c *Client) KeywordSuggestion(keyword string, ratings types.Ratings, underscore bool) ([]KeywordAutocomplete, error) {
	type params struct {
		Keyword          string          `json:"keyword"`
		Ratings          types.Ratings   `json:"ratingsmask"`
		UnderscoreSpaces types.BooleanYN `json:"underscorespaces"`
	}
	param := params{
		Keyword:          keyword,
		Ratings:          ratings,
		UnderscoreSpaces: types.BooleanYN(underscore),
	}
	type results struct {
		Results []KeywordAutocomplete `json:"results"`
	}
	response, err := PostDecode[results](c, ApiUrl("search_autosuggest"), utils.StructToUrlValues(param))
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
func KeywordSuggestion(keyword string, ratings types.Ratings, underscore bool) ([]KeywordAutocomplete, error) {
	return DefaultClient.KeywordSuggestion(keyword, ratings, underscore)
}
