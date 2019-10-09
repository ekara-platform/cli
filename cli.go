package main

import (
	"github.com/ekara-platform/cli/cmd"
	"github.com/ekara-platform/cli/common"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:generate go run generate/generate.go

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		common.ShowError("Interrupted by user after %s!", common.HumanizeDuration(time.Since(common.StartTime)))
		common.NoProgress = true
		cmd.StopCurrentContainerIfRunning()
		os.Exit(124)
	}()

	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			common.ShowError("Panic: %v", err)
			cmd.StopCurrentContainerIfRunning()
			os.Exit(125)
		} else {
			cmd.StopCurrentContainerIfRunning()
			os.Exit(0)
		}
	}()

	cmd.Execute()
}
