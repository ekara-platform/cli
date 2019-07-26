package cmd

import (
	"github.com/ekara-platform/cli/folder"
	"github.com/ekara-platform/cli/message"
	"github.com/ekara-platform/engine"
	"github.com/spf13/cobra"
)

const (
	dumpExchangeFolder string = "dump"
)

var dumpCmd = &cobra.Command{
	Use:   "dump <descriptor-repository-url>",
	Short: "Dump an existing environment descriptor.",
	Run: func(cmd *cobra.Command, args []string) {

		logger.Printf(message.LOG_DUMPING_FROM, cr.Descriptor.Url)
		ef := folder.CreateEF(dumpExchangeFolder, logger)
		cr.User.Output = true
		starterStart(*ef, "dump", cr.Descriptor, engine.ActionDumpID, cr)
	},
	PersistentPreRun:  showHeader,
	PersistentPostRun: logDone,
	Args:              cobra.ExactArgs(1),
}
