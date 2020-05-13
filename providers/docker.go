package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	dockerAPIVersion  string = "v1.40"
	suiEnabledLabel   string = "sui.enabled"
	suiProtectedLabel string = "sui.protected"
	suiURLLabel       string = "sui.url"
	suiIconLabel      string = "sui.icon"
	suiNameLabel      string = "sui.name"
)

type DockerProvider struct {
	Path   string
	Client *http.Client
}

type dockerVersion struct {
	Version string `json:"Version"`
	Os      string `json:"Os"`
}

type dockerContainer struct {
	Name   []string          `json:"Names"`
	Labels map[string]string `json:"Labels"`
}

func NewDockerProvider(placement uint8, path string) (*Provider, error) {

	var dockerClient DockerProvider
	dockerClient.Path = path
	dockerClient.Client = createDockerClient(path)

	var provider Provider
	provider.Title = "docker"
	provider.Placement = placement
	provider.Type = Docker
	provider.Config = &dockerClient
	provider.Apps = make(map[string]*App)

	version, err := checkDockerVersion(dockerClient.Client)
	if err != nil {
		return nil, fmt.Errorf("Could not communicate with docker over path %s", path)
	}

	log.Infof("Added Docker Provider\n")
	log.Debugf("docker version: %s\n", version.Version)

	return &provider, nil
}

func fetchDockerApps(p *Provider) error {
	log.Debugln("fetching apps from docker provider")
	cnf, valid := p.Config.(*DockerProvider)
	if !valid {
		return errors.New("Docker provider has invalid config")
	}

	containers := getContainerList(cnf.Client, true)
	if containers == nil {
		return fmt.Errorf("Could not fetch container list")
	}
	debugContainerList(containers)
	containerListToApps(p, containers)
	return nil
}

func getContainerList(client *http.Client, suiEnabled bool) []*dockerContainer {
	var containers []*dockerContainer
	response, err := requestFromSocket(client, "containers/json")
	if err != nil {
		log.Errorf("Failed to fetch container list from docker socket")
		return nil
	}
	err = json.NewDecoder(response.Body).Decode(&containers)
	if err != nil {
		log.Errorf("Failed to decode container list from docker socket")
		return nil
	}
	if !suiEnabled {
		return containers
	}
	for idx, container := range containers {
		isEnabled, ok := container.Labels[suiEnabledLabel]
		if !ok || isEnabled != "true" {
			containers = append(containers[:idx], containers[idx+1:]...)
		}
	}
	return containers
}

func containerListToApps(provider *Provider, dcl []*dockerContainer) {
	for _, container := range dcl {
		//parse name from labels
		name, ok := container.Labels[suiNameLabel]
		if !ok {
			if len(container.Name) > 0 {
				name = container.Name[0][1:]
			} else {
				log.Errorf("An enabled container has no name!")
				continue
			}
		}

		//parse protected from labels
		var protectBool bool
		protect, ok := container.Labels[suiProtectedLabel]
		if ok {
			var err error
			protectBool, err = strconv.ParseBool(protect)
			if err != nil {
				log.Errorf("Provided 'protect' (%s) value is not a bool!", protect)
				continue
			}
		} else {
			protectBool = false
		}

		//parse icon from labels
		icon, ok := container.Labels[suiIconLabel]
		if !ok {
			icon = "application"
		}

		//parse URL from labels
		URL, ok := container.Labels[suiURLLabel]
		if !ok {
			URL = "https://google.com"
		}

		provider.AddApp(name, icon, URL, protectBool)
	}
}

func debugContainerList(dcl []*dockerContainer) {
	if log.GetLevel() == log.DebugLevel {

	}
}

func checkDockerVersion(httpClient *http.Client) (dockerVersion, error) {
	response, err := requestFromSocket(httpClient, "version")
	if err == nil {
		var version dockerVersion
		err = json.NewDecoder(response.Body).Decode(&version)
		if err == nil {
			return version, nil
		}
	}
	return dockerVersion{}, err
}

func createDockerClient(path string) *http.Client {
	httpClient := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", path)
			},
		},
	}
	return &httpClient
}

func requestFromSocket(httpClient *http.Client, path string) (*http.Response, error) {
	path = fmt.Sprintf("http://127.0.0.1/%s/%s", dockerAPIVersion, path)
	response, err := httpClient.Get(path)
	return response, err
}
