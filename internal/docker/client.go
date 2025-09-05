package docker

import (
	"fmt"

	"github.com/docker/docker/client"
)

func ConnectToDaemon() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("could not connect to docker daemon: %w", err)
	}
	return cli, nil
}
