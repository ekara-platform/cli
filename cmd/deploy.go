package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/ekara-platform/cli/folder"
	"github.com/ekara-platform/cli/message"
	"github.com/ekara-platform/engine"
	"github.com/ekara-platform/engine/util"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use: "deploy <descriptor-repository-url>",

	Short: "Deploy a new environment.",
	Long: `The deploy command will provision the environment nodes,
install the orchestrator and deploy stacks.`,
	Run: func(cmd *cobra.Command, args []string) {

		qName := parseHeader()
		ef := folder.CreateEF(qName, logger)

		logger.Printf(message.LOG_DEPLOYING_FROM, cr.Descriptor.Url)

		if cr.Host.PrivateSSHKey != "" && cr.Host.PublicSSHKey != "" {
			// Move the ssh keys into the exchange folder input
			err := Copy(cr.Host.PublicSSHKey, filepath.Join(ef.Input.Path(), util.SSHPuplicKeyFileName))
			if err != nil {
				logger.Fatal(fmt.Errorf(message.ERROR_COPYING_SSH_PUB, cr.Host.PublicSSHKey))
			}

			err = Copy(cr.Host.PrivateSSHKey, filepath.Join(ef.Input.Path(), util.SSHPrivateKeyFileName))
			if err != nil {
				logger.Fatal(fmt.Errorf(message.ERROR_COPYING_SSH_PRIV, cr.Host.PrivateSSHKey))
			}
		}
		starterStart(*ef, qName, cr.Descriptor, engine.ActionDeployID, cr)

	},
	PersistentPreRun:  showHeader,
	PersistentPostRun: logDone,
	Args:              cobra.ExactArgs(1),
}
