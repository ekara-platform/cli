package docker

import (
	"log"

	"github.com/ekara-platform/cli/message"
)

const (
	envDockerHost string = "DOCKER_HOST"
)

type DaemonParams struct {
	// The docker host
	Host string
	// The docker certificates location
	Cert string
}

func (d DaemonParams) checkAndLog(logger *log.Logger) error {
	// If we use flags then these 2 are required
	if d.Host != "" || d.Cert != "" {
		if e := checkFlag(d.Host, "host"); e != nil {
			return e
		}
		if d.Cert != "" {
			logger.Printf(message.LOG_FLAG_CONFIRMATION, "cert", d.Cert)
		}
		logger.Printf(message.LOG_FLAG_CONFIRMATION, "host", d.Host)
		logger.Printf(message.LOG_INIT_FLAGGED_DOCKER_CLIENT)
		initFlaggedClient(d.Host, d.Cert)
	} else {
		// if the flags are not used then we will ensure
		// that the environment variables are well defined
		if e := checkEnvVar(envDockerHost); e != nil {
			return e
		}
		logger.Printf(message.LOG_INIT_DOCKER_CLIENT)
		initClient()
	}
	return nil
}
