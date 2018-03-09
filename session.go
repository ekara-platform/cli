package main

import (
	"encoding/json"
	"log"
	"os"
)

const (
	// Name of the lagoon persisted session file
	sessionFileName string = "lagoon_session.cli"
)

// Structure of the file containing the session details
type Session struct {
	Url   string `json:"api-url"`
	User  string `json:"logged_user"`
	Token string `json:"token"`
}

// isLogged returns true if the user is already logged in a environment manager API;
// If the user is logged it will also return the url where the login occured
// and the logged user
func isLogged() (logged bool, url string, user string) {

	if _, err := os.Stat(fullSessionFileName); err == nil {
		if data, err := os.Open(fullSessionFileName); err == nil {
			var s Session
			defer data.Close()
			err = json.NewDecoder(data).Decode(&s)
			if err != nil {
				log.Fatal(err.Error())
			}
			logged = true
			url = s.Url
			user = s.User
			return
		} else {
			log.Fatal(err.Error())
		}
	}
	logged = false
	url = ""
	user = ""
	return

}

// saveLogged saves the session details into the session file
func saveLogged(s Session) {
	b, err := json.Marshal(s)
	if err != nil {
		// TODO add real error message here
		log.Fatal(err.Error())
	}
	f, err := os.Create(fullSessionFileName)
	if err != nil {
		// TODO add real error message here
		log.Fatal(err.Error())
	}

	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		// TODO add real error message here
		log.Fatal(err.Error())
	}
}
