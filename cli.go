package main

import (
	"bufio"
	"fmt"
	"log"

	"os"
	"strings"

	"github.com/lagoon-platform/engine"
	"github.com/lagoon-platform/model"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	VALIDATION_OUTPUT_FILE = "validation_output.json"

	// Environment variables used by default by the docker client
	// "github.com/docker/docker/client"
	envCertPath   string = "DOCKER_CERT_PATH"
	envApiVersion string = "DOCKER_API_VERSION"
	envDockerHost string = "DOCKER_HOST"
	envHttpProxy  string = "HTTP_PROXY"
	envHttpsProxy string = "HTTPS_PROXY"

	// Flags keys for Commands
	deployFlagKey = "create"
	updateFlagKey = "update"
	checkFlagKey  = "check"
	loginFlagKey  = "login"
	logoutFlagKey = "logout"
	statusFlagKey = "status"

	// Flags keys for Arguments
	descriptorFlagKey     = "descriptor"
	environmentUrlFlagKey = "url"
	certPathFlagKey       = "cert"
	apiVersionFlagKey     = "api"
	dockerHostFlagKey     = "host"
	httpProxyFlagKey      = "http_proxy"
	httpsProxyFlagKey     = "https_proxy"
	userFlagKey           = "user"
	apiUrlFlagKey         = "url"

	containerFileFlagKey   = "file"
	containerOutputFlagKey = "output"

	// Name of the lagoon starter image
	starterImageName string = "lagoon-platform/installer:latest"
)

var (
	// Commands
	deploy *kingpin.CmdClause
	update *kingpin.CmdClause
	check  *kingpin.CmdClause
	login  *kingpin.CmdClause
	logout *kingpin.CmdClause
	status *kingpin.CmdClause

	fullSessionFileName string

	// Arguments
	p      *DockerParams
	c      *CheckParams
	l      *Login
	logger *log.Logger
)

func initFlags(app *kingpin.Application) {

	p = &DockerParams{}
	deploy = app.Command(deployFlagKey, "Create a new environment.")
	deploy.Arg(descriptorFlagKey, "The environment descriptor url").Required().StringVar(&p.url)
	deploy.Flag(certPathFlagKey, "The location of the docker certificates (optional)").StringVar(&p.cert)
	deploy.Flag(apiVersionFlagKey, "The version of the docker API (optional)").StringVar(&p.api)
	deploy.Flag(dockerHostFlagKey, "The url of the docker host (optional)").StringVar(&p.host)
	deploy.Flag(httpProxyFlagKey, "The http proxy(optional)").StringVar(&p.httpProxy)
	deploy.Flag(httpsProxyFlagKey, "The https proxy (optional)").StringVar(&p.httpsProxy)
	deploy.Flag(containerOutputFlagKey, "\"true\" to write the container logs into a local file, defaulted to  \"false\"").BoolVar(&p.output)
	deploy.Flag(containerFileFlagKey, "The output file where to write the logs, if missing the content will be written in \"container.log\"").StringVar(&p.file)
	deploy.Action(p.checkDockerParams)

	update = app.Command(updateFlagKey, "Update an existing environment.")
	update.Arg(descriptorFlagKey, "The environment descriptor url").Required().StringVar(&p.url)
	update.Flag(certPathFlagKey, "The location of the docker certificates (optional)").StringVar(&p.cert)
	update.Flag(apiVersionFlagKey, "The version of the docker API (optional)").StringVar(&p.api)
	update.Flag(dockerHostFlagKey, "The url of the docker host (optional)").StringVar(&p.host)
	deploy.Action(p.checkDockerParams)

	c = &CheckParams{}
	check = app.Command(checkFlagKey, "Valid an existing environment descriptor.")
	check.Arg(descriptorFlagKey, "The environment descriptor url").Required().StringVar(&c.url)

	l = &Login{}
	login = app.Command(loginFlagKey, "Login into an environment manager API.")
	login.Arg(apiUrlFlagKey, "The url of the environment manager API").Required().StringVar(&l.url)
	login.Flag(userFlagKey, "The user (optional)").StringVar(&l.user)
	login.Action(l.checkLoginParams)

	logout = app.Command(logoutFlagKey, "Logout from an environment manager API.")

	status = app.Command(statusFlagKey, "Status of the environment manager API.")
}

func main() {
	logger = log.New(os.Stdout, "Lagoon CLI: ", log.Ldate|log.Ltime)

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
	case check.FullCommand():
		runCheck()
	case login.FullCommand():
		runLogin()
	case logout.FullCommand():
		runLogout()
	case status.FullCommand():
		runStatus()
	}
	log.Println(LOG_COMMAND_COMPLETED)
}

// runCreate starts the installer in order to create an environement
func runCreate() {
	b, url, user := isLogged()
	if b {
		log.Printf(LOG_LOGGED_AS, user, url)
		log.Printf(LOG_LOGOUT_REQUIRED)
	} else {
		log.Printf(LOG_DEPLOYING_FROM, p.url)
		err := parseLocation(p.url)
		if err == nil {
			starterStart(p.url, true)
		} else {
			log.Fatalf(ERROR_PARSING_ENVIRONMENT, err.Error())
		}
	}
}

// runUpdate starts the installer in order to update an environment
// The user must be logged into the environment which he wants to update
func runUpdate() {
	b, url, _ := isLogged()
	if b {
		log.Printf(LOG_UPDATING_FROM, url)
		err := parseLocation(p.url)
		if err == nil {
			starterStart(p.url, false)
		} else {
			log.Fatalf(ERROR_PARSING_ENVIRONMENT, err.Error())
		}
	} else {
		log.Printf(LOG_LOGIN_REQUIRED)
	}
}

// Parameters required to check the environment descriptor contant
type CheckParams struct {
	url    string
	file   string
	output bool
}

// runCheck checks the validity of the environment descriptor content
func runCheck() {
	log.Printf(LOG_CHECKING_FROM, c.url)
	err := parseLocation(c.url)
	if err != nil {
		log.Fatalf(ERROR_PARSING_ENVIRONMENT, err.Error())
	}
}

func parseLocation(location string) error {
	log.Printf(LOG_PARSING)

	_, e := engine.Create(logger, location)

	vErrs, ok := e.(model.ValidationErrors)

	// if the error is not a "validation error" then we return it
	if e != nil && !ok {
		return e
	}

	if ok {
		// print both errors and warnings into the report file
		log.Printf(e.Error())
		log.Printf(LOG_SEE_REPORT_IN, VALIDATION_OUTPUT_FILE)
		f, e := os.Create(VALIDATION_OUTPUT_FILE)
		if e != nil {
			logger.Fatal(fmt.Errorf(ERROR_CREATING_REPORT_FILE, VALIDATION_OUTPUT_FILE))
		}
		defer f.Close()
		b, e := vErrs.JSonContent()
		if e != nil {
			logger.Fatal(e)
		}
		_, e = f.Write(b)
		if e != nil {
			logger.Fatal(e)
		}
		if vErrs.HasErrors() {
			// in case of validation error we also return the error
			return e
		}
	}
	return nil
}

func starterStart(location string, create bool) {
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
	startContainer(starterImageName, done, create, location)
	<-done
	log.Printf(LOG_OK_STARTED)
}

func getHttpProxy() string {
	if p.httpProxy == "" {
		return os.Getenv(envHttpProxy)
	}
	return p.httpProxy
}

func getHttpsProxy() string {
	if p.httpsProxy == "" {
		return os.Getenv(envHttpsProxy)
	}
	return p.httpsProxy
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
