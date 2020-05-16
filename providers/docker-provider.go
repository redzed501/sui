package providers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

////----- Constants --->
const (
	dockerAPIVersion string = "v1.24"
)

var (
	dockerEnabledLabel string = fmt.Sprintf("%s.enabled", labelRoot)
	dockerNameLabel    string = fmt.Sprintf("%s.name", labelRoot)
	dockerURLLabel     string = fmt.Sprintf("%s.url", labelRoot)
	dockerIconLabel    string = fmt.Sprintf("%s.icon", labelRoot)
)

////----- end

////----- Models --->
type Docker struct {
	Host   string
	Client *http.Client
	User   string
	Pass   string
}
type DockerConfig struct {
	ConnType string `json:"connection"`
	ConnPath string `json:"path"`
	ConnURL  string `json:"url"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
}
type DockerContainerInfo struct {
	ID    string   `json:"Id"`
	Names []string `json:"Names"`
	Image string   `json:"Image"`
	//State  string            `json:"State"`
	Labels map[string]string `json:"Labels"`
}
type DockerContainerConfig struct {
	Labels map[string]string `json:"Labels"`
}
type DockerIndividualInfo struct {
	ID     string                `json:"Id"`
	Name   string                `json:"Names"`
	Config DockerContainerConfig `json:"Config"`
}
type DockerVersionInfo struct {
	Version string `json:"Version"`
	Os      string `json:"Os"`
}

////----- end

////----- Common App Provider Functions --->
func newDocker(name string) (*Docker, error) {
	log.Debugf("creating new docker provider")
	config, err := loadDockerConfig(name)
	if err != nil {
		return nil, err
	}
	docker := Docker{
		Host: config.ConnURL,
		User: config.User,
		Pass: config.Pass,
	}
	docker.Client, err = config.createClient()
	if err != nil {
		return nil, err
	}
	if !docker.TestConnection(true) {
		return nil, fmt.Errorf("could not create docker connection | %s", name)
	}
	return &docker, nil
}

func toDocker(cnf interface{}) (*Docker, error) {
	var docker *Docker
	docker, valid := cnf.(*Docker)
	if !valid {
		return nil, fmt.Errorf("conversion to docker type not valid")
	}
	return docker, nil
}

// GetApps fetches all the containers visible fromthe provided Docker client
// It will however not provide any apps back where the enabled flag was false
// Returns a map of App formatted containers
func (dkr *Docker) GetApps() map[string]*App {

	containers, err := dkr.getContainerList()
	if err != nil {
		return nil
	}

	apps := make(map[string]*App)
	for _, container := range containers {
		app := newApp()
		var name string
		if len(container.Names) > 0 {
			name = container.Names[0][1:]
		}
		newName, upName := dkr.UpgradeApp(name, app)
		if upName {
			name = newName
		}
		if app.Enabled {
			apps[name] = app
		}
	}

	return apps
}

// UpgradeApp takes an already existing app and replaces data with defined data
// an example use case is where you want to overwrite an apps info with the info
// from docker labels
// In this case (docker), match name should be the name of the container you want
// to use
// Returns true if returning a new suggested app name
func (dkr *Docker) UpgradeApp(matchName string, app *App) (string, bool) {
	info, err := dkr.getContainerInfo(matchName)
	if err != nil {
		return "", false
	}
	if info != nil {
		lIcon, icex := info.Config.Labels[dockerIconLabel]
		lURL, urlex := info.Config.Labels[dockerURLLabel]
		lEnab, enabex := info.Config.Labels[dockerEnabledLabel]

		if icex {
			app.Icon = lIcon
		}
		if urlex {
			app.URL = lURL
		}
		if enabex {
			lEnabB, err := strconv.ParseBool(lEnab)
			if err == nil {
				app.Enabled = lEnabB
			} else {
				enabex = false
			}
		}
		lName, namex := info.Config.Labels[dockerNameLabel]
		if namex {
			return lName, true
		}
	}
	return "", false
}

// TestConnection checks to see if the docker client can be communicated with via
// the given http client. This is done by checking the version
// Returns true if communication was possible
func (dkr *Docker) TestConnection(output bool) bool {
	version, err := dkr.getDockerVersion()
	if err != nil {
		return false
	}
	log.Debugf("docker version checked | %s", version.Version)
	return true
}

func loadDockerConfig(name string) (*DockerConfig, error) {
	configPath := fmt.Sprintf("%s/%s.json", getFileConfigRoot(), name)
	dkrCnf := newDockerConfig()
	cnfFile, err := os.Open(configPath)
	if err != nil {
		log.Errorf("no config file found | %s", configPath)
		return nil, err
	}
	defer cnfFile.Close()
	err = json.NewDecoder(cnfFile).Decode(dkrCnf)
	if err != nil {
		log.Errorf("config file could not be parsed | %s", configPath)
		return nil, err
	}
	return dkrCnf, nil
}

////----- end

////----- Provider Specific ---> eof
func newDockerConfig() *DockerConfig {
	return &DockerConfig{
		ConnPath: "/var/run/docker.sock",
		ConnType: "socker",
		ConnURL:  "",
		User:     "",
		Pass:     "",
	}
}

func (cnf *DockerConfig) createClient() (*docker.Client, error) {
	var path string
	if strings.ToLower(cnf.ConnType) == "unix" {
		if !cnf.dockerSocketExist() {
			log.Errorf("you must mount the docker socket (%s) to use unix type", path)
			return nil, fmt.Errorf("docker socket is not mounted (correctly)")
		}
		path = fmt.Sprintf("%s://%s", strings.ToLower(cnf.ConnType), cnf.ConnPath)
	} else if strings.ToLower(cnf.ConnType) == "tcp" {
		if !cnf.dockerTCPOkay() {
			log.Errorf("you must enter a valid [ip]:[port] to use tcp type docker")
			return nil, fmt.Errorf("docker host is not valid")
		}
		path = fmt.Sprintf("%s://%s", strings.ToLower(cnf.ConnType), cnf.ConnURL)
	} else {
		return nil, fmt.Errorf("type must be unix or tcp | given -> %s", cnf.ConnType)
	}
	dkrClient, err := docker.NewClient(path)
	if err != nil {
		return nil, err
	}
	return dkrClient, nil
}

func (cnf *DockerConfig) dockerSocketExist() bool {
	if _, err := os.Stat(cnf.ConnPath); err == nil {
		return true
	}
	return false
}

func (cnf *DockerConfig) dockerTCPOkay() bool {
	parts := strings.Split(cnf.ConnURL, ":")
	log.Infoln(parts)
	if len(parts) != 2 {
		return false
	}
	ip := net.ParseIP(parts[0])
	if ip != nil {
		_, err := strconv.ParseInt(parts[1], 10, 32)
		if err == nil {
			return true
		}
	}
	return false
}

func (dkr *Docker) requestFromDocker(path string) (*http.Response, error) {
	//TODO: support non localhost hosts....
	host := dkr.Host
	prefix := "tcp"
	if host == "" {
		host = "127.0.0.1"
		prefix = "http"
	}
	path = fmt.Sprintf("%s://%s/%s/%s", prefix, host, dockerAPIVersion, path)
	req, err := http.NewRequest("GET", path, nil)
	if dkr.User != "" {
		req.SetBasicAuth(dkr.User, dkr.Pass)
	}
	response, err := dkr.Client.Do(req)
	if err != nil {
		panic(err)
	}
	return response, err
}

func (dkr *Docker) getContainerList() ([]*DockerContainerInfo, error) {
	response, err := dkr.requestFromDocker("containers/json")
	if err != nil || response.StatusCode != 200 {
		log.Errorf("failed to fetch docker container list")
		return nil, err
	}
	var containerList []*DockerContainerInfo
	err = json.NewDecoder(response.Body).Decode(&containerList)
	if err != nil {
		log.Errorf("docker container list could not be parsed")
		return nil, err
	}
	return containerList, nil
}

func (dkr *Docker) getContainerInfo(name string) (*DockerIndividualInfo, error) {
	response, err := dkr.requestFromDocker(fmt.Sprintf("containers/%s/json", name))
	if err != nil || response.StatusCode != 200 {
		//log.Errorf("failed to fetch docker container info")
		return nil, err
	}
	var containerInfo *DockerIndividualInfo
	err = json.NewDecoder(response.Body).Decode(&containerInfo)
	if err != nil {
		log.Errorf("docker container info could not be parsed")
		return nil, err
	}
	return containerInfo, nil
}

func (dkr *Docker) getDockerVersion() (*DockerVersionInfo, error) {
	response, err := dkr.requestFromDocker("version")
	if err != nil || response.StatusCode != 200 {
		log.Errorf("failed to fetch docker version")
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
