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
	defer func() {
		if err := recover(); err != nil { //catch
			common.CliFeedbackNotifier.Error("Panic: %v", err)
			cmd.StopCurrentContainerIfRunning()
			os.Exit(125)
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		common.CliFeedbackNotifier.Error("Interrupted by user after %s!", common.HumanizeDuration(time.Since(common.StartTime)))
		common.NoFeedback = true
		cmd.StopCurrentContainerIfRunning()
		os.Exit(124)
	}()

	if err := cmd.Execute(); err != nil {
		common.CliFeedbackNotifier.Error("Error: %s", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
