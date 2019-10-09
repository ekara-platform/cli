package cmd

import (
	"github.com/ekara-platform/cli/common"
	"github.com/ekara-platform/cli/docker"
	"github.com/ekara-platform/engine/action"
	"github.com/ekara-platform/engine/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
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
		common.ShowWorking(common.LOG_APPLYING_ENV)
		if common.Flags.SSH.PrivateSSHKey != "" && common.Flags.SSH.PublicSSHKey != "" {
			// Move the ssh keys into the exchange folder input
			err := copyFile(common.Flags.SSH.PublicSSHKey, filepath.Join(ef.Input.Path(), util.SSHPuplicKeyFileName))
			if err != nil {
				common.ShowError("Error copying the SSH public key")
				os.Exit(1)
			}
			err = copyFile(common.Flags.SSH.PrivateSSHKey, filepath.Join(ef.Input.Path(), util.SSHPrivateKeyFileName))
			if err != nil {
				common.ShowError("Error copying the SSH private key")
				os.Exit(1)
			}
		}
		status := starterStart(args[0], ef, action.ApplyActionID)
		if status == 0 {
			common.ShowDone("Done in %s!", common.HumanizeDuration(time.Since(common.StartTime)))
		} else {
			common.ShowError("Errored (%d) after %s!", status, common.HumanizeDuration(time.Since(common.StartTime)))
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

func starterStart(url string, ef util.ExchangeFolder, action action.ActionID) int {
	docker.EnsureDockerInit()

	done := make(chan bool, 1)
	go docker.ImagePull(starterImageName, done)
	<-done

	if id, running := docker.ContainerRunningByImageName(starterImageName); running {
		text := common.ShowPrompt(common.PROMPT_RESTART)
		if strings.ToUpper(strings.TrimSpace(text)) == "Y" {
			done := make(chan bool, 1)
			go docker.StopContainerById(id, done)
			<-done
		} else {
			panic(errors.New(common.LOG_FAIL_ON_PROMPT_RESTART))
		}
	}

	done = make(chan bool, 1)
	status := docker.StartContainer(url, starterImageName, done, ef, action)
	<-done
	return status
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
