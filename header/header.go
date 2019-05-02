package header

import (
	"log"

	"github.com/ekara-platform/cli/image"
	"github.com/ekara-platform/cli/message"
)

func ShowHeader() {

	log.Printf(message.LOG_CLI_IMAGE, image.StarterImageName)

	// this comes from http://www.kammerl.de/ascii/AsciiSignature.php
	// the font used id "standard"
	log.Println(" _____ _                   ")
	log.Println("| ____| | ____ _ _ __ __ _ ")
	log.Println("|  _| | |/ / _` | '__/ _` |")
	log.Println("| |___|   < (_| | | | (_| |")
	log.Println(`|_____|_|\_\__,_|_|  \__,_|`)

	log.Println(`  ____ _     ___ `)
	log.Println(` / ___| |   |_ _|`)
	log.Println(`| |   | |    | | `)
	log.Println(`| |___| |___ | | `)
	log.Println(` \____|_____|___|`)

}
