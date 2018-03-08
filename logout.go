package main

import (
	"fmt"
	"log"
	"os"
)

func runLogout() {
	b, _, _ := isLogged()
	if !b {
		log.Println(LOG_ALREADY_LOGGED_OUT)
	} else {
		// delete the saved file containing the session token
		if err := os.Remove(fullSessionFileName); err != nil {
			log.Fatal(fmt.Errorf(ERROR_SESSION_NOT_CLOSED, fullSessionFileName))
		}
	}
}
