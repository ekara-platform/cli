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
	HEADER_PARSING_FOLDER string = "parsingGName"
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
		c.PersistentFlags().StringVarP(&cr.Descriptor.ParamFile, "param", "p", "", " Location of the parameters file that will be substitutable in the descriptor.")
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
		c.PersistentFlags().StringVarP(&cr.Descriptor.ParamFile, "param", "p", "", " Location of the parameters file that will be substitutable in the descriptor.")
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

func starterStart(ef util.ExchangeFolder, name string, descriptor string, file string, action engine.ActionID, cp *docker.CreateParams) {
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
	docker.StartContainer(image.StarterImageName, done, name, descriptor, file, ef, cp, action, logger)
	<-done
}

// parseHeader parses the environment descriptor in order to get the qualified
// environement name
func parseHeader() string {
	ef := folder.CreateEF(HEADER_PARSING_FOLDER, logger)
	defer ef.Delete()

	p, err := ansible.ParseParams(cr.Descriptor.ParamFile)
	if err != nil {
		logger.Fatalf(message.ERROR_UNREACHABLE_PARAM_FILE, err.Error())
	}

	engine, err := engine.Create(logger, ef.Output.Path(), p)
	if err != nil {
		ef.Delete()
		logger.Fatalf(message.ERROR_CREATING_EKARA_ENGINE, err.Error())
	}

	err = engine.Init(cr.Descriptor.Url, "", cr.Descriptor.File)
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
