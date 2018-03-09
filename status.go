package main

import (
	"log"
)

// runStatus returns the status of the environment API wherein the user is logged
func runStatus() {
	b, url, _ := isLogged()
	if !b {
		log.Println(LOG_LOGIN_REQUIRED)
	} else {
		// TODO get the status of the environment
		log.Printf(LOG_GETTING_STATUS, url)
	}
}
