package docker

import (
	"github.com/docker/docker/client"
)

func ConnectToDaemon() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	return cli
}
