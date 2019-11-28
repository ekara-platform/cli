package main

import (
	"fmt"
	"os"

	"text/template"
)

// Content is the data structure used for go:generate templating
type Content struct {
	Version string
	Count   uint
}

func main() {
	c := Content{}
	tag := os.Getenv("TRAVIS_TAG")
	if len(tag) > 0 {
		c.Version = tag
	} else {
		commit := os.Getenv("TRAVIS_COMMIT")
		if commit != "" {
			c.Version = "Commit:" + commit
		} else {
			c.Version = "<none>"
		}
	}

	fmt.Printf("Generating the CLI version %s\n", c.Version)

	w, err := os.Create("cmd/version.go")
	if err != nil {
		panic(err)
	}
	defer w.Close()

	t, err := template.ParseFiles("generate/version.tmpl")
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, c)
	if err != nil {
		panic(err)
	}
}
