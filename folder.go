package main

import (
	"fmt"

	"github.com/ekara-platform/engine/util"
)

func createEF(folder string) *util.ExchangeFolder {
	ef, e := util.CreateExchangeFolder(ROOT_EXCHANGE_FOLDER, folder)
	if e != nil {
		logger.Fatal(fmt.Errorf(ERROR_CREATING_EXCHANGE_FOLDER, folder))
	}
	e = ef.Create()
	if e != nil {
		logger.Fatal(fmt.Errorf(ERROR_CREATING_EXCHANGE_FOLDER, e.Error()))
	}
	return ef
}
