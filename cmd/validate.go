package cmd

import (
	"fmt"
	"github.com/fatih/color"
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
		color.New(color.FgHiWhite).Println(common.LOG_VALIDATING_ENV)
		dir, err := ioutil.TempDir(os.TempDir(), "ekara_validate")
		if err != nil {
			common.CliFeedbackNotifier.Error("Unable to create temporary directory: %s", err.Error())
			os.Exit(1)
		}
		defer os.RemoveAll(dir)

		e := initLocalEngine(dir, args[0])
		res, err := e.ActionManager().Run(action.ValidateActionID)
		if err != nil {
			common.CliFeedbackNotifier.Error("Unable to run validate action: %s", err.Error())
			os.Exit(1)
		}

		text, err := res.AsPlainText()
		if err != nil {
			common.CliFeedbackNotifier.Error("Unable to format text result from validate action: %s", err.Error())
			os.Exit(1)
		}

		if len(text) > 0 {
			for _, line := range text {
				fmt.Println(line)
			}
			common.CliFeedbackNotifier.Info("Validation problem(s) were encountered")
			os.Exit(2)
		} else {
			common.CliFeedbackNotifier.Info("No validation problem encountered")
			os.Exit(0)
		}
	},
}
