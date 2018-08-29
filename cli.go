package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/lagoon-platform/engine"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (

	// Environment variables used by default by the docker client
	// "github.com/docker/docker/client"
	envCertPath   string = "DOCKER_CERT_PATH"
	envApiVersion string = "DOCKER_API_VERSION"
	envDockerHost string = "DOCKER_HOST"
	envHttpProxy  string = "HTTP_PROXY"
	envHttpsProxy string = "HTTPS_PROXY"
	envNoProxy    string = "NO_PROXY"

	// Flags keys for Commands
	deployFlagKey = "create"
	updateFlagKey = "update"
	checkFlagKey  = "check"
	loginFlagKey  = "login"
	logoutFlagKey = "logout"
	statusFlagKey = "status"

	// Flags keys for Arguments
	descriptorFlagKey     = "descriptor"
	descriptorNameFlagKey = "file"

	certPathFlagKey      = "cert"
	apiVersionFlagKey    = "api"
	dockerHostFlagKey    = "host"
	clientFlagKey        = "client"
	httpProxyFlagKey     = "http_proxy"
	httpsProxyFlagKey    = "https_proxy"
	publicSSHKeyFlagKey  = "public_ssh"
	privateSSHKeyFlagKey = "private_ssh"

	noProxyFlagKey   = "no_proxy"
	userFlagKey      = "user"
	apiUrlFlagKey    = "url"
	paramFileFlagKey = "param"

	containerFileFlagKey   = "logfile"
	containerOutputFlagKey = "output"

	// Name of the lagoon starter image
	starterImageName string = "lagoonplatform/installer:latest"
)

var (
	// Commands
	deploy *kingpin.CmdClause
	update *kingpin.CmdClause
	check  *kingpin.CmdClause
	login  *kingpin.CmdClause
	logout *kingpin.CmdClause
	status *kingpin.CmdClause

	fullLoginFileName string

	// Arguments
	cr     *DockerCreateParams
	up     *DockerUpdateParams
	ch     *DockerCheckParams
	l      *Login
	logger *log.Logger
)

func initFlags(app *kingpin.Application) {

	cr = &DockerCreateParams{}
	deploy = app.Command(deployFlagKey, "Create a new environment.")
	deploy.Arg(descriptorFlagKey, "The environment descriptor url (the root folder location)").Required().StringVar(&cr.url)
	deploy.Flag(descriptorNameFlagKey, "The name of the environment descriptor, if missing we will look for a descriptor named \""+engine.DescriptorFileName+"\"").Default(engine.DescriptorFileName).StringVar(&cr.file)
	deploy.Flag(clientFlagKey, "The name of the environment client (required)").StringVar(&cr.client)
	deploy.Flag(certPathFlagKey, "The location of the docker certificates (optional)").StringVar(&cr.cert)
	deploy.Flag(apiVersionFlagKey, "The version of the docker API (optional)").StringVar(&cr.api)
	deploy.Flag(dockerHostFlagKey, "The url of the docker host (optional)").StringVar(&cr.host)
	deploy.Flag(paramFileFlagKey, "The parameters file (optional)").StringVar(&cr.container.paramFile)
	deploy.Flag(httpProxyFlagKey, "The http proxy(optional)").StringVar(&cr.container.httpProxy)
	deploy.Flag(httpsProxyFlagKey, "The https proxy (optional)").StringVar(&cr.container.httpsProxy)
	deploy.Flag(noProxyFlagKey, "The no proxy (optional)").StringVar(&cr.container.noProxy)
	deploy.Flag(containerOutputFlagKey, "\"true\" to write the container logs into a local file, defaulted to  \"false\"").BoolVar(&cr.container.output)
	deploy.Flag(containerFileFlagKey, "The output file where to write the logs, if missing the log content will be written in \""+DefaultContainerLogFileName+"\"").StringVar(&cr.container.file)

	deploy.Flag(publicSSHKeyFlagKey, "The public SSH key to connect the created machines  (optional)").StringVar(&cr.publicSSHKey)
	deploy.Flag(privateSSHKeyFlagKey, "The private SSH key to connect the created machines  (optional)").StringVar(&cr.privateSSHKey)
	deploy.Action(cr.checkParams)

	up = &DockerUpdateParams{}
	update = app.Command(updateFlagKey, "Update an existing environment.")
	update.Arg(descriptorFlagKey, "The environment descriptor url (the root folder location)").Required().StringVar(&up.url)
	update.Flag(descriptorNameFlagKey, "The name of the environment descriptor, if missing we will look for a descriptor named \""+engine.DescriptorFileName+"\"").Default(engine.DescriptorFileName).StringVar(&up.file)
	update.Flag(containerOutputFlagKey, "\"true\" to write the container logs into a local file, defaulted to  \"false\"").BoolVar(&up.container.output)
	update.Flag(containerFileFlagKey, "The output file where to write the logs, if missing the log content will be written in \""+DefaultContainerLogFileName+"\"").StringVar(&up.container.file)
	update.Action(up.checkParams)

	ch = &DockerCheckParams{}
	check = app.Command(checkFlagKey, "Valid an existing environment descriptor.")
	check.Arg(descriptorFlagKey, "The environment descriptor url (the root folder location)").Required().StringVar(&ch.url)
	check.Flag(descriptorNameFlagKey, "The name of the environment descriptor, if missing we will look for a descriptor named \""+engine.DescriptorFileName+"\"").Default(engine.DescriptorFileName).StringVar(&ch.file)
	check.Flag(certPathFlagKey, "The location of the docker certificates (optional)").StringVar(&ch.cert)
	check.Flag(apiVersionFlagKey, "The version of the docker API (optional)").StringVar(&ch.api)
	check.Flag(dockerHostFlagKey, "The url of the docker host (optional)").StringVar(&ch.host)
	check.Flag(paramFileFlagKey, "The environment variables file (optional)").StringVar(&ch.container.paramFile)
	check.Flag(httpProxyFlagKey, "The http proxy(optional)").StringVar(&ch.container.httpProxy)
	check.Flag(httpsProxyFlagKey, "The https proxy (optional)").StringVar(&ch.container.httpsProxy)
	check.Flag(noProxyFlagKey, "The no proxy (optional)").StringVar(&ch.container.noProxy)
	check.Flag(containerOutputFlagKey, "\"true\" to write the container logs into a local file, defaulted to  \"false\"").BoolVar(&ch.container.output)
	check.Flag(containerFileFlagKey, "The output file where to write the logs, if missing the logcontent will be written in \""+DefaultContainerLogFileName+"\"").StringVar(&ch.container.file)
	check.Action(ch.checkParams)

	l = &Login{}
	login = app.Command(loginFlagKey, "Login into an environment manager API.")
	login.Arg(apiUrlFlagKey, "The url of the environment manager API").Required().StringVar(&l.url)
	login.Flag(userFlagKey, "The user (optional)").StringVar(&l.user)
	//login.Action(l.checkParams)

	logout = app.Command(logoutFlagKey, "Logout from an environment manager API.")

	status = app.Command(statusFlagKey, "Status of the environment manager API.")
}

func main() {
	logger = log.New(os.Stdout, "Lagoon CLI: ", log.Ldate|log.Ltime)

	fullLoginFileName = path.Join("", loginFileName)
	// this comes from http://www.kammerl.de/ascii/AsciiSignature.php
	// the font used id "standard"
	if _, err := os.Stat(fullLoginFileName); os.IsNotExist(err) {
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
		ef, e := engine.CreateExchangeFolder("out", cr.client)
		if e != nil {
			logger.Fatal(fmt.Errorf(ERROR_CREATING_EXCHANGE_FOLDER, cr.client))
		}
		ef.Create()

		log.Printf(LOG_DEPLOYING_FROM, cr.url)
		b, session := engine.HasCreationSession(*ef)
		log.Printf("Has a session for client %v", b)
		if b {
			reader := bufio.NewReader(os.Stdin)
			c := session.CreationSession.Client
			fmt.Printf(PROMPT_UPDATE_SESSION, c)
			text, _ := reader.ReadString('\n')
			if strings.TrimSpace(text) != "Y" {
				log.Printf("Cleaning the session for client %s", c)
				if err := os.Remove(session.File); err != nil {
					log.Fatal(fmt.Errorf(ERROR_CLIENT_SESSION_NOT_CLOSED, c, session.File))
				}
				ef.CleanAll()
			}
		}

		if cr.privateSSHKey != "" && cr.publicSSHKey != "" {
			// Move the ssh keys into the exchange folder input
			err := Copy(cr.publicSSHKey, filepath.Join(ef.Input.Path(), engine.SSHPuplicKeyFileName))
			if err != nil {
				logger.Fatal(fmt.Errorf(ERROR_COPYING_SSH_PUB, cr.publicSSHKey))
			}

			err = Copy(cr.privateSSHKey, filepath.Join(ef.Input.Path(), engine.SSHPrivateKeyFileName))
			if err != nil {
				logger.Fatal(fmt.Errorf(ERROR_COPYING_SSH_PRIV, cr.privateSSHKey))
			}
		}
		starterStart(*ef, cr.url, cr.file, cr.client, engine.ActionCreate, cr.container)

	}
}

// runUpdate starts the installer in order to update an environment
// The user must be logged into the environment which he wants to update
func runUpdate() {
	b, url, _ := isLogged()
	if b {
		log.Printf(LOG_UPDATING_FROM, url)
		// TODO GET REAL CLIENT NAME FROM LOGGIN
		dummyClientName := "DUMMY_CLIENT_NAME"
		ef, e := engine.CreateExchangeFolder("out", dummyClientName)
		if e != nil {
			logger.Fatal(fmt.Errorf(ERROR_CREATING_EXCHANGE_FOLDER, dummyClientName))
		}
		ef.Create()
		// TODO CALL THE API HERE IN ORDER TO START THE ENVIRONMENT UPDATE
	} else {
		log.Printf(LOG_LOGIN_REQUIRED)
	}
}

// runCheck checks the validity of the environment descriptor content
func runCheck() {
	log.Printf(LOG_CHECKING_FROM, ch.url)
	ef, e := engine.CreateExchangeFolder("out", "check")
	if e != nil {
		logger.Fatal(fmt.Errorf(ERROR_CREATING_EXCHANGE_FOLDER, cr.client))
	}
	ef.Create()
	starterStart(*ef, ch.url, ch.file, "", engine.ActionCheck, ch.container)
}

func starterStart(ef engine.ExchangeFolder, descriptor string, file string, client string, action engine.EngineAction, cp ContainerParam) {
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
	startContainer(starterImageName, done, descriptor, file, ef, client, cp, action)
	<-done
}

func getHttpProxy(param string) string {
	if param == "" {
		return os.Getenv(envHttpProxy)
	}
	return param
}

func getHttpsProxy(param string) string {
	if param == "" {
		return os.Getenv(envHttpsProxy)
	}
	return param
}

func getNoProxy(param string) string {
	if param == "" {
		return os.Getenv(envNoProxy)
	}
	return param
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

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
