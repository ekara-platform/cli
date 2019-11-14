package cmd

import (
	"github.com/ekara-platform/cli/common"
	"github.com/ekara-platform/engine/action"
	"github.com/ekara-platform/engine/util"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func init() {
	// This is a descriptor-based command
	applyDescriptorFlags(destroyCmd)

	// Docker flags
	destroyCmd.PersistentFlags().StringVar(&common.Flags.Docker.Host, "docker-host", getEnvDockerHost(), "Docker daemon host")
	destroyCmd.PersistentFlags().StringVar(&common.Flags.Docker.Cert, "docker-cert-path", os.Getenv("DOCKER_CERT_PATH"), "Location of the Docker certificates")
	destroyCmd.PersistentFlags().BoolVar(&common.Flags.Docker.TLS, "docker-tls-verify", os.Getenv("DOCKER_TLS_VERIFY") == "", "If present TLS is enforced for Docker daemon communication")
	destroyCmd.PersistentFlags().StringVar(&common.Flags.Docker.APIVersion, "docker-api-version", os.Getenv("DOCKER_API_VERSION"), "Docker daemon API version")

	// SSH flags
	destroyCmd.PersistentFlags().StringVar(&common.Flags.SSH.PublicSSHKey, "public-key", "", "Custom public SSH key for the environment")
	destroyCmd.PersistentFlags().StringVar(&common.Flags.SSH.PrivateSSHKey, "private-key", "", "Custom private SSH key for the environment")

	rootCmd.AddCommand(destroyCmd)
}

var destroyCmd = &cobra.Command{
	Use:   "destroy <repository-url>",
	Short: "Destroy the existing environment infrastructure.",
	Long:  `The destroy command will ensure that every resource from the environment is destroyed.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		color.New(color.FgHiWhite).Println(common.LOG_DESTROYING_ENV)
		if common.Flags.SSH.PrivateSSHKey != "" && common.Flags.SSH.PublicSSHKey != "" {
			// Move the ssh keys into the exchange folder input
			err := copyFile(common.Flags.SSH.PublicSSHKey, filepath.Join(ef.Input.Path(), util.SSHPublicKeyFileName))
			if err != nil {
				common.CliFeedbackNotifier.Error("Error copying the SSH public key")
				os.Exit(1)
			}
			err = copyFile(common.Flags.SSH.PrivateSSHKey, filepath.Join(ef.Input.Path(), util.SSHPrivateKeyFileName))
			if err != nil {
				common.CliFeedbackNotifier.Error("Error copying the SSH private key")
				os.Exit(1)
			}
		}
		status, err := execAndWait(args[0], ef, action.DestroyActionID)
		if err != nil {
			common.CliFeedbackNotifier.Error("Unable to start installer: %s", err.Error())
			return
		}

		if status == 0 {
			common.CliFeedbackNotifier.Info("Destroy done in %s!", common.HumanizeDuration(time.Since(common.StartTime)))
			if ef.Output.Contains("result.json") {
				result, err := ioutil.ReadFile(filepath.Join(ef.Output.AdaptedPath(), "result.json"))
				if err != nil {
					common.CliFeedbackNotifier.Error("Unable to read destroy result: %s", err.Error())
					return
				}
				destroyResult := action.DestroyResult{}
				err = destroyResult.FromJson(string(result))
				if err != nil {
					common.CliFeedbackNotifier.Error("Unable to parse destroy result: %s", err.Error())
					return
				}
			}
		} else {
			common.CliFeedbackNotifier.Error("Errored (%d) after %s!", status, common.HumanizeDuration(time.Since(common.StartTime)))
		}
	},
}
