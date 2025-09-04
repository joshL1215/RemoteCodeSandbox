package docker

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/joshL1215/RemoteCodeSandbox/internal/models"
	"github.com/joshL1215/RemoteCodeSandbox/internal/wrapper"
)

// Spawns a code execution container
func SpawnJob(cli *client.Client, language string, code string, cases []models.Case) (string, error) {

	// 30 sec timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Selects image for language
	images := map[string]string{
		"python": "python:3.12-alpine", // TODO: probably add pre-pulled images to main.go
		"node":   "node:20-alpine",
	}

	image, ok := images[language]
	if !ok {
		return "", fmt.Errorf("%s is not a supported language", language)
	}

	// Temp directory creation
	jobDir, err := wrapper.WrapJobDir(language, code, cases)
	if err != nil {
		return "", fmt.Errorf("Encountered an error while wrapping code: %w", err)
	}
	defer os.RemoveAll(jobDir)

	config = &container.Config{
		Image:      image,
		Cmd:        []string{"python3", "/app/main.py"}, // TODO: other language support, this only allows Python
		WorkingDir: "/app",
		Tty:        false,
		User:       "1000:1000",
	}

	hostConfig := &container.HostConfig{ // TODO: Add gVisor runtime for isolation
		Binds: []string{
			fmt.Sprintf("%s:/app:ro", jobDir),
		},
		AutoRemove: true,
		Tmpfs: map[string]string{
			"/tmp": "size=64m",
		},
		CapDrop: []string{"ALL"},
		Resources: container.Resources{
			Memory:   64 * 1024 * 1024, // give container 64mb memory max
			NanoCPUs: 100_000_000,      // give container .1 cpus
		},
		NetworkMode:    "none",
		ReadonlyRootfs: true,
	}
}
