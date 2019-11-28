package cmd

import (
	"errors"
	"github.com/ekara-platform/cli/common"
	"github.com/ekara-platform/cli/docker"
	"github.com/ekara-platform/engine/action"
	"github.com/ekara-platform/engine/util"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	starterImageName         string = "ekaraplatform/installer:latest"
	defaultWindowsDockerHost string = "npipe:////./pipe/docker_engine"
	defaultUnixDockerHost    string = "unix:///var/run/docker.sock"
)

func init() {
	// This is a descriptor-based command
	applyDescriptorFlags(applyCmd)

	// Skipping flags
	applyCmd.PersistentFlags().BoolVar(&common.Flags.Skipping.SkipCreate, "skip-create", false, "If true, apply will forcibly skip infrastructure creation")
	applyCmd.PersistentFlags().BoolVar(&common.Flags.Skipping.SkipInstall, "skip-install", false, "If true, apply will forcibly skip infrastructure creation and orchestrator installation")
	applyCmd.PersistentFlags().BoolVar(&common.Flags.Skipping.SkipDeploy, "skip-deploy", false, "If true, apply will forcibly skip infrastructure creation, orchestrator installation and stack deployment")

	// Docker flags
	applyCmd.PersistentFlags().StringVar(&common.Flags.Docker.Host, "docker-host", getEnvDockerHost(), "Docker daemon host")
	applyCmd.PersistentFlags().StringVar(&common.Flags.Docker.Cert, "docker-cert-path", os.Getenv("DOCKER_CERT_PATH"), "Location of the Docker certificates")
	applyCmd.PersistentFlags().BoolVar(&common.Flags.Docker.TLS, "docker-tls-verify", os.Getenv("DOCKER_TLS_VERIFY") == "", "If present TLS is enforced for Docker daemon communication")
	applyCmd.PersistentFlags().StringVar(&common.Flags.Docker.APIVersion, "docker-api-version", os.Getenv("DOCKER_API_VERSION"), "Docker daemon API version")

	// SSH flags
	applyCmd.PersistentFlags().StringVar(&common.Flags.SSH.PublicSSHKey, "public-key", "", "Custom public SSH key for the environment")
	applyCmd.PersistentFlags().StringVar(&common.Flags.SSH.PrivateSSHKey, "private-key", "", "Custom private SSH key for the environment")

	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply <repository-url>",
	Short: "Apply the descriptor to obtain the desired environment.",
	Long:  `The apply command will ensure that everything declared in the descriptor matches reality by taking the necessary actions.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		color.New(color.FgHiWhite).Println(common.LOG_APPLYING_ENV)
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
		status, err := execAndWait(args[0], ef, action.ApplyActionID)
		if err != nil {
			common.CliFeedbackNotifier.Error("Unable to start installer: %s", err.Error())
			return
		}

		if status == 0 {
			common.CliFeedbackNotifier.Info("Apply done in %s!", common.HumanizeDuration(time.Since(common.StartTime)))
			if ef.Output.Contains("result.json") {
				result, err := ioutil.ReadFile(filepath.Join(ef.Output.AdaptedPath(), "result.json"))
				if err != nil {
					common.CliFeedbackNotifier.Error("Unable to read apply result: %s", err.Error())
					return
				}
				applyResult := action.ApplyResult{}
				err = applyResult.FromJson(string(result))
				if err != nil {
					common.CliFeedbackNotifier.Error("Unable to parse apply result: %s", err.Error())
					return
				}
				showInventory(applyResult.Inventory)
			}
		} else {
			common.CliFeedbackNotifier.Error("Errored (%d) after %s!", status, common.HumanizeDuration(time.Since(common.StartTime)))
		}
	},
}

func getEnvDockerHost() string {
	host := os.Getenv("DOCKER_HOST")
	if host == "" {
		if runtime.GOOS == "windows" {
			host = defaultWindowsDockerHost
		} else {
			host = defaultUnixDockerHost
		}
	}
	return host
}

func execAndWait(url string, ef util.ExchangeFolder, action action.ActionID) (int, error) {
	docker.EnsureDockerInit()

	// check if container already running
	id, running, err := docker.ContainerRunningByImageName(starterImageName)
	if err != nil {
		return 0, err
	}

	if running {
		text := common.CliFeedbackNotifier.Prompt(common.PROMPT_RESTART)
		if strings.ToUpper(strings.TrimSpace(text)) == "Y" {
			done := make(chan bool, 1)
			go func() {
				if err := docker.StopContainerById(id, done); err != nil {
					common.CliFeedbackNotifier.Error("Unable to stop running container: %s", err.Error())
				}
			}()
			<-done
		} else {
			return 0, errors.New(common.LOG_FAIL_ON_PROMPT_RESTART)
		}
	}

	// check again
	id, running, err = docker.ContainerRunningByImageName(starterImageName)
	if err != nil {
		return 0, err
	}

	// pull the image if necessary
	done := make(chan bool, 1)
	failed := make(chan error, 1)
	go docker.ImagePull(starterImageName, done, failed)
	select {
	case err = <-failed:
		return 0, err
	case <-done:
	}

	// start the container
	done = make(chan bool, 1)
	status, err := docker.StartContainer(url, starterImageName, done, ef, action)
	if err != nil {
		common.CliFeedbackNotifier.Error("Unable to start container: %s", err.Error())
		return 0, err
	}
	<-done

	return status, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
