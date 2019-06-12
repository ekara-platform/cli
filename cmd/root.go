package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"strings"

	"github.com/ekara-platform/cli/docker"
	"github.com/ekara-platform/cli/folder"
	"github.com/ekara-platform/cli/header"
	"github.com/ekara-platform/cli/image"
	"github.com/ekara-platform/cli/message"
	"github.com/ekara-platform/engine"
	"github.com/ekara-platform/engine/ansible"
	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
	"github.com/spf13/cobra"
)

var (
	logger *log.Logger
	cr     *docker.CreateParams
)

const (
	//HeaderParsingFolder is th folder name where the header name will be parsed.
	// It's supposed to be deleted once the parsing is over.
	HeaderParsingFolder string = "parsingGName"
)

func Execute(l *log.Logger) {
	logger = l

	cr = &docker.CreateParams{}

	addCreationFlags(deployCmd)
	addSSHFlags(deployCmd)

	addLightFlags(checkCmd, dumpCmd)

	var rootCmd = &cobra.Command{}
	rootCmd.AddCommand(deployCmd, checkCmd, dumpCmd, versionCmd)
	deployCmd.AddCommand(createCmd, installCmd)
	rootCmd.Execute()
}

func addCreationFlags(cs ...*cobra.Command) {
	for _, c := range cs {
		c.PersistentFlags().StringVarP(&cr.Descriptor.File, "descriptor", "d", model.DefaultDescriptorName, "The name of the environment descriptor, if missing we will look for the defaulted name.")
		c.PersistentFlags().StringVarP(&cr.Descriptor.ParamFile, "param", "p", "", "Location of the parameters file that will be substitutable in the descriptor.")
		c.PersistentFlags().StringVarP(&cr.Descriptor.Login, "user", "U", "", "User to log into the descriptor repository.")
		c.PersistentFlags().StringVarP(&cr.Descriptor.Password, "password", "P", "", "Password to log into the descriptor repository.")
		c.PersistentFlags().StringVarP(&cr.Daemon.Cert, "cert", "c", "", "Location of the docker certificates (optional, can be substituted by an environment variable).")
		c.PersistentFlags().StringVarP(&cr.Daemon.Host, "host", "H", "", "URL of the docker host(optional, can be substituted by an environment variable).")
		c.PersistentFlags().BoolVarP(&cr.User.Output, "logs", "l", false, "Allows to turn on the installer logs.")
		c.PersistentFlags().StringVarP(&cr.User.File, "log-file", "L", "", "The output file where to write the logs, if missing the log content will be written in \""+docker.DefaultContainerLogFileName+"\".")

		c.PersistentFlags().StringVar(&cr.Installer.HttpProxy, "http-proxy", "", "The http proxy(optional).")
		c.PersistentFlags().StringVar(&cr.Installer.HttpsProxy, "https-proxy", "", "The https proxy(optional).")
		c.PersistentFlags().StringVar(&cr.Installer.NoProxy, "no-proxy", "", "The no proxy(optional).")

	}
}

func addLightFlags(cs ...*cobra.Command) {
	for _, c := range cs {
		c.PersistentFlags().StringVarP(&cr.Descriptor.File, "descriptor", "d", model.DefaultDescriptorName, "The name of the environment descriptor, if missing we will look for the defaulted name.")
		c.PersistentFlags().StringVarP(&cr.Descriptor.ParamFile, "param", "p", "", "Location of the parameters file that will be substitutable in the descriptor.")
		c.PersistentFlags().StringVarP(&cr.Descriptor.Login, "user", "U", "", "User to log into the descriptor repository.")
		c.PersistentFlags().StringVarP(&cr.Descriptor.Password, "password", "P", "", "Password to log into the descriptor repository.")
		c.PersistentFlags().StringVarP(&cr.Daemon.Cert, "cert", "c", "", "Location of the docker certificates (optional, can be substituted by an environment variable).")
		c.PersistentFlags().StringVarP(&cr.Daemon.Host, "host", "H", "", "URL of the docker host(optional, can be substituted by an environment variable).")
	}
}

func addSSHFlags(cs ...*cobra.Command) {
	for _, c := range cs {
		c.PersistentFlags().StringVar(&cr.Host.PublicSSHKey, "public-ssh", "", "Path to the public SSH key to be used for remote node access (if none given a key will be generated).")
		c.PersistentFlags().StringVar(&cr.Host.PrivateSSHKey, "private-ssh", "", "Path to the private SSH key to be used for remote node access (if none given a key will be generated).")
	}
}

func showHeader(cmd *cobra.Command, args []string) {
	header.ShowHeader()
	cr.Descriptor.Url = args[0]
	if e := cr.CheckAndLog(logger); e != nil {
		logger.Fatal(e)
	}
}

func logDone(cmd *cobra.Command, args []string) {
	logger.Println(message.LOG_COMMAND_COMPLETED)
}

func starterStart(ef util.ExchangeFolder, name string, descParam docker.DescriptorParams, action engine.ActionID, cp *docker.CreateParams) {
	logger.Printf(message.LOG_GET_IMAGE)
	done := make(chan bool, 1)
	go docker.ImagePull(image.StarterImageName, done, logger)
	<-done

	if id, running := docker.ContainerRunningByImageName(image.StarterImageName); running {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(message.PROMPT_RESTART)
		text, _ := reader.ReadString('\n')
		if strings.TrimSpace(text) == "Y" {
			done := make(chan bool, 1)
			go docker.StopContainerById(id, done, logger)
			<-done
		} else {
			logger.Printf(message.LOG_FAIL_ON_PROMPT_RESTART)
			return
		}
	}

	done = make(chan bool, 1)
	docker.StartContainer(image.StarterImageName, done, name, descParam, ef, cp, action, logger)
	<-done
}

// parseHeader parses the environment descriptor in order to get the qualified
// environement name
func parseHeader() string {
	ef := folder.CreateEF(HeaderParsingFolder, logger)
	defer ef.Delete()

	p, err := ansible.ParseParams(cr.Descriptor.ParamFile)
	if err != nil {
		logger.Fatalf(message.ERROR_UNREACHABLE_PARAM_FILE, err.Error())
	}
	vars := model.CreateContext(p)

	engine, err := engine.Create(logger, ef.Output.Path(), vars)
	if err != nil {
		ef.Delete()
		logger.Fatalf(message.ERROR_CREATING_EKARA_ENGINE, err.Error())
	}

	ctx := cliContext{
		locationContent: cr.Descriptor.Url,
		name:            cr.Descriptor.File,
		user:            cr.Descriptor.Login,
		password:        cr.Descriptor.Password,
	}

	err = engine.Init(ctx)
	if err != nil {
		ef.Delete()
		logger.Fatalf(message.ERROR_INITIALIZING_EKARA_ENGINE, err.Error())
	}
	qName := engine.ComponentManager().Environment().QualifiedName().String()
	logger.Printf(message.LOG_QUALIFIED_NAME, qName)
	return qName
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

type (
	//cliContext simulates the LaunchContext for testing purposes
	cliContext struct {
		efolder              *util.ExchangeFolder
		logger               *log.Logger
		qualifiedNameContent string
		locationContent      string
		sshPublicKeyContent  string
		sshPrivateKeyContent string
		engine               engine.Engine
		name                 string
		user                 string
		password             string
		templateContext      *model.TemplateContext
		ekaraError           error
	}
)

//Name implements the corresponding method in LaunchContext for testing purposes
func (lC cliContext) Name() string {
	return lC.name
}

//User implements the corresponding method in LaunchContext for testing purposes
func (lC cliContext) User() string {
	return lC.user
}

//Password implements the corresponding method in LaunchContext for testing purposes
func (lC cliContext) Password() string {
	return lC.password
}

//Log implements the corresponding method in LaunchContext for testing purposes
func (lC cliContext) Log() *log.Logger {
	return lC.logger
}

//Ef implements the corresponding method in LaunchContext for testing purposes
func (lC cliContext) Ef() *util.ExchangeFolder {
	return lC.efolder
}

//Ekara implements the corresponding method in LaunchContext for testing purposes
func (lC cliContext) Ekara() engine.Engine {
	return lC.engine
}

//QualifiedName implements the corresponding method in LaunchContext for testing purposes
func (lC cliContext) QualifiedName() string {
	return lC.qualifiedNameContent
}

//Location implements the corresponding method in LaunchContext for testing purposes
func (lC cliContext) Location() string {
	return lC.locationContent
}

//HTTPProxy implements the corresponding method in LaunchContext for testing purposes
func (lC cliContext) HTTPProxy() string {
	return ""
}

//HTTPSProxy implements the corresponding method in LaunchContext for testing purposes
func (lC cliContext) HTTPSProxy() string {
	return ""
}

//NoProxy implements the corresponding method in LaunchContext for testing purposes
func (lC cliContext) NoProxy() string {
	return ""
}

//SSHPublicKey implements the corresponding method in LaunchContext for testing purposes
func (lC cliContext) SSHPublicKey() string {
	return lC.sshPublicKeyContent
}

//SSHPrivateKey implements the corresponding method in LaunchContext for testing purposes
func (lC cliContext) SSHPrivateKey() string {
	return lC.sshPrivateKeyContent
}

//TemplateContext implements the corresponding method in LaunchContext for testing purposes
func (lC cliContext) TemplateContext() *model.TemplateContext {
	return lC.templateContext
}

//Error implements the corresponding method in LaunchContext for testing purposes
func (lC cliContext) Error() error {
	return lC.ekaraError
}
