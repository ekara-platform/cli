package main

import (
	"encoding/json"
	"log"
	"os"
)

type Session struct {
	Url   string `json:"api-url"`
	User  string `json:"logged_user"`
	Token string `json:"token"`
}

func isLogged() (bool, string, string) {
	if _, err := os.Stat(fullSessionFileName); err == nil {
		if data, err := os.Open(fullSessionFileName); err == nil {
			var s Session
			defer data.Close()
			err = json.NewDecoder(data).Decode(&s)
			if err != nil {
				log.Fatal(err.Error())
			}
			return true, s.Url, s.User
		} else {
			log.Fatal(err.Error())
		}
	}
	return false, "", ""
}

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
