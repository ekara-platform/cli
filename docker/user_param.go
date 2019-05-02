package docker

import (
	"log"

	"github.com/ekara-platform/cli/message"
)

type UserParams struct {
	Output bool
	File   string
}

func (u UserParams) checkAndLog(logger *log.Logger) error {
	// The file can be used only if the output is required
	if !u.Output && u.File != "" {
		logger.Printf(message.LOG_OUTPUT_FILE_IGNORED, u.File)

	}

	if u.Output {
		if u.File != "" {
			logger.Printf(message.LOG_CONTAINER_LOG_IS, u.File)
		} else {
			logger.Printf(message.LOG_CONTAINER_LOG_IS, DefaultContainerLogFileName)
		}
	}
	return nil
}
