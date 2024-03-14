package main

import "github.com/docker/docker/api/types/container"

var resourceConfig = container.Resources{
	CPUShares: 2,
	Memory:    MBtoBytes(500),
}

var hostConfig = container.HostConfig{
	Resources: resourceConfig,
}

func MBtoBytes(memory int64) int64 {
	return memory * 1e+6
}
