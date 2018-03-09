package main

import (
	"log"
)

// runUpdate starts the installer in order to update an environment
// The user must be logged into the environment which he wants to update
func runUpdate() {
	b, url, _ := isLogged()
	if b {
		log.Printf(LOG_UPDATING_FROM, url)
		starterStart(false)
	} else {
		log.Printf(LOG_LOGIN_REQUIRED)
	}
}
