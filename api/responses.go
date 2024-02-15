package api

type Credentials struct {
	Sid      string `json:"sid"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type LogoutResponse struct {
	Sid    string `json:"sid"`
	Logout string `json:"logout"`
}

type WatchlistResponse struct {
	Watches []struct {
		Username string `json:"username"`
	} `json:"watches"`
}

type User struct {
	ID         string `json:"id"`
	Value      string `json:"value"`
	Icon       string `json:"icon"`
	Info       string `json:"info"`
	SingleWord string `json:"singleword"`
	SearchTerm string `json:"searchterm"`
}
