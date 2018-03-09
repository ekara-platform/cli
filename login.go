package main

import (
	"fmt"
	"log"
	"net/url"
	"os/user"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Parameters required to login into the envirment manager API
type Login struct {
	url      string
	user     string
	password string
}

// checkLoginParams checks the coherence of the parameters received to deal with the
// environment API
func (l *Login) checkLoginParams(c *kingpin.ParseContext) error {
	b, _, _ := isLogged()
	if !b {
		checkFlag(l.url, apiUrlFlagKey)

		if _, err := url.ParseRequestURI(l.url); err != nil {
			log.Fatal(err.Error())
		}
		if l.user == "" {
			if u, e := user.Current(); e == nil {
				l.user = u.Username
			}
		}

		if l.user == "" {
			log.Fatal(fmt.Errorf(ERROR_NO_PROVIDED_USER, userFlagKey))
		}
		fmt.Printf(PROMPT_PASSWORD, l.user)
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		log.Printf("\n")
		if err != nil {
			log.Fatal(ERROR_READING_PASSWORD)
		}
		l.password = string(bytePassword)

	}
	return nil
}

// runLogin performs a login into the environment manager API
func runLogin() {
	b, url, user := isLogged()
	if b {
		log.Printf(LOG_ALREADY_LOGGED_AS, user, url)
	} else {
		log.Printf(LOG_LOGGING_INTO, l.url)
		// TODO
		// - Connect to the api with the received url/user/password
		// - Save received token into the session file
		s := &Session{
			Url:   l.url,
			User:  l.user,
			Token: "save token here",
		}
		saveLogged(*s)
	}
}
