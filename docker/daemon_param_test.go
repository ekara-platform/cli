package docker

import (
	"log"
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoHost(t *testing.T) {
	logger := log.New(os.Stdout, "Ekara CLI: ", log.Ldate|log.Ltime)

	dp := DaemonParams{Host: ""}
	e := dp.checkAndLog(logger)
	assert.NotNil(t, e)
}
