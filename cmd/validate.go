package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ekara-platform/cli/common"
	"github.com/ekara-platform/engine/action"
	"github.com/spf13/cobra"
)

func init() {
	// This is a descriptor-based command
	applyDescriptorFlags(validateCmd)

	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate <repository-url>",
	Short: "Validate an existing environment descriptor.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		common.Logger.Printf(common.LOG_VALIDATING_ENV, args[0])
		dir, err := ioutil.TempDir(os.TempDir(), "ekara_validate")
		if err != nil {
			fmt.Println("unable to create temporary directory", err)
			os.Exit(1)
		}
		defer os.RemoveAll(dir)

		e := initLocalEngine(dir, args[0])
		res, err := e.ActionManager().Run(action.ValidateActionID)
		if err != nil {
			fmt.Println("unable to run validate action", err)
			os.Exit(1)
		}

		text, err := res.AsPlainText()
		if err != nil {
			fmt.Println("unable to format text result from validate action", err)
			os.Exit(1)
		}
		if len(text) > 0 {
			common.Logger.Println("Validation error(s) and warning(s) follow")
			for _, line := range text {
				fmt.Println(line)
			}
		} else {
			fmt.Println("No validation error or warning encountered")
		}

		if res.IsSuccess() {
			os.Exit(2)
		}
		os.Exit(0)
	},
}
