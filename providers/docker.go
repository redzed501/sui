package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

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
	debugContainerList(containers)
	return nil
}

func getContainerList(client *http.Client, suiEnabled bool) []*dockerContainer {
	var containers []*dockerContainer
	response, err := requestFromSocket(client, "containers/json")
	err = json.NewDecoder(response.Body).Decode(&containers)
	if err != nil {
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

func debugContainerList(dcl []*dockerContainer) {
	if log.GetLevel() == log.DebugLevel {
		for _, dc := range dcl {
			log.Debugf("found container: %s\n", dc.Name[0][1:])
		}
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
