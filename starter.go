package main

import (
	"bufio"
	"fmt"
	"log"

	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	// Environment variables used by default by the docker client
	// "github.com/docker/docker/client"
	envCertPath   string = "DOCKER_CERT_PATH"
	envApiVersion string = "DOCKER_API_VERSION"
	envDockerHost string = "DOCKER_HOST"

	// Flags keys for Commands
	deployFlagKey = "create"
	updateFlagKey = "update"
	loginFlagKey  = "login"
	logoutFlagKey = "logout"
	statusFlagKey = "status"

	// Flags keys for Arguments
	descriptorFlagKey     = "descriptor"
	environmentUrlFlagKey = "url"
	certPathFlagKey       = "cert"
	apiVersionFlagKey     = "api"
	dockerHostFlagKey     = "host"
	userFlagKey           = "user"
	apiUrlFlagKey         = "url"

	// Name of the lagoon starter image
	starterImageName string = "redis:latest"
)

var (
	// Commands
	deploy *kingpin.CmdClause
	update *kingpin.CmdClause
	login  *kingpin.CmdClause
	logout *kingpin.CmdClause
	status *kingpin.CmdClause

	fullSessionFileName string

	// Arguments
	p *DockerParams
	l *Login
)

func initFlags(app *kingpin.Application) {

	p = &DockerParams{}
	deploy = app.Command(deployFlagKey, "Create a new environment.")
	deploy.Arg(descriptorFlagKey, "The environment descriptor url").Required().StringVar(&p.url)
	deploy.Flag(certPathFlagKey, "The location of the docker certificates (optional)").StringVar(&p.cert)
	deploy.Flag(apiVersionFlagKey, "The version of the docker API (optional)").StringVar(&p.api)
	deploy.Flag(dockerHostFlagKey, "The url of the docker host (optional)").StringVar(&p.host)
	deploy.Action(p.checkDockerParams)

	update = app.Command(updateFlagKey, "Update an existing environment.")
	update.Arg(descriptorFlagKey, "The environment descriptor url").Required().StringVar(&p.url)
	update.Flag(certPathFlagKey, "The location of the docker certificates (optional)").StringVar(&p.cert)
	update.Flag(apiVersionFlagKey, "The version of the docker API (optional)").StringVar(&p.api)
	update.Flag(dockerHostFlagKey, "The url of the docker host (optional)").StringVar(&p.host)
	deploy.Action(p.checkDockerParams)

	l = &Login{}
	login = app.Command(loginFlagKey, "Login into an environment manager API.")
	login.Arg(apiUrlFlagKey, "The url of the environment manager API").Required().StringVar(&l.url)
	login.Flag(userFlagKey, "The user (optional)").StringVar(&l.user)
	login.Action(l.checkLoginParams)

	logout = app.Command(logoutFlagKey, "Logout from an environment manager API.")

	status = app.Command(statusFlagKey, "Status of the environment manager API.")
}

func main() {

	fullSessionFileName = "./" + sessionFileName
	// this comes from http://www.kammerl.de/ascii/AsciiSignature.php
	// the font used id "standard"
	if _, err := os.Stat(fullSessionFileName); os.IsNotExist(err) {
		log.Println(` _                                  `)
		log.Println(`| |    __ _  __ _  ___   ___  _ __  `)
		log.Println(`| |   / _  |/ _  |/ _ \ / _ \| '_ \ `)
		log.Println(`| |__| (_| | (_| | (_) | (_) | | | |`)
		log.Println(`|_____\__,_|\__, |\___/ \___/|_| |_|`)
		log.Println(`            |___/                   `)

		log.Println(`  ____ _     ___ `)
		log.Println(` / ___| |   |_ _|`)
		log.Println(`| |   | |    | | `)
		log.Println(`| |___| |___ | | `)
		log.Println(` \____|_____|___|`)
	}

	app := kingpin.New("lagoon", CLI_DESCRIPTION)
	initFlags(app)
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case deploy.FullCommand():
		runCreate()
	case update.FullCommand():
		runUpdate()
	case login.FullCommand():
		runLogin()
	case logout.FullCommand():
		runLogout()
	case status.FullCommand():
		runStatus()
	}
	log.Println(LOG_COMMAND_COMPLETED)
}

func starterStart(create bool) {
	log.Printf(LOG_GET_IMAGE)
	done := make(chan bool, 1)
	go imagePull(starterImageName, done)
	<-done

	if id, running := containerRunningByImageName(starterImageName); running {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(PROMPT_RESTART)
		text, _ := reader.ReadString('\n')
		if strings.TrimSpace(text) == "Y" {
			done := make(chan bool, 1)
			go stopContainerById(id, done)
			<-done
		} else {
			log.Printf(LOG_FAIL_ON_PROMPT_RESTART)
			return
		}
	}

	done = make(chan bool, 1)
	startContainer(starterImageName, done, create)
	<-done
	log.Printf(LOG_OK_STARTED)
}

func checkFlag(val string, flagKey string) {
	if val == "" {
		log.Fatal(fmt.Errorf(ERROR_REQUIRED_FLAG, flagKey))
	}
}

func checkEnvVar(key string) {
	if os.Getenv(key) == "" {
		log.Fatal(fmt.Errorf(ERROR_REQUIRED_ENV, key))
	}
}
