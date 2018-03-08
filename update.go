package main

import (
	"log"
)

func runUpdate() {
	b, url, _ := isLogged()
	if b {
		log.Printf(LOG_UPDATING_FROM, url)
		starterStart(false)
	} else {
		log.Printf(LOG_LOGIN_REQUIRED)
	}
}
