package common

import (
	"fmt"
	"log"
)

// Flags holds actual CLI flag values
var Flags = AllFlags{}

// AllFlags regroups all possible CLI flags
type AllFlags struct {
	Debug      bool
	Docker     DockerFlags
	Descriptor DescriptorFlags
	Logging    LoggingFlags
	Proxy      ProxyFlags
	SSH        SSHFlags
}

func (p AllFlags) checkAndLog(logger *log.Logger) error {
	if e := p.SSH.checkAndLog(logger); e != nil {
		return e
	}
	return nil
}

// DescriptorFlags regroups descriptor-related flags
type DescriptorFlags struct {
	// The name of the environment descriptor
	File string
	// The name of the parameters files
	ParamFile string
	// The login required to reach the descriptor
	Login string
	// The password required to reach the descriptor
	Password string
}

// LoggingFlags regroups logging-related flags
type LoggingFlags struct {
	Verbose     bool
	VeryVerbose bool
	File        string
}

// ShouldOutputLogs returns true if (very) verbose mode is enabled
func (l LoggingFlags) ShouldOutputLogs() bool {
	return l.Verbose || l.VeryVerbose
}

// VerbosityLevel returns the numeric verbosity level (0, 1 or 2)
func (l LoggingFlags) VerbosityLevel() int {
	verbosity := 0
	if l.VeryVerbose {
		verbosity = 2
	} else if l.Verbose {
		verbosity = 1
	}
	return verbosity
}

// SSHFlags regroups SSH-related flags
type SSHFlags struct {
	// The public SSH key used to log on the created environment
	PublicSSHKey string
	// The private SSH key used to log on the created environment
	PrivateSSHKey string
}

func (s SSHFlags) checkAndLog(logger *log.Logger) error {
	// The SSH public key always comes with the private
	if s.PrivateSSHKey != "" || s.PublicSSHKey != "" {
		if s.PrivateSSHKey == "" {
			return fmt.Errorf(ERROR_REQUIRED_SSH_PRIVATE)
		}
		if s.PublicSSHKey == "" {
			return fmt.Errorf(ERROR_REQUIRED_SSH_PUBLIC)
		}
		logger.Printf(LOG_SSH_CONFIRMATION)
	}
	return nil
}

// DockerFlags regroups docker-related flags
type DockerFlags struct {
	// The docker host
	Host string
	// The docker certificates location
	Cert string
	// TLS using for daemon communication
	TLS bool
	// Docker daemon API version if any
	APIVersion string
}

// ProxyFlags regroups proxy-related flags
type ProxyFlags struct {
	HTTP       string
	HTTPS      string
	Exclusions string
}
