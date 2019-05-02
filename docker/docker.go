package docker

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	_ "gopkg.in/alecthomas/kingpin.v2"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/mount"
	"github.com/docker/go-connections/tlsconfig"

	"github.com/ekara-platform/cli/message"
	"github.com/ekara-platform/engine"
	"github.com/ekara-platform/engine/util"
)

// The docker client used within the whole application
var cli docker.Client

const (
	DefaultContainerLogFileName string = "installer.log"
	envHttpProxy                string = "HTTP_PROXY"
	envHttpsProxy               string = "HTTPS_PROXY"
	envNoProxy                  string = "NO_PROXY"
)

type CreateParams struct {
	Daemon     DaemonParams
	Descriptor DescriptorParams
	Installer  InstallerParams
	Host       HostParams
	User       UserParams
}

func (c CreateParams) CheckAndLog(logger *log.Logger) error {
	if e := c.Daemon.checkAndLog(logger); e != nil {
		return e
	}
	if e := c.Descriptor.checkAndLog(logger); e != nil {
		return e
	}
	if e := c.Installer.checkAndLog(logger); e != nil {
		return e
	}
	if e := c.Host.checkAndLog(logger); e != nil {
		return e
	}
	if e := c.User.checkAndLog(logger); e != nil {
		return e
	}
	return nil
}
func checkEnvVar(key string) error {
	if os.Getenv(key) == "" {

		return fmt.Errorf(message.ERROR_REQUIRED_ENV, key)
	}
	return nil
}

func checkFlag(val string, flagKey string) error {
	if val == "" {
		return fmt.Errorf(message.ERROR_REQUIRED_FLAG, flagKey)
	}
	return nil
}

// initClient initializes the docker client using the environment variables
func initClient() {
	c, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}
	cli = *c
}

// initFlaggedClient initializes the docker client using the flaged values
func initFlaggedClient(host string, cert string) {

	var err error
	var c *docker.Client
	if cert != "" {
		options := tlsconfig.Options{
			CAFile:             filepath.Join(cert, "ca.pem"),
			CertFile:           filepath.Join(cert, "cert.pem"),
			KeyFile:            filepath.Join(cert, "key.pem"),
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
		c, err = docker.NewClient(host, "", httpClient, nil)
	} else {
		c, err = docker.NewClient(host, "", nil, nil)
	}

	if err != nil {
		panic(err)
	}
	cli = *c
}

// ContainerRunningByImageName returns true if a container, built
// on the given image, is running
func ContainerRunningByImageName(name string) (string, bool) {
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
func StopContainerById(id string, done chan bool, logger *log.Logger) {
	if err := cli.ContainerStop(context.Background(), id, nil); err != nil {
		panic(err)
	}
	if err := cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}
	for {
		logger.Printf(message.LOG_WAITING_STOP)
		time.Sleep(500 * time.Millisecond)
		if stillRunning := containerRunningById(id); stillRunning == false {
			logger.Printf(message.LOG_STOPPED)
			done <- true
			return
		}
	}
}

// startContainer builds or updates a container base on the provided image name
// Once built the container will be started.
// The method will wait until the container is started and
// will notify it using the chanel
func StartContainer(imageName string, done chan bool, name string, descriptor string, file string, ef util.ExchangeFolder, p *CreateParams, action engine.ActionID, logger *log.Logger) {

	envVar := []string{}
	envVar = append(envVar, util.StarterEnvVariableKey+"="+descriptor)
	envVar = append(envVar, util.StarterEnvNameVariableKey+"="+file)
	envVar = append(envVar, util.StarterEnvQualifiedVariableKey+"="+name)

	envVar = append(envVar, util.ActionEnvVariableKey+"="+action.String())
	envVar = append(envVar, "http_proxy="+getHttpProxy(p.Installer.HttpProxy, logger))
	envVar = append(envVar, "https_proxy="+getHttpsProxy(p.Installer.HttpsProxy, logger))
	envVar = append(envVar, "no_proxy="+getNoProxy(p.Installer.NoProxy, logger))

	logger.Printf(message.LOG_PASSING_CONTAINER_ENVARS, envVar)

	// Check if we need to load parameters from the comand line
	if p.Descriptor.ParamFile != "" {
		copyExtraParameters(p.Descriptor.ParamFile, ef, logger)
	}

	startedAt := time.Now().UTC()
	startedAt = startedAt.Add(time.Second * -2)
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:      imageName,
		WorkingDir: util.InstallerVolume,
		Env:        envVar,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: ef.Location.AdaptedPath(),
				Target: util.InstallerVolume,
			},
			{
				Type:   mount.TypeBind,
				Source: "/var/run/docker.sock",
				Target: "/var/run/docker.sock",
			},
		},
	}, nil, "")

	if err != nil {
		panic(err)
	}

	// Chan used to turn off the rolling log
	exitCh := make(chan bool)

	loggerNoHearder := log.New(os.Stdout, "", 0)

	if p.User.Output {
		// Rolling output of the container logs
		go func(start time.Time, exit chan bool) {
			logMap := make(map[string]string)
			// Trick to avoid tracing twice the same log line
			notExist := func(s string) (bool, string) {
				tab := strings.Split(s, util.InstallerLogPrefix)
				if len(tab) > 1 {
					sTrim := strings.Trim(tab[1], " ")
					if _, ok := logMap[sTrim]; ok {
						return false, ""
					}
					logMap[sTrim] = ""
					return true, util.InstallerLogPrefix + sTrim
				} else {
					return true, s
				}
			}

			// Request to get the logs content from the container
			req := func(sr string) {
				out, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{Since: sr, ShowStdout: true, ShowStderr: true})
				if err != nil {
					// TODO REMOVE PANIC
					panic(err)
				}
				s := bufio.NewScanner(out)
				for s.Scan() {
					str := s.Text()
					if b, sTrim := notExist(str); b {
						loggerNoHearder.Print(sTrim)
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
	defer logAllFromContainer(resp.ID, ef, done, p, logger)
	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if p.User.Output {
			exitCh <- true
		}
		if err != nil {
			panic(err)
		}
	case <-statusCh:
		if p.User.Output {
			exitCh <- true
		}
	}
}

func logAllFromContainer(id string, ef util.ExchangeFolder, done chan bool, p *CreateParams, logger *log.Logger) {
	if p.User.Output {
		out, err := cli.ContainerLogs(context.Background(), id, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
		if err != nil {
			panic(err)
		}

		logFile, err := ContainerLog(ef, p.User.File)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()

		_, err = io.Copy(logFile, out)
		if err != nil {
			panic(err)
		}
		logger.Printf(message.LOG_CONTAINER_LOG_WRITTEN, logFile.Name())
	}
	// We are done!
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

// ImagePull pulls the image corresponding to th given name
// and wait for the download to be completed.
//
// The completion of the download will be notified using the chanel
func ImagePull(taggedName string, done chan bool, logger *log.Logger) {
	if img := imageExistsByName(taggedName); !img {
		if _, err := cli.ImagePull(context.Background(), taggedName, types.ImagePullOptions{}); err != nil {
			panic(err)
		}
		for {
			logger.Printf(message.LOG_WAITING_DOWNLOAD)
			time.Sleep(500 * time.Millisecond)
			if img := imageExistsByName(taggedName); img {
				logger.Printf(message.LOG_DOWNLOAD_COMPLETED)
				done <- true
				return
			}
		}
	}
	done <- true
}

func getHttpProxy(param string, logger *log.Logger) string {
	if param == "" {
		s := os.Getenv(envHttpsProxy)
		logger.Printf(message.LOG_GETTING_HTTP_PROXY, s)
		return s
	}
	return param
}

func getHttpsProxy(param string, logger *log.Logger) string {
	if param == "" {
		s := os.Getenv(envHttpsProxy)
		logger.Printf(message.LOG_GETTING_HTTPS_PROXY, s)
		return s
	}
	return param
}

func getNoProxy(param string, logger *log.Logger) string {
	if param == "" {
		s := os.Getenv(envNoProxy)
		logger.Printf(message.LOG_GETTING_NO_PROXY, s)
		return s
	}
	return param
}

func copyExtraParameters(file string, ef util.ExchangeFolder, logger *log.Logger) {
	// Check if the parameter file exist
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			logger.Fatalf(message.ERROR_UNREACHABLE_PARAM_FILE, file)
		}
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = ef.Location.Write(b, util.CliParametersFileName)
	if err != nil {
		panic(err)
	}
}

func ContainerLog(ef util.ExchangeFolder, fileName string) (*os.File, error) {
	var file string
	if fileName == "" {
		file = filepath.Join(ef.Output.Path(), DefaultContainerLogFileName)
	} else {
		file = filepath.Join(ef.Output.Path(), fileName)
	}
	f, e := os.Create(file)
	if e != nil {
		return nil, e
	}
	return f, nil
}
