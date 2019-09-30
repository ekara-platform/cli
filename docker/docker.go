package docker

import (
	"bufio"
	"fmt"
	"github.com/ekara-platform/engine/action"
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

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/mount"
	"github.com/docker/go-connections/tlsconfig"

	"github.com/ekara-platform/cli/common"
	"github.com/ekara-platform/engine/util"
)

// The docker client used within the whole application
var client *docker.Client
var runningContainerId string

//EnsureDockerInit ensures that the Docker client is properly initialized
func EnsureDockerInit() {
	if client == nil {
		var err error
		var c *docker.Client
		if common.Flags.Docker.Cert != "" {
			options := tlsconfig.Options{
				CAFile:             filepath.Join(common.Flags.Docker.Cert, "ca.pem"),
				CertFile:           filepath.Join(common.Flags.Docker.Cert, "cert.pem"),
				KeyFile:            filepath.Join(common.Flags.Docker.Cert, "key.pem"),
				InsecureSkipVerify: common.Flags.Docker.TLS,
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
			c, err = docker.NewClient(common.Flags.Docker.Host, "", httpClient, nil)
		} else {
			c, err = docker.NewClient(common.Flags.Docker.Host, "", nil, nil)
		}

		if err != nil {
			panic(err)
		}
		client = c
	}
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
func StopContainerById(id string, done chan bool) {
	if err := client.ContainerStop(context.Background(), id, nil); err != nil {
		panic(err)
	}
	if err := client.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}
	for {
		common.Logger.Printf(common.LOG_WAITING_STOP)
		time.Sleep(500 * time.Millisecond)
		if stillRunning := containerRunningById(id); stillRunning == false {
			common.Logger.Printf(common.LOG_STOPPED)
			done <- true
			return
		}
	}
}

// StartContainer builds or updates a container base on the provided image name
// Once built the container will be started.
// The method will wait until the container is started and
// will notify it using the chanel
func StartContainer(url string, imageName string, done chan bool, ef util.ExchangeFolder, action action.ActionID) {
	envVar := []string{}
	envVar = append(envVar, util.StarterEnvVariableKey+"="+url)
	envVar = append(envVar, util.StarterEnvNameVariableKey+"="+common.Flags.Descriptor.File)
	envVar = append(envVar, util.StarterEnvLoginVariableKey+"="+common.Flags.Descriptor.Login)
	envVar = append(envVar, util.StarterEnvPasswordVariableKey+"="+common.Flags.Descriptor.Password)

	envVar = append(envVar, util.ActionEnvVariableKey+"="+action.String())
	envVar = append(envVar, "http_proxy="+common.Flags.Proxy.HTTP)
	envVar = append(envVar, "https_proxy="+common.Flags.Proxy.HTTPS)
	envVar = append(envVar, "no_proxy="+common.Flags.Proxy.Exclusions)

	common.Logger.Printf(common.LOG_PASSING_CONTAINER_ENVARS, envVar)

	// Check if we need to load parameters from the comand line
	fmt.Println(common.Flags.Descriptor.ParamFile)
	if common.Flags.Descriptor.ParamFile != "" {
		copyExtraParameters(common.Flags.Descriptor.ParamFile, ef)
	}

	startedAt := time.Now().UTC()
	startedAt = startedAt.Add(time.Second * -2)
	resp, err := client.ContainerCreate(context.Background(), &container.Config{
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
	runningContainerId = resp.ID

	if err != nil {
		panic(err)
	}

	// Chan used to turn off the rolling log
	exitCh := make(chan bool)

	loggerNoHearder := log.New(os.Stdout, "", 0)

	if common.Flags.Logging.ShouldOutputLogs() {
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
				out, err := client.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{Since: sr, ShowStdout: true, ShowStderr: true})
				if err != nil {
					exitCh <- true
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

	defer logAllFromContainer(resp.ID, ef, done)
	if err := client.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := client.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if common.Flags.Logging.ShouldOutputLogs() {
			exitCh <- true
		}
		panic(err)
	case <-statusCh:
		if common.Flags.Logging.ShouldOutputLogs() {
			exitCh <- true
		}
	}
}

func logAllFromContainer(id string, ef util.ExchangeFolder, done chan bool) {
	if common.Flags.Logging.ShouldOutputLogs() {
		out, err := client.ContainerLogs(context.Background(), id, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
		if err != nil {
			// we stop now
			done <- true
			return
		}

		logFile, err := containerLog(ef, common.Flags.Logging.File)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()

		_, err = io.Copy(logFile, out)
		if err != nil {
			panic(err)
		}
		common.Logger.Printf(common.LOG_CONTAINER_LOG_WRITTEN, logFile.Name())
	}

	// We are done!
	done <- true
}

// getContainers returns the detail of all running containers
func getContainers() []types.Container {
	containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{})
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
	images, err := client.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	return images
}

// ImagePull pulls the image corresponding to th given name
// and wait for the download to be completed.
//
// The completion of the download will be notified using the chanel
func ImagePull(taggedName string, done chan bool) {
	if img := imageExistsByName(taggedName); !img {
		if _, err := client.ImagePull(context.Background(), taggedName, types.ImagePullOptions{}); err != nil {
			panic(err)
		}
		for {
			common.Logger.Printf(common.LOG_WAITING_DOWNLOAD)
			time.Sleep(500 * time.Millisecond)
			if img := imageExistsByName(taggedName); img {
				common.Logger.Printf(common.LOG_DOWNLOAD_COMPLETED)
				done <- true
				return
			}
		}
	}
	done <- true
}

func copyExtraParameters(file string, ef util.ExchangeFolder) {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			common.Logger.Fatalf(common.ERROR_UNREACHABLE_PARAM_FILE, file)
		}
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = ef.Location.Write(b, util.ExternalVarsFilename)
	if err != nil {
		panic(err)
	}
}

func containerLog(ef util.ExchangeFolder, fileName string) (*os.File, error) {
	f, e := os.Create(filepath.Join(ef.Output.Path(), common.Flags.Logging.File))
	if e != nil {
		return nil, e
	}
	return f, nil
}
