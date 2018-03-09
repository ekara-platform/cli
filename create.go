package main

import (
	"log"
)

// runCreate starts the installer in order to create an environement
func runCreate() {
	b, url, user := isLogged()
	if b {
		log.Printf(LOG_LOGGED_AS, user, url)
		log.Printf(LOG_LOGOUT_REQUIRED)
	} else {
		log.Printf(LOG_DEPLOYING_FROM, p.url)
		starterStart(true)
	}
}
