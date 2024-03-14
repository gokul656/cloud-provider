package main

import (
	"log"

	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
)

const imageName = "ubuntu:latest"
const containerName = "go-container"

var dockerClient *DockerClient

func init() {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("[ERROR] %s", err.Error())
	}

	dockerClient = &DockerClient{
		client: client,
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("[INFO] Recovering from panic...")
			log.Fatalln("[ERROR]", r)
		}
	}()

	defer dockerClient.client.Close()

	// start a container
	containerID, err := dockerClient.startContainer(imageName)
	if err != nil {
		if errdefs.IsConflict(err) {
			log.Printf("[WARN] container name '%s' is already in use.", containerName)
			log.Println("[INFO] Trying to remove exisisting container..")

			// Trying to force remove container
			containerID := dockerClient.extractContainerID(err.Error())
			dockerClient.stopContainr(containerID)
			dockerClient.removeContainer(containerID)
		}
	} else {
		log.Println("[INFO] Container started successfully", containerID)
	}

	containerID, err = dockerClient.startContainer(imageName)
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		log.Println("[INFO] Container started successfully", containerID)
	}
}
