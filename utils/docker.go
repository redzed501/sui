package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	localDockerSocketPath string = "/var/run/docker.sock"
	localDockerAPIVersion string = "v1.40"
)

// GetDockerSocketConn creates a http connection to a docker socket
// An empty string for path will use the default local socket path
// This returns a http client for future requests
func GetDockerSocketConn(path string) (*http.Client, error) {
	if path == "" {
		path = localDockerSocketPath
	}
	if _, err := os.Stat(localDockerSocketPath); err == nil {
		httpClient := http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", path)
				},
			},
		}
		return &httpClient, nil
	} else if os.IsNotExist(err) {
		log.Errorf("You must mount the Docker Socket\n")
		log.Debugf("(mount /var/run/docker.sock as read only\n")
		return nil, fmt.Errorf("Docker Socket is not mounted (correctly)")
	} else {
		log.Errorf("Can not find thr Docker Socket\n")
		log.Debugf("(mount /var/run/docker.sock as read only\n")
		return nil, fmt.Errorf("Docker socket not found")
	}
}

func LocalDockerSockRequest(path string) (*http.Response, error) {
	client, err := GetDockerSocketConn("")
	if err != nil {
		log.Errorf("Failed to create docker socket connection\n")
		return nil, err
	}
	path = fmt.Sprintf("http://127.0.0.1/%s/%s", localDockerAPIVersion, path)
	response, err := client.Get(path)
	return response, err
}

func GetLocalDockerVersionInfo() (*DockerVersionInfo, error) {
	response, err := LocalDockerSockRequest("version")
	if err != nil || response.StatusCode != 200 {
		log.Errorf("Failed to fetch local docker version")
		return nil, err
	}
	var versionInfo *DockerVersionInfo
	err = json.NewDecoder(response.Body).Decode(&versionInfo)
	if err != nil {
		log.Errorf("Docker version info could not be parsed\n")
		return nil, err
	}
	return versionInfo, nil
}

func GetLocalContainerList() (*ContainerList, error) {
	response, err := LocalDockerSockRequest("containers/json")
	if err != nil || response.StatusCode != 200 {
		log.Errorf("Failed to fetch local docker container list")
		return nil, err
	}
	var containerList *ContainerList
	err = json.NewDecoder(response.Body).Decode(&containerList)
	if err != nil {
		log.Errorf("Docker container list could not be parsed\n")
		return nil, err
	}
	return containerList, nil
}

func DockerOk() bool {
	versionInfo, err := GetLocalDockerVersionInfo()
	if err != nil {
		log.Errorf("Docker Connection Not Avaliable")
		return false
	}
	log.Debugf("Docker Connection Made\n\tdocker version: %s\n", versionInfo.Version)
	return true
}
