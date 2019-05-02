package docker

import (
	"log"
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoUrl(t *testing.T) {
	logger := log.New(os.Stdout, "Ekara CLI: ", log.Ldate|log.Ltime)

	dp := DescriptorParams{Url: ""}
	e := dp.checkAndLog(logger)
	assert.NotNil(t, e)
}

func TestOkay(t *testing.T) {
	logger := log.New(os.Stdout, "Ekara CLI: ", log.Ldate|log.Ltime)

	dp := DescriptorParams{Url: "content"}
	e := dp.checkAndLog(logger)
	assert.Nil(t, e)
}
