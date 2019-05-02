package main

import (
	"fmt"
	"os"

	"text/template"

	"github.com/ekara-platform/cli/image"
)

type Content struct {
	Version   string
	ImageName string
	Count     uint
}

func main() {

	c := Content{}
	c.ImageName = image.StarterImageName
	tag := os.Getenv("TRAVIS_TAG")
	if len(tag) > 0 {
		c.Version = tag
	} else {
		commit := os.Getenv("TRAVIS_COMMIT")
		if commit != "" {
			c.Version = "Commit:" + commit
		} else {
			c.Version = "unset"
		}
	}

	fmt.Printf("Generating the CLI version %s\n", c.Version)

	w, err := os.Create("./cmd/version.go")
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
