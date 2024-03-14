package main

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type DockerClient struct {
	client *client.Client
}

func (dc *DockerClient) startContainer(image string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	response, err := dc.client.ContainerCreate(ctx, &container.Config{
		Image: image,
		Tty:   true,
	}, &hostConfig, &network.NetworkingConfig{}, &v1.Platform{}, containerName)

	return response.ID, err
}

func (dc *DockerClient) stopContainr(containerID string) {
	log.Println("[INFO] Stopping container", containerID)
	var timeout int = 3
	err := dc.client.ContainerStop(context.Background(), containerID, container.StopOptions{
		Timeout: &timeout,
	})

	if err != nil {
		log.Println("[ERROR] Unable to stop container.")
	}
}

func (dc *DockerClient) removeContainer(containerID string) {
	log.Println("[INFO] Removing container", containerID)
	err := dc.client.ContainerRemove(context.Background(), containerID, container.RemoveOptions{})

	if err != nil {
		log.Println("[ERROR] Unable to remove the container", containerID)
		// log.Fatalln(err)
	}
}

func (dc *DockerClient) extractContainerID(errMsg string) string {
	re := regexp.MustCompile(`[0-9a-fA-F]{64}`)
	match := re.FindStringSubmatch(errMsg)
	if len(match) >= 1 {
		return match[0]
	}

	panic("Unable to extract container ID.")
}
