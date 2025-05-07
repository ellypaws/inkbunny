package types

type LogoutResponse struct {
	SID    string `json:"sid"`
	Logout string `json:"logout"`
}

type UsernameID struct {
	UserID   string `json:"user_id" query:"user_id"`
	Username string `json:"username" query:"username"`
}

type Autocomplete struct {
	ID    IntString `json:"id"`    // User ID of the suggested user.
	Value string    `json:"value"` // The user input with the partial username replaced by the suggested username.
	// Path and file name of user icon (if account has a user icon set). Note that
	// this is a relative path like "27/27014_fred.jpg". You need to prepend the
	// full location to the start of this string get the icon size you want. Eg:
	//  - "/usericons/tiny/27/27014_fred.jpg" for tiny icon 20x20px,
	//  - "/usericons/small/27/27014_fred.jpg" for small icon 50x50px,
	//  - "/usericons/large/27/27014_fred.jpg" for large icon 100x100px.
	Icon       string `json:"icon"`
	Info       string `json:"info"`       // Additional information about the selection (usually blank).
	SingleWord string `json:"singleword"` // The single username being suggested.
	SearchTerm string `json:"searchterm"` // They keyword identified in the user input being used to generate this suggestion.
}
