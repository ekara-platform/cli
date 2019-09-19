package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ekara-platform/cli/cmd"
)

//go:generate go run generate/generate.go

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Fprintln(os.Stderr, "Interrupted by user")
		cmd.StopCurrentContainerIfRunning()
		os.Exit(1)
	}()

	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
	}()
	cmd.Execute()
}
