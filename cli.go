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

	"github.com/ekara-platform/engine"
	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/util"

	"gopkg.in/alecthomas/kingpin.v2"
)

//go:generate go run generate/generate.go

const (
	ROOT_EXCHANGE_FOLDER  string = "out"
	CHECK_EXCHANGE_FOLDER string = "check"
	HEADER_PARSING_FOLDER string = "parsingGName"

	// Environment variables used by default by the docker client
	envCertPath   string = "DOCKER_CERT_PATH"
	envDockerHost string = "DOCKER_HOST"
	envHttpProxy  string = "HTTP_PROXY"
	envHttpsProxy string = "HTTPS_PROXY"
	envNoProxy    string = "NO_PROXY"

	// Flags keys for Commands
	createFlagKey  = "create"
	installFlagKey = "install"
	deployFlagKey  = "deploy"
	updateFlagKey  = "update"
	checkFlagKey   = "check"
	loginFlagKey   = "login"
	logoutFlagKey  = "logout"
	statusFlagKey  = "status"
	versionFlagKey = "version"

	// Flags keys for Arguments
	descriptorFlagKey     = "descriptor"
	descriptorNameFlagKey = "file"

	certPathFlagKey      = "cert"
	dockerHostFlagKey    = "host"
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

	// Name of the ekara starter image
	starterImageName string = "ekaraplatform/installer:latest"
)

var (
	// Commands
	create  *kingpin.CmdClause
	install *kingpin.CmdClause
	deploy  *kingpin.CmdClause
	update  *kingpin.CmdClause
	check   *kingpin.CmdClause
	login   *kingpin.CmdClause
	logout  *kingpin.CmdClause
	status  *kingpin.CmdClause
	version *kingpin.CmdClause

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
	create = app.Command(createFlagKey, "Create a new environment.")
	addFlags(create)
	addSSHFlags(create)
	create.Action(cr.checkParams)

	install = app.Command(installFlagKey, "Install a new environment.")
	addFlags(install)
	addSSHFlags(install)
	install.Action(cr.checkParams)

	deploy = app.Command(deployFlagKey, "Deploy a new environment.")
	addFlags(deploy)
	addSSHFlags(deploy)
	deploy.Action(cr.checkParams)

	up = &DockerUpdateParams{}
	update = app.Command(updateFlagKey, "Update an existing environment.")
	update.Arg(descriptorFlagKey, "The environment descriptor url (the root folder location)").Required().StringVar(&up.url)
	update.Flag(descriptorNameFlagKey, "The name of the environment descriptor, if missing we will look for a descriptor named \""+util.DescriptorFileName+"\"").Default(util.DescriptorFileName).StringVar(&up.file)
	update.Flag(containerOutputFlagKey, "\"true\" to write the container logs into a local file, defaulted to  \"false\"").BoolVar(&up.container.output)
	update.Flag(containerFileFlagKey, "The output file where to write the logs, if missing the log content will be written in \""+DefaultContainerLogFileName+"\"").StringVar(&up.container.file)
	update.Action(up.checkParams)

	ch = &DockerCheckParams{}
	check = app.Command(checkFlagKey, "Valid an existing environment descriptor.")
	addFlags(check)
	check.Action(ch.checkParams)

	l = &Login{}
	login = app.Command(loginFlagKey, "Login into an environment manager API.")
	login.Arg(apiUrlFlagKey, "The url of the environment manager API").Required().StringVar(&l.url)
	login.Flag(userFlagKey, "The user (optional)").StringVar(&l.user)
	//login.Action(l.checkParams)

	logout = app.Command(logoutFlagKey, "Logout from an environment manager API.")

	status = app.Command(statusFlagKey, "Status of the environment manager API.")

	version = app.Command(versionFlagKey, "The version details of the CLI.")
}

func addFlags(c *kingpin.CmdClause) {
	c.Arg(descriptorFlagKey, "The environment descriptor url (the root folder location)").Required().StringVar(&cr.url)
	c.Flag(descriptorNameFlagKey, "The name of the environment descriptor, if missing we will look for a descriptor named \""+util.DescriptorFileName+"\"").Default(util.DescriptorFileName).StringVar(&cr.file)
	c.Flag(certPathFlagKey, "The location of the docker certificates (optional)").StringVar(&cr.cert)
	c.Flag(dockerHostFlagKey, "The url of the docker host (optional)").StringVar(&cr.host)
	c.Flag(paramFileFlagKey, "The parameters file (optional)").StringVar(&cr.container.paramFile)
	c.Flag(httpProxyFlagKey, "The http proxy(optional)").StringVar(&cr.container.httpProxy)
	c.Flag(httpsProxyFlagKey, "The https proxy (optional)").StringVar(&cr.container.httpsProxy)
	c.Flag(noProxyFlagKey, "The no proxy (optional)").StringVar(&cr.container.noProxy)
	c.Flag(containerOutputFlagKey, "\"true\" to write the container logs into a local file, defaulted to  \"false\"").BoolVar(&cr.container.output)
	c.Flag(containerFileFlagKey, "The output file where to write the logs, if missing the log content will be written in \""+DefaultContainerLogFileName+"\"").StringVar(&cr.container.file)
}

func addSSHFlags(c *kingpin.CmdClause) {
	c.Flag(publicSSHKeyFlagKey, "The public SSH key to connect the created machines  (optional)").StringVar(&cr.publicSSHKey)
	c.Flag(privateSSHKeyFlagKey, "The private SSH key to connect the created machines  (optional)").StringVar(&cr.privateSSHKey)
}
func showHeader() {

	log.Printf("Ekara installation based on the Docker image: %s\n", starterImageName)

	fullLoginFileName = path.Join("", loginFileName)
	// this comes from http://www.kammerl.de/ascii/AsciiSignature.php
	// the font used id "standard"
	if _, err := os.Stat(fullLoginFileName); os.IsNotExist(err) {

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
}

func main() {
	logger = log.New(os.Stdout, "Ekara CLI: ", log.Ldate|log.Ltime)

	app := kingpin.New("ekara", CLI_DESCRIPTION)
	initFlags(app)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case create.FullCommand():
		showHeader()
		runCreate()
	case install.FullCommand():
		showHeader()
		runInstall()
	case deploy.FullCommand():
		showHeader()
		runDeploy()
	case update.FullCommand():
		showHeader()
		runUpdate()
	case check.FullCommand():
		showHeader()
		runCheck()
	case login.FullCommand():
		showHeader()
		runLogin()
	case logout.FullCommand():
		showHeader()
		runLogout()
	case status.FullCommand():
		showHeader()
		runStatus()
	case version.FullCommand():
		runVersion()
	}
	log.Println(LOG_COMMAND_COMPLETED)
}

// parseHeader parses the environment descriptor in order to get the qualified
// environement name
func parseHeader() string {
	ef := createEF(HEADER_PARSING_FOLDER)
	defer ef.Delete()

	p, err := ansible.ParseParams(cr.container.paramFile)
	if err != nil {
		logger.Fatalf(ERROR_UNREACHABLE_PARAM_FILE, err.Error())
	}

	logger.Printf("--> Loaded parameters %v \n", p)

	engine, err := engine.Create(logger, ef.Output.Path(), p)
	if err != nil {
		ef.Delete()
		logger.Fatalf(ERROR_CREATING_EKARA_ENGINE, err.Error())
	}

	err = engine.Init(cr.url, "", cr.file)
	if err != nil {
		ef.Delete()
		logger.Fatalf(ERROR_INITIALIZING_EKARA_ENGINE, err.Error())
	}
	qName := engine.ComponentManager().Environment().QualifiedName().String()
	logger.Printf(LOG_QUALIFIED_NAME, qName)
	return qName
}

// runCreate starts the installer in order to create an environement
func runCreate() {
	b, url, user := isLogged()
	if b {
		log.Printf(LOG_LOGGED_AS, user, url)
		log.Printf(LOG_LOGOUT_REQUIRED)
	} else {
		qName := parseHeader()
		ef := createEF(qName)

		log.Printf(LOG_CREATING_FROM, cr.url)

		if cr.privateSSHKey != "" && cr.publicSSHKey != "" {
			// Move the ssh keys into the exchange folder input
			err := Copy(cr.publicSSHKey, filepath.Join(ef.Input.Path(), util.SSHPuplicKeyFileName))
			if err != nil {
				logger.Fatal(fmt.Errorf(ERROR_COPYING_SSH_PUB, cr.publicSSHKey))
			}

			err = Copy(cr.privateSSHKey, filepath.Join(ef.Input.Path(), util.SSHPrivateKeyFileName))
			if err != nil {
				logger.Fatal(fmt.Errorf(ERROR_COPYING_SSH_PRIV, cr.privateSSHKey))
			}
		}
		starterStart(*ef, qName, cr.url, cr.file, engine.ActionCreateId, cr.container)
	}
}

// runInstall starts the installer in order to install an environement
func runInstall() {
	b, url, user := isLogged()
	if b {
		log.Printf(LOG_LOGGED_AS, user, url)
		log.Printf(LOG_LOGOUT_REQUIRED)
	} else {
		qName := parseHeader()
		ef := createEF(qName)

		log.Printf(LOG_INSTALLING_FROM, cr.url)

		if cr.privateSSHKey != "" && cr.publicSSHKey != "" {
			// Move the ssh keys into the exchange folder input
			err := Copy(cr.publicSSHKey, filepath.Join(ef.Input.Path(), util.SSHPuplicKeyFileName))
			if err != nil {
				logger.Fatal(fmt.Errorf(ERROR_COPYING_SSH_PUB, cr.publicSSHKey))
			}

			err = Copy(cr.privateSSHKey, filepath.Join(ef.Input.Path(), util.SSHPrivateKeyFileName))
			if err != nil {
				logger.Fatal(fmt.Errorf(ERROR_COPYING_SSH_PRIV, cr.privateSSHKey))
			}
		}
		starterStart(*ef, qName, cr.url, cr.file, engine.ActionInstallId, cr.container)
	}
}

// runDeploy starts the installer in order to deploy an environement
func runDeploy() {
	b, url, user := isLogged()
	if b {
		log.Printf(LOG_LOGGED_AS, user, url)
		log.Printf(LOG_LOGOUT_REQUIRED)
	} else {
		qName := parseHeader()
		ef := createEF(qName)

		log.Printf(LOG_DEPLOYING_FROM, cr.url)

		if cr.privateSSHKey != "" && cr.publicSSHKey != "" {
			// Move the ssh keys into the exchange folder input
			err := Copy(cr.publicSSHKey, filepath.Join(ef.Input.Path(), util.SSHPuplicKeyFileName))
			if err != nil {
				logger.Fatal(fmt.Errorf(ERROR_COPYING_SSH_PUB, cr.publicSSHKey))
			}

			err = Copy(cr.privateSSHKey, filepath.Join(ef.Input.Path(), util.SSHPrivateKeyFileName))
			if err != nil {
				logger.Fatal(fmt.Errorf(ERROR_COPYING_SSH_PRIV, cr.privateSSHKey))
			}
		}
		starterStart(*ef, qName, cr.url, cr.file, engine.ActionDeployId, cr.container)
	}
}

// runUpdate starts the installer in order to update an environment
// The user must be logged into the environment which he wants to update
func runUpdate() {
	b, url, _ := isLogged()
	if b {
		log.Printf(LOG_UPDATING_FROM, url)
		// TODO GET REAL QUALIFIED NAME FROM THE DESCRIPTOR
		dummyQualifiedName := "DUMMY_QUALIFIED_NAME"
		_ = createEF(dummyQualifiedName)

		// TODO CALL THE API HERE IN ORDER TO START THE ENVIRONMENT UPDATE
	} else {
		log.Printf(LOG_LOGIN_REQUIRED)
	}
}

// runCheck checks the validity of the environment descriptor content
func runCheck() {
	log.Printf(LOG_CHECKING_FROM, ch.url)
	ef := createEF(CHECK_EXCHANGE_FOLDER)
	starterStart(*ef, "check", ch.url, ch.file, engine.ActionCheckId, ch.container)
}

func starterStart(ef util.ExchangeFolder, name string, descriptor string, file string, action engine.ActionId, cp ContainerParam) {
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
	startContainer(starterImageName, done, name, descriptor, file, ef, cp, action)
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
