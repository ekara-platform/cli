package common

import (
	"log"
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoSSH(t *testing.T) {
	logger := log.New(os.Stdout, "Ekara CLI: ", log.Ldate|log.Ltime)

	hp := SSHFlags{PublicSSHKey: "", PrivateSSHKey: ""}
	e := hp.checkAndLog(logger)
	assert.Nil(t, e)
}

func TestSSHOkay(t *testing.T) {
	logger := log.New(os.Stdout, "Ekara CLI: ", log.Ldate|log.Ltime)

	hp := SSHFlags{PublicSSHKey: "content", PrivateSSHKey: "content"}
	e := hp.checkAndLog(logger)
	assert.Nil(t, e)
}

func TestSSHKoNoPrivate(t *testing.T) {
	logger := log.New(os.Stdout, "Ekara CLI: ", log.Ldate|log.Ltime)

	hp := SSHFlags{PublicSSHKey: "content", PrivateSSHKey: ""}
	e := hp.checkAndLog(logger)
	assert.NotNil(t, e)
}

func TestSSHKoNoPublic(t *testing.T) {
	logger := log.New(os.Stdout, "Ekara CLI: ", log.Ldate|log.Ltime)

	hp := SSHFlags{PublicSSHKey: "", PrivateSSHKey: "content"}
	e := hp.checkAndLog(logger)
	assert.NotNil(t, e)
}
