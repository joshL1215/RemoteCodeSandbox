package docker

import (
	"github.com/docker/docker/client"
)

// Spawns a code execution
func SpawnJob(cli *client.Client, language string, code string) {
	cli, err := client.NewClientWithOpts(client.FromEnv); err != nil {
		panic(err)
	}
}
