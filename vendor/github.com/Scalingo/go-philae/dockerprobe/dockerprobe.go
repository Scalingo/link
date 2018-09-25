package dockerprobe

import (
	docker "github.com/fsouza/go-dockerclient"
	errgo "gopkg.in/errgo.v1"
)

type DockerProbe struct {
	name     string
	endpoint string
}

func NewDockerProbe(name, endpoint string) DockerProbe {
	return DockerProbe{
		name:     name,
		endpoint: endpoint,
	}
}

func (p DockerProbe) Name() string {
	return p.name
}

func (p DockerProbe) Check() error {
	client, err := docker.NewClient(p.endpoint)
	if err != nil {
		return errgo.Notef(err, "Unable to create")
	}

	_, err = client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return errgo.Notef(err, "Unable to contact docker")
	}

	return nil
}
