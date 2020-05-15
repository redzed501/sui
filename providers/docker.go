package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/willfantom/sui/config"
)

const (
	dockerAPIVersion string = "v1.24"
	nameFromLabel    string = "sui.name"
	iconFromLabel    string = "sui.icon"
	urlFromLabel     string = "sui.url"
	protecFromLabel  string = "sui.protect"
)

type DockerProvider struct {
	Client *http.Client
	User   string
	Pass   string
}

type ContainerInfo struct {
	ID     string            `json:"Id"`
	Names  []string          `json:"Names"`
	Image  string            `json:"Image"`
	State  string            `json:"State"`
	Labels map[string]string `json:"Labels"`
}

type DockerVersionInfo struct {
	Version string `json:"Version"`
	Os      string `json:"Os"`
}

func NewDockerProvider(cnf *config.DockerConfig) (*AppProvider, error) {
	ap := newAppProvider()
	ap.PType = Docker

	var dp DockerProvider
	client, err := createDockerClient(cnf.Path, cnf.DType)
	if err != nil {
		return nil, fmt.Errorf("could not create docker app provider")
	}
	dp.Client = client
	dp.User = cnf.User
	dp.Pass = cnf.Pass

	ap.TypeConfig = &dp

	if !dp.TestDockerConn() {
		return nil, fmt.Errorf("could not create docker app provider")
	}
	return ap, nil
}

func NewDockerProviderLite(cnf *config.DockerConfig) (*DockerProvider, error) {
	var dp DockerProvider
	client, err := createDockerClient(cnf.Path, cnf.DType)
	if err != nil {
		return nil, fmt.Errorf("could not create docker provider")
	}
	dp.Client = client
	dp.User = cnf.User
	dp.Pass = cnf.Pass
	if !dp.TestDockerConn() {
		return nil, fmt.Errorf("could not create docker app provider")
	}
	return &dp, nil
}

func (dp *DockerProvider) GetApps(list map[string]*App) error {

	containers, err := dp.GetContainerList()
	if err != nil {
		return err
	}
	for _, container := range containers {

		app := newApp()

		var name string
		if len(container.Names) > 0 {
			name = container.Names[0][1:]
		}
		labelName, exist := container.Labels[nameFromLabel]
		if exist {
			name = labelName
		}
		labelIcon, exist := container.Labels[iconFromLabel]
		if exist {
			app.Icon = labelIcon
		} else {
			app.Icon = getDefaultIcon(name)
		}
		labelUrl, exist := container.Labels[urlFromLabel]
		if exist {
			app.URL = labelUrl
		}
		labelProtec, exist := container.Labels[protecFromLabel]
		if exist {
			labelProtecBool, err := strconv.ParseBool(labelProtec)
			if err != nil {
				log.Errorf("Provided 'protect' (%s) value is not a bool!", labelProtec)
				continue
			}
			app.Protected = labelProtecBool
		}
		list[name] = app
	}

	return nil
}

func (app *App) UpdateFromDockerLabels(ci *ContainerInfo) bool {
	lIcon, icex := ci.Labels[iconFromLabel]
	lURL, urlex := ci.Labels[urlFromLabel]
	lProc, procex := ci.Labels[protecFromLabel]

	if icex {
		app.Icon = lIcon
	}
	if urlex {
		app.URL = lURL
	}
	if procex {
		lProcB, err := strconv.ParseBool(lProc)
		if err == nil {
			app.Protected = lProcB
		} else {
			procex = false
		}
	}
	if icex || urlex || procex {
		return true
	}
	return false
}

func (dp *DockerProvider) TestDockerConn() bool {
	version, err := dp.GetDockerVersion()
	if err != nil {
		return false
	}
	log.Debugf("docker version found: %s", version.Version)
	return true
}

func (dp *DockerProvider) GetContainerList() ([]*ContainerInfo, error) {
	response, err := requestFromDocker(dp.Client, "containers/json")
	if err != nil || response.StatusCode != 200 {
		log.Errorf("failed to fetch local docker container list")
		return nil, err
	}
	var containerList []*ContainerInfo
	err = json.NewDecoder(response.Body).Decode(&containerList)
	if err != nil {
		log.Errorf("docker container list could not be parsed")
		return nil, err
	}
	return containerList, nil
}

func (dp *DockerProvider) GetContainerInfo(name string) (*ContainerInfo, error) {
	response, err := requestFromDocker(dp.Client, fmt.Sprintf("containers/%s/json", name))
	if err != nil || response.StatusCode != 200 {
		log.Errorf("failed to fetch local docker container info")
		return nil, err
	}
	var containerInfo *ContainerInfo
	err = json.NewDecoder(response.Body).Decode(&containerInfo)
	if err != nil {
		log.Errorf("docker container info could not be parsed")
		return nil, err
	}
	return containerInfo, nil
}

func (dp *DockerProvider) GetDockerVersion() (*DockerVersionInfo, error) {
	response, err := requestFromDocker(dp.Client, "version")
	if err != nil || response.StatusCode != 200 {
		log.Errorf("failed to fetch local docker version")
		return nil, err
	}
	var versionInfo *DockerVersionInfo
	err = json.NewDecoder(response.Body).Decode(&versionInfo)
	if err != nil {
		log.Errorf("docker version info could not be parsed\n")
		return nil, err
	}
	return versionInfo, nil
}

func createDockerClient(path string, dType config.DockerType) (*http.Client, error) {
	typeString := "tcp"
	if dType == config.Socket {
		if _, err := os.Stat(path); err == nil {
			log.Debugf("docker socket path found :)")
			typeString = "unix"
		} else if os.IsNotExist(err) {
			log.Errorf("you must mount the docker socket (%s)\n", path)
			return nil, fmt.Errorf("docker socket is not mounted (correctly)")
		} else {
			log.Errorf("can not find the docker socket (%s)\n", path)
			return nil, fmt.Errorf("docker socket not found")
		}
	}
	httpClient := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial(typeString, path)
			},
		},
	}
	return &httpClient, nil
}

func requestFromDocker(c *http.Client, path string) (*http.Response, error) {
	//TODO: Convert Get to Do to add Basic Auth
	//TODO: Allow non Localhost host to be used
	path = fmt.Sprintf("http://127.0.0.1/%s/%s", dockerAPIVersion, path)
	response, err := c.Get(path)
	return response, err
}
