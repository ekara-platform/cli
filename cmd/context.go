package cmd

import (
	"github.com/ekara-platform/cli/common"
	"log"

	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
)

type (
	//cliContext simulates the LaunchContext
	cliContext struct {
		fN                   util.FeedbackNotifier
		logger               *log.Logger
		ef                   util.ExchangeFolder
		descriptorName       string
		location             string
		user                 string
		password             string
		sshPublicKeyContent  string
		sshPrivateKeyContent string
		extVars              model.Parameters
	}
)

//Progress implements the corresponding method in LaunchContext
func (lC cliContext) Feedback() util.FeedbackNotifier {
	return lC.fN
}

//Skip implements the corresponding method in LaunchContext
func (lC cliContext) Skipping() int {
	return common.Flags.Skipping.SkippingLevel()
}

//Verbosity implements the corresponding method in LaunchContext
func (lC cliContext) Verbosity() int {
	return common.Flags.Logging.VerbosityLevel()
}

//DescriptorName implements the corresponding method in LaunchContext
func (lC cliContext) DescriptorName() string {
	return lC.descriptorName
}

//User implements the corresponding method in LaunchContext
func (lC cliContext) User() string {
	return lC.user
}

//Password implements the corresponding method in LaunchContext
func (lC cliContext) Password() string {
	return lC.password
}

//Log implements the corresponding method in LaunchContext
func (lC cliContext) Log() *log.Logger {
	return lC.logger
}

//Ef implements the corresponding method in LaunchContext
func (lC cliContext) Ef() util.ExchangeFolder {
	return lC.ef
}

//Location implements the corresponding method in LaunchContext
func (lC cliContext) Location() string {
	return lC.location
}

//Proxy returns launch context proxy settings
func (lC cliContext) Proxy() model.Proxy {
	return model.Proxy{}
}

//SSHPublicKey implements the corresponding method in LaunchContext
func (lC cliContext) SSHPublicKey() string {
	return lC.sshPublicKeyContent
}

//SSHPrivateKey implements the corresponding method in LaunchContext
func (lC cliContext) SSHPrivateKey() string {
	return lC.sshPrivateKeyContent
}

func (lC cliContext) ExternalVars() model.Parameters {
	return lC.extVars
}
