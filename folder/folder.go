package folder

import (
	"fmt"
	"log"

	"github.com/ekara-platform/cli/message"
	"github.com/ekara-platform/engine/util"
)

const (
	ROOT_EXCHANGE_FOLDER string = "out"
)

func CreateEF(folder string, logger *log.Logger) *util.ExchangeFolder {
	ef, e := util.CreateExchangeFolder(ROOT_EXCHANGE_FOLDER, folder)
	if e != nil {
		logger.Fatal(fmt.Errorf(message.ERROR_CREATING_EXCHANGE_FOLDER, folder))
	}
	e = ef.Create()
	if e != nil {
		logger.Fatal(fmt.Errorf(message.ERROR_CREATING_EXCHANGE_FOLDER, e.Error()))
	}
	return ef
}
