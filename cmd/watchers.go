package main

import (
	"fmt"
	"inkbunny/api"
	"log"
)

func main() {
	user, err := (&api.Credentials{
		Username: "",
		Password: "",
	}).Login()
	if err != nil {
		log.Fatalf("error logging in: %v", err)
	}
	log.Printf("logged in as %s, session id: %s", user.Username, user.Sid)

	var username string
	var watchers []string
	fmt.Print("Enter username to get watchers for: ")
	fmt.Scanln(&username)
	if username == "" || username == user.Username || username == "self" {
		watchers, err = user.GetWatchlist()
	} else {
		watchers, err = user.GetWatchers(username)
	}
	if err != nil {
		log.Fatalf("error getting watchers: %v", err)
	}

	log.Printf("watchers: %v", watchers)
}
