package docker

import (
	"fmt"
	"log"

	"github.com/ekara-platform/cli/message"
)

type HostParams struct {
	// The public SSH key used to log on the created environment
	PublicSSHKey string
	// The private SSH key used to log on the created environment
	PrivateSSHKey string
}

func (h HostParams) checkAndLog(logger *log.Logger) error {
	// The SSH public key always comes with the private
	if h.PrivateSSHKey != "" || h.PublicSSHKey != "" {

		if h.PrivateSSHKey == "" {
			return fmt.Errorf(message.ERROR_REQUIRED_SSH_PRIVATE)
		}

		if h.PublicSSHKey == "" {
			return fmt.Errorf(message.ERROR_REQUIRED_SSH_PUBLIC)
		}
		logger.Printf(message.LOG_SSH_PUBLIC_CONFIRMATION, h.PublicSSHKey)
		logger.Printf(message.LOG_SSH_PRIVATE_CONFIRMATION, h.PrivateSSHKey)
	}
	return nil
}
