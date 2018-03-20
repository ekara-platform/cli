package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"github.com/docker/go-connections/tlsconfig"
	"golang.org/x/net/context"

	"gopkg.in/alecthomas/kingpin.v2"
)

// The docker client used within the whole application
var cli docker.Client

// initFlaggedClient initializes the docker client using the flaged values
func initFlaggedClient(host string, api string, path string) {

	options := tlsconfig.Options{
		CAFile:             filepath.Join(path, "ca.pem"),
		CertFile:           filepath.Join(path, "cert.pem"),
		KeyFile:            filepath.Join(path, "key.pem"),
		InsecureSkipVerify: false,
	}

	tlsc, err := tlsconfig.Client(options)
	if err != nil {
		panic(err)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsc,
		},
		CheckRedirect: docker.CheckRedirect,
	}

	c, err := docker.NewClient(host, api, httpClient, nil)
	if err != nil {
		panic(err)
	}
	cli = *c
}

// initClient initializes the docker client using the environment variables
func initClient() {
	c, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}
	cli = *c
}

// containerRunningByImageName returns true if a container, built
// on the given image, is running
func containerRunningByImageName(name string) (string, bool) {
	containers := getContainers()
	for _, container := range containers {
		if container.Image == name || container.Image+":latest" == name {
			return container.ID, true
		}
	}
	return "", false
}

//containerRunningById returns true if a container with the given id is running
func containerRunningById(id string) bool {
	containers := getContainers()
	for _, container := range containers {
		if container.ID == id {
			return true
		}
	}
	return false
}

//stopContainerById stops a container corresponding to the provider id
func stopContainerById(id string, done chan bool) {
	if err := cli.ContainerStop(context.Background(), id, nil); err != nil {
		panic(err)
	}
	if err := cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}
	for {
		log.Printf(LOG_WAITING_STOP)
		time.Sleep(500 * time.Millisecond)
		if stillRunning := containerRunningById(id); stillRunning == false {
			log.Printf(LOG_STOPPED)
			done <- true
			return
		}
	}
}

// startContainer builds or updates a container base on the provided image name
// Once built the container will be started.
// The method will wait until the container is started and
// will notify it using the chanel
func startContainer(imageName string, done chan bool, create bool, descriptor []byte) {

	if create {
		log.Printf(LOG_START_CREATION)
	} else {
		log.Printf(LOG_START_UPDATE)
	}
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: imageName,
		Env:   []string{starterEnvVariableKey + "=" + string(descriptor)},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	for {
		log.Printf(LOG_WAITING_START)
		time.Sleep(500 * time.Millisecond)
		if _, isRunning := containerRunningByImageName(imageName); isRunning {
			log.Printf(LOG_STARTED)
			done <- true
			return
		}
	}
}

// getContainers returns the detail of all running containers
func getContainers() []types.Container {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	return containers
}

// imageExistsByName returns true if an image corresponding
// to the given name has been already downloaded
func imageExistsByName(name string) bool {
	images := getImages()
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == name {
				return true
			}
		}
	}
	return false
}

// getImages returns the summary of all images already downloaded
func getImages() []types.ImageSummary {
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	return images
}

// imagePull pulls the image corresponding to th given name
// and wait for the download to be completed.
//
// The completion of the download will be notified using the chanel
func imagePull(taggedName string, done chan bool) {
	if _, err := cli.ImagePull(context.Background(), taggedName, types.ImagePullOptions{}); err != nil {
		panic(err)
	}
	for {
		log.Printf(LOG_WAITING_DOWNLOAD)
		time.Sleep(500 * time.Millisecond)
		if img := imageExistsByName(starterImageName); img {
			log.Printf(LOG_DOWNLOAD_COMPLETED)
			done <- true
			return
		}
	}
}

// Parameters required to connect with the docker API
type DockerParams struct {
	url  string
	cert string
	api  string
	host string
}

// Parameters required to check the environment descriptor contant
type CheckParams struct {
	url    string
	file   string
	output bool
}

// checkDockerParams checks the coherence of the parameters received to deal with docker
// using the flags and/or the environment variables
func (n *DockerParams) checkDockerParams(c *kingpin.ParseContext) error {
	log.Printf("Create or update of:%v\n", n.url)
	log.Printf("Lauched to run docker with cert:%v, api:%v, on daemon:%v\n", n.cert, n.api, n.host)

	// The environment descriptor is always required
	if n.url == "" {
		log.Fatal(fmt.Errorf(ERROR_REQUIRED_CONFIG))
	} else {
		log.Printf(LOG_CONFIG_CONFIRMATION, descriptorFlagKey, n.url)
	}

	// If we use flags then these 3 are required
	if n.host != "" || n.api != "" || n.cert != "" {
		checkFlag(n.cert, certPathFlagKey)
		checkFlag(n.host, dockerHostFlagKey)
		checkFlag(n.api, apiVersionFlagKey)
		log.Printf(LOG_FLAG_CONFIRMATION, certPathFlagKey, n.cert)
		log.Printf(LOG_FLAG_CONFIRMATION, dockerHostFlagKey, n.host)
		log.Printf(LOG_FLAG_CONFIRMATION, apiVersionFlagKey, n.api)
		log.Printf(LOG_INIT_FLAGGED_DOCKER_CLIENT)
		initFlaggedClient(n.host, n.api, n.cert)
	} else {
		// if the flags are not used then we will ensure
		// that the environment variables are well defined
		checkEnvVar(envCertPath)
		checkEnvVar(envDockerHost)
		log.Printf(LOG_INIT_DOCKER_CLIENT)
		initClient()
	}
	return nil
}
