package utils

type ContainerList struct {
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
