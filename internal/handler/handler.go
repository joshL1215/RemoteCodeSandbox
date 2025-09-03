package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func InputHandler(cli *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
		if err != nil {
			panic(err)
		}

		for _, ctr := range containers {
			fmt.Printf("%s %s\n", ctr.ID, ctr.Image)
		}
	}
}
