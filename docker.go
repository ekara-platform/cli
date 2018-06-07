package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/mount"
	"github.com/docker/go-connections/tlsconfig"
	"golang.org/x/net/context"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/lagoon-platform/engine"
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
func startContainer(imageName string, done chan bool, descriptor string, ef engine.ExchangeFolder, client string, p ContainerParam, action engine.EngineAction) {

	envVar := []string{}
	envVar = append(envVar, engine.ClientEnvVariableKey+"="+client)
	envVar = append(envVar, engine.StarterEnvVariableKey+"="+descriptor)
	envVar = append(envVar, engine.ActionEnvVariableKey+"="+action.String())
	envVar = append(envVar, "http_proxy="+getHttpProxy(p.httpProxy))
	envVar = append(envVar, "https_proxy="+getHttpsProxy(p.httpsProxy))
	envVar = append(envVar, "no_proxy="+getNoProxy(p.noProxy))

	awsDir, err := filepath.Abs(string(path.Join(path.Dir(""), ".aws")))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Adapted output dir %s", engine.AdaptPath(awsDir))

	// **************************************************************************
	// Time added stuff - start
	// Taken from here https://blog.shameerc.com/2017/03/quick-tip-fixing-time-drift-issue-on-docker-for-mac
	// Where the working solution ws : "docker run --rm --privileged alpine hwclock -s"

	syncTime, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: "alpine",
		Cmd:   []string{"hwclock", "-s"},
	}, &container.HostConfig{
		AutoRemove: true,
		Privileged: true,
	}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(context.Background(), syncTime.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	statusCh, errCh := cli.ContainerWait(context.Background(), syncTime.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	// Time added stuff - end
	// **************************************************************************

	startedAt := time.Now().UTC()
	startedAt = startedAt.Add(time.Second * -2)
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:      imageName,
		WorkingDir: engine.InstallerVolume,
		Env:        envVar,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: ef.Location.AdaptedPath(),
				Target: engine.InstallerVolume,
			},
			{
				Type:   mount.TypeBind,
				Source: engine.AdaptPath(awsDir),
				// TODO Removed this mounting point
				Target: "/root/.aws",
			},
		},
	}, nil, "")

	if err != nil {
		panic(err)
	}

	// Chan used to turn off the rolling log
	exitCh := make(chan bool)
	if p.output {
		// Rolling output of the container logs
		go func(start time.Time, exit chan bool) {
			logMap := make(map[string]string)
			// Trick to avoid tracing twice the same log line
			notExist := func(s string) (bool, string) {
				tab := strings.Split(s, engine.InstallerLogPrefix)
				if len(tab) > 1 {
					sTrim := strings.Trim(tab[1], " ")
					if _, ok := logMap[sTrim]; ok {
						return false, ""
					}
					logMap[sTrim] = ""
					return true, engine.InstallerLogPrefix + sTrim
				} else {
					return true, s
				}
			}

			// Request to get the logs content from the container
			req := func(sr string) {
				out, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{Since: sr, ShowStdout: true, ShowStderr: true})
				if err != nil {
					panic(err)
				}
				s := bufio.NewScanner(out)
				for s.Scan() {
					str := s.Text()
					if b, sTrim := notExist(str); b {
						log.Print(sTrim)
					}
				}
				out.Close()
			}
		Loop:
			for {
				select {
				case <-exit:
					// Last call to be sure to get the end of the logs content
					now := time.Now()
					now = now.Add(time.Second * -1)
					sinceReq := strconv.FormatInt(now.Unix(), 10)
					req(sinceReq)
					break Loop
				default:
					// Running call to trace the container logs every 500ms
					sinceReq := strconv.FormatInt(start.Unix(), 10)
					start = start.Add(time.Millisecond * 500)
					req(sinceReq)
					time.Sleep(time.Millisecond * 500)
				}
			}
		}(startedAt, exitCh)
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh = cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if p.output {
			exitCh <- true
		}
		if err != nil {
			panic(err)
		}
	case <-statusCh:
		if p.output {
			exitCh <- true
		}
	}

	out, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}

	if p.output {
		logFile, err := ContainerLog(ef, p.file)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()

		_, err = io.Copy(logFile, out)
		if err != nil {
			panic(err)
		}
		log.Printf(LOG_CONTAINER_LOG_WRITTEN, logFile.Name())
	}
	done <- true
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
	if img := imageExistsByName(starterImageName); !img {
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
	done <- true
}

// Parameters required to connect with the docker API; in creation mode
type DockerCreateParams struct {
	url       string
	cert      string
	api       string
	host      string
	client    string
	container ContainerParam
}

// Parameters required to connect with the docker API; in update mode
type DockerUpdateParams struct {
	url       string
	cert      string
	api       string
	host      string
	container ContainerParam
}

// Parameters required to connect with the docker API; in Check mode
type DockerCheckParams struct {
	url       string
	cert      string
	api       string
	host      string
	container ContainerParam
}

type ContainerParam struct {
	httpProxy  string
	httpsProxy string
	noProxy    string
	output     bool
	file       string
}

// checkParams checks the coherence of the parameters received to deal with docker
// using the flags and/or the environment variables
func (n *DockerCreateParams) checkParams(c *kingpin.ParseContext) error {
	log.Printf("Creation of:%v\n", n.url)
	log.Printf("Lauched to run docker with cert:%v, api:%v, on daemon:%v\n", n.cert, n.api, n.host)
	checkDescriptor(n.url)
	// The client name is always required
	if n.client == "" {
		log.Fatal(fmt.Errorf(ERROR_REQUIRED_CLIENT))
	} else {
		log.Printf(LOG_CLIENT_CONFIRMATION, n.client)
	}

	checkDockerStuff(n.cert, n.host, n.api)
	return nil
}

// checkParams checks the coherence of the parameters received to deal with docker
// using the flags and/or the environment variables
func (n *DockerCheckParams) checkParams(c *kingpin.ParseContext) error {
	log.Printf("Creation of:%v\n", n.url)
	log.Printf("Lauched to run docker with cert:%v, api:%v, on daemon:%v\n", n.cert, n.api, n.host)
	checkDescriptor(n.url)
	checkDockerStuff(n.cert, n.host, n.api)
	return nil
}

// checkParams checks the coherence of the parameters received to deal with docker
// using the flags and/or the environment variables
func (n *DockerUpdateParams) checkParams(c *kingpin.ParseContext) error {
	log.Printf("Update of:%v\n", n.url)
	checkDescriptor(n.url)
	checkDockerStuff(n.cert, n.host, n.api)
	return nil
}

func checkDescriptor(d string) {
	// The environment descriptor is always required
	if d == "" {
		log.Fatal(fmt.Errorf(ERROR_REQUIRED_CONFIG))
	} else {
		log.Printf(LOG_CONFIG_CONFIRMATION, descriptorFlagKey, d)
	}

}

func checkDockerStuff(cert string, host string, api string) {

	// If we use flags then these 3 are required
	if host != "" || api != "" || cert != "" {
		checkFlag(cert, certPathFlagKey)
		checkFlag(host, dockerHostFlagKey)
		checkFlag(api, apiVersionFlagKey)
		log.Printf(LOG_FLAG_CONFIRMATION, certPathFlagKey, cert)
		log.Printf(LOG_FLAG_CONFIRMATION, dockerHostFlagKey, host)
		log.Printf(LOG_FLAG_CONFIRMATION, apiVersionFlagKey, api)
		log.Printf(LOG_INIT_FLAGGED_DOCKER_CLIENT)
		initFlaggedClient(host, api, cert)
	} else {
		// if the flags are not used then we will ensure
		// that the environment variables are well defined
		checkEnvVar(envCertPath)
		checkEnvVar(envDockerHost)
		log.Printf(LOG_INIT_DOCKER_CLIENT)
		initClient()
	}
}
