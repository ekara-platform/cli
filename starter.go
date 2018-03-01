package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

const (
	// Environment variables used by default by the docker client
	// "github.com/docker/docker/client"
	envCertPath   string = "DOCKER_CERT_PATH"
	envApiVersion string = "DOCKER_API_VERSION"
	envDockerHost string = "DOCKER_HOST"

	// Flags keys
	certPathFlagKey   = "cert"
	apiVersionFlagKey = "api"
	dockerHostFlagKey = "host"
	configFlagKey     = "config"

	// Name of lagoon starter image
	starterImageName string = "redis:latest"
)

var (
	CertPath              string
	APIVersion            string
	DockerHost            string
	EnvironmentDescriptor string
)

// checkEnv checks the coherence of the parameters received
// using the flags and/or the environment variables
func checkEnv() {
	// The environment descriptor is always required
	if EnvironmentDescriptor == "" {
		log.Fatal(fmt.Errorf(ERROR_REQUIRED_CONFIG))
	} else {
		log.Printf(LOG_CONFIG_CONFIRMATION, configFlagKey, EnvironmentDescriptor)
	}

	if _, err := url.ParseRequestURI(EnvironmentDescriptor); err != nil {
		log.Fatal(err.Error())
	}

	// check if we should use the flags content
	if DockerHost != "" || APIVersion != "" || CertPath != "" {
		checkFlag(CertPath, certPathFlagKey)
		checkFlag(DockerHost, dockerHostFlagKey)
		checkFlag(APIVersion, apiVersionFlagKey)
		log.Printf(LOG_FLAG_CONFIRMATION, certPathFlagKey, CertPath)
		log.Printf(LOG_FLAG_CONFIRMATION, dockerHostFlagKey, DockerHost)
		log.Printf(LOG_FLAG_CONFIRMATION, apiVersionFlagKey, APIVersion)
		log.Printf(LOG_INIT_FLAGGED_DOCKER_CLIENT)
		initFlaggedClient(DockerHost, APIVersion, CertPath)
	} else {
		// if the flags are not used then we will ensure
		// that the environment variables are well definned
		checkEnvVar(envCertPath)
		checkEnvVar(envDockerHost)
		log.Printf(LOG_INIT_DOCKER_CLIENT)
		initClient()
	}
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

func starterStart() {
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
	startContainer(starterImageName, done)
	<-done
	log.Printf(LOG_OK_STARTED)
}

func main() {

	// this comes from http://www.kammerl.de/ascii/AsciiSignature.php
	// the font used id "standard"
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

	flag.StringVar(&CertPath, certPathFlagKey, "", "The location of the docker certificates")
	flag.StringVar(&APIVersion, apiVersionFlagKey, "", "The version of the docker API")
	flag.StringVar(&DockerHost, dockerHostFlagKey, "", "The docker host")
	flag.StringVar(&EnvironmentDescriptor, configFlagKey, "", "The http location of the environment descriptor")
	flag.Parse()

	checkEnv()
	starterStart()
}
