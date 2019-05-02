package docker

import (
	"log"
)

type InstallerParams struct {
	HttpProxy  string
	HttpsProxy string
	NoProxy    string
}

func (u InstallerParams) checkAndLog(logger *log.Logger) error {
	return nil
}
