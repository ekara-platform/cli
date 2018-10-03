package main

import (
	"fmt"
	"os"
	"text/template"
)

type Content struct {
	Version string
	Count   uint
}

func main() {

	c := Content{}

	c.Version = os.Getenv("LAGOON_CLI_VERSION")

	if c.Version == "" {
		c.Version = "unset"
	}
	fmt.Printf("Generating the CLI version %s\n", c.Version)

	w, err := os.Create("version.go")
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
