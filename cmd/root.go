package cmd

import (
	"fmt"
	"github.com/ekara-platform/cli/docker"
	"log"
	"os"

	"github.com/ekara-platform/cli/common"
	"github.com/ekara-platform/engine"
	"github.com/ekara-platform/engine/util"
	"github.com/ekara-platform/model"
	"github.com/spf13/cobra"
)

const (
	envHTTPProxy       string = "http_proxy"
	envHTTPSProxy      string = "https_proxy"
	envNoProxy         string = "no_proxy"
	defaultLogFileName string = "installer.log"
	defaultVarFileName string = "vars.yaml"
	rootExchangeFolder string = "ekara"
)

var rootCmd = &cobra.Command{
	Use:   "ekara",
	Short: "Ekara is a lightweight platform for deploying cloud applications.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		info, e := os.Stdout.Stat()
		if e != nil {
			return
		} else if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
			if common.Flags.Logging.ShouldOutputLogs() {
				common.Logger = log.New(os.Stdout, "CLI  > ", log.Ldate|log.Ltime)
			}

			// this comes from http://www.kammerl.de/ascii/AsciiSignature.php
			// the font used id "standard"
			fmt.Println(" _____ _                   ")
			fmt.Println("| ____| | ____ _ _ __ __ _ ")
			fmt.Println("|  _| | |/ / _` | '__/ _` |")
			fmt.Println("| |___|   < (_| | | | (_| |")
			fmt.Println(`|_____|_|\_\__,_|_|  \__,_|`)
			if isDescriptorCommand(cmd, args) {
				fmt.Println(args[0])
			}
			fmt.Println("")
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Done")
	},
}

func isDescriptorCommand(cmd *cobra.Command, args []string) bool {
	return len(args) > 0 && (cmd.Name() == "apply" || cmd.Name() == "dump" || cmd.Name() == "validate")
}

func init() {
	// Proxy flags
	rootCmd.PersistentFlags().StringVar(&common.Flags.Proxy.HTTP, "http-proxy", os.Getenv(envHTTPProxy), "HTTP proxy url")
	rootCmd.PersistentFlags().StringVar(&common.Flags.Proxy.HTTPS, "https-proxy", os.Getenv(envHTTPSProxy), "HTTPS proxy url")
	rootCmd.PersistentFlags().StringVar(&common.Flags.Proxy.Exclusions, "no-proxy", os.Getenv(envNoProxy), "Proxy exclusion(s)")

	// Logging flags
	rootCmd.PersistentFlags().BoolVar(&common.Flags.Logging.Verbose, "verbose", false, "Verbose standard output")
	rootCmd.PersistentFlags().BoolVar(&common.Flags.Logging.VeryVerbose, "very-verbose", false, "Very verbose standard output")
	rootCmd.PersistentFlags().StringVar(&common.Flags.Logging.File, "logfile", defaultLogFileName, "Installer logfile")

	// Debug flags
	rootCmd.PersistentFlags().BoolVar(&common.Flags.Debug, "debug", false, "Installer logfile")
}

// Execute launchs the adequate command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func StopCurrentContainerIfRunning() {
	if id, running := docker.ContainerRunningByImageName(starterImageName); running {
		done := make(chan bool, 1)
		go docker.StopContainerById(id, done)
		<-done
	}
}

func applyDescriptorFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&common.Flags.Descriptor.File, "descriptor", "d", model.DefaultDescriptorName, "Name of the main environment descriptor")
	cmd.PersistentFlags().StringVarP(&common.Flags.Descriptor.ParamFile, "vars", "v", checkDefaultVarFile(), "Path to the external variable file")
	cmd.PersistentFlags().StringVarP(&common.Flags.Descriptor.Login, "user", "u", "", "Username for the main descriptor repository")
	cmd.PersistentFlags().StringVarP(&common.Flags.Descriptor.Password, "password", "p", "", "Password for the main descriptor repository")
}

func checkDefaultVarFile() string {
	if _, err := os.Stat(defaultVarFileName); !os.IsNotExist(err) {
		return defaultVarFileName
	}
	return ""
}

func createEF(folder string) util.ExchangeFolder {
	ef, e := util.CreateExchangeFolder(folder, "")
	if e != nil {
		common.Logger.Fatal(fmt.Errorf(common.ERROR_CREATING_EXCHANGE_FOLDER, folder))
	}
	e = ef.Create()
	if e != nil {
		common.Logger.Fatal(fmt.Errorf(common.ERROR_CREATING_EXCHANGE_FOLDER, e.Error()))
	}
	return ef
}

func initLocalEngine(workDir string, descriptorURL string) engine.Ekara {
	var p model.Parameters
	if common.Flags.Descriptor.ParamFile != "" {
		var err error
		p, err = model.ParseParameters(common.Flags.Descriptor.ParamFile)
		if err != nil {
			common.Logger.Fatalf(common.ERROR_UNREACHABLE_PARAM_FILE, err.Error())
		}
	} else {
		p = model.Parameters{}
	}

	e, err := engine.Create(&cliContext{
		ef:             createEF(rootExchangeFolder),
		logger:         common.Logger,
		location:       descriptorURL,
		descriptorName: common.Flags.Descriptor.File,
		user:           common.Flags.Descriptor.Login,
		password:       common.Flags.Descriptor.Password,
		extVars:        p,
	}, workDir)
	if err != nil {
		common.Logger.Fatalf(common.ERROR_CREATING_EKARA_ENGINE, err.Error())
	}

	err = e.Init()
	if err != nil {
		common.Logger.Fatalf(common.ERROR_INITIALIZING_EKARA_ENGINE, err.Error())
	}

	return e
}
