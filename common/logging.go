package common

import (
	"io/ioutil"
	"log"
)

var Logger = log.New(ioutil.Discard, "", log.LstdFlags)
