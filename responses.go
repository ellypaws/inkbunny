package inkbunny

// LogoutResponse is returned by User.Logout.
type LogoutResponse struct {
	// SID is the session ID that was logged out.
	SID string `json:"sid"`
	// Logout is "success" when the logout completed.
	Logout string `json:"logout"`
}

// UsernameID is the package's lightweight user reference type.
type UsernameID struct {
	// UserID is the user ID.
	UserID string `json:"user_id" query:"user_id"`
	// Username is the current username.
	Username string `json:"username" query:"username"`
}

// Autocomplete is one username suggestion returned by SearchMembers.
type Autocomplete struct {
	ID    IntString `json:"id"`    // ID is the suggested user's ID.
	Value string    `json:"value"` // Value is the completed input string returned by Inkbunny.
	// Path and file name of user icon (if account has a user icon set). Note that
	// this is a relative path like "27/27014_fred.jpg". You need to prepend the
	// full location to the start of this string get the icon size you want. Eg:
	//  - "/usericons/tiny/27/27014_fred.jpg" for tiny icon 20x20px,
	//  - "/usericons/small/27/27014_fred.jpg" for small icon 50x50px,
	//  - "/usericons/large/27/27014_fred.jpg" for large icon 100x100px.
	Icon       string `json:"icon"`
	Info       string `json:"info"`       // Info is optional extra metadata, usually blank.
	SingleWord string `json:"singleword"` // SingleWord is the suggested username token.
	SearchTerm string `json:"searchterm"` // SearchTerm is the matched fragment from the user's input.
}
