package docker

import (
	"fmt"
	"log"

	"github.com/ekara-platform/cli/message"
	"github.com/ekara-platform/model"
)

type DescriptorParams struct {
	// The url of the repository containing the environment descriptor
	Url string
	// The name of the environment descriptor
	File string
	// The name of the parameters files
	ParamFile string
}

func (d DescriptorParams) checkAndLog(logger *log.Logger) error {
	// The environment descriptor is always required
	if d.Url == "" {
		return fmt.Errorf(message.ERROR_REQUIRED_DESCRIPTOR_URL)
	} else {
		logger.Printf(message.LOG_CONFIG_CONFIRMATION, "descriptor location", d.Url)
	}

	if d.File != "" {
		logger.Printf(message.LOG_CONFIG_CONFIRMATION, "descriptor name", d.File)
	} else {
		logger.Printf(message.LOG_CONFIG_CONFIRMATION, "descriptor name", model.DefaultDescriptorName)
	}

	if d.ParamFile != "" {
		logger.Printf(message.LOG_CONFIG_CONFIRMATION, "descriptor parameters", d.ParamFile)
	}

	return nil
}
