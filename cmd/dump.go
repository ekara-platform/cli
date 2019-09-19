package cmd

import (
	"fmt"
	"github.com/ekara-platform/cli/common"
	"github.com/ekara-platform/engine/action"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

const (
	dumpExchangeFolder string = "dump"
)

func init() {
	// This is a descriptor-based command
	applyDescriptorFlags(dumpCmd)

	rootCmd.AddCommand(dumpCmd)
}

var dumpCmd = &cobra.Command{
	Use:   "dump <repository-url>",
	Short: "Dump an existing environment descriptor.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		common.Logger.Printf(common.LOG_VALIDATING_ENV, args[0])
		dir, err := ioutil.TempDir(os.TempDir(), "ekara_dump")
		if err != nil {
			fmt.Println("unable to create temporary directory", err)
			os.Exit(1)
		}
		defer os.RemoveAll(dir)

		e := initLocalEngine(dir, args[0])
		res, err := e.ActionManager().Run(action.DumpActionID)
		if err != nil {
			fmt.Println("unable to run dump action", err)
			os.Exit(1)
		}

		text, err := res.AsPlainText()
		if err != nil {
			fmt.Println("unable to format text result from dump action", err)
			os.Exit(1)
		}
		if len(text) > 0 {
			for _, line := range text {
				fmt.Println(line)
			}
		} else {
			fmt.Println("No result")
		}

		if !res.IsSuccess() {
			os.Exit(2)
		}
		os.Exit(0)
	},
}
