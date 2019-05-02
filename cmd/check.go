package cmd

import (
	"github.com/ekara-platform/cli/folder"
	"github.com/ekara-platform/cli/message"
	"github.com/ekara-platform/engine"
	"github.com/spf13/cobra"
)

const (
	//CheckExchangeFolder is the folder where the check result will be written
	CheckExchangeFolder string = "check"
)

var checkCmd = &cobra.Command{
	Use:   "check <descriptor-repository-url>",
	Short: "Validate an existing environment descriptor.",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Printf(message.LOG_CHECKING_FROM, cr.Descriptor.Url)
		ef := folder.CreateEF(CheckExchangeFolder, logger)
		cr.User.Output = true
		starterStart(*ef, "check", cr.Descriptor.Url, cr.Descriptor.File, engine.ActionCheckID, cr)
	},
	PersistentPreRun:  showHeader,
	PersistentPostRun: logDone,
	Args:              cobra.ExactArgs(1),
}
