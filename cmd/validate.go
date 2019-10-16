package cmd

import (
	"fmt"
	"github.com/ekara-platform/model"
	"github.com/fatih/color"
	"io/ioutil"
	"math"
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

		errCount := len(res.(action.ValidateResult).Errors)
		if errCount > 0 {
			fmt.Println("Validation problem(s) were encountered")
			for _, vErr := range res.(action.ValidateResult).Errors {
				if vErr.ErrorType == model.Error {
					color.New(color.FgHiRed).Printf("ERROR " + vErr.Message)
				} else {
					color.New(color.FgHiYellow).Printf("WARN  " + vErr.Message)
				}
			}
			os.Exit(int(math.Min(float64(errCount), 99)))
		} else {
			color.New(color.FgHiWhite).Println("No validation problem encountered")
			os.Exit(0)
		}
	},
}
