package cmd

import (
	"github.com/ekara-platform/engine/action"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestValidation(t *testing.T) {
	e := initLocalEngine("out/engine", "https://github.com/ekara-platform/demo")
	defer os.RemoveAll("out/engine")
	res, err := e.ActionManager().Run(action.ValidateActionID)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestDump(t *testing.T) {
	e := initLocalEngine("out/engine", "https://github.com/ekara-platform/demo")
	defer os.RemoveAll("out/engine")
	res, err := e.ActionManager().Run(action.DumpActionID)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}
