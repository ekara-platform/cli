package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ekara-platform/cli/cmd"
)

//go:generate go run generate/generate.go

const ()

func main() {
	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
	}()
	logger := log.New(os.Stdout, "Ekara CLI: ", log.Ldate|log.Ltime)
	cmd.Execute(logger)

}
