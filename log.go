package main

import (
	"os"

	"path/filepath"

	"github.com/ekara-platform/engine/util"
)

const (
	DefaultContainerLogFileName string = "container.log"
)

func ContainerLog(ef util.ExchangeFolder, fileName string) (*os.File, error) {
	var file string
	if fileName == "" {
		file = filepath.Join(ef.Output.Path(), DefaultContainerLogFileName)
	} else {
		file = filepath.Join(ef.Output.Path(), fileName)
	}
	f, e := os.Create(file)
	if e != nil {
		return nil, e
	}
	return f, nil
}
