package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/joshL1215/RemoteCodeSandbox/internal/models"
	"github.com/joshL1215/RemoteCodeSandbox/internal/wrapper"
)

// Spawns a code execution container
func RunJudgeJob(cli *client.Client, language string, code string, cases []models.Case) (string, error) {

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

	config := &container.Config{
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
	// TODO: maybe put configs in a different file

	cont, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("Failed to create container: %w", err)
	}

	containerID := cont.ID

	if err := cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("Failed to run container: %w", err)
	}

	resultCh, errCh := cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)

	// Watches over channels for container results
	select {
	case <-resultCh:
		// Got a result so move onto logging
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("Container failed during run: %w")
		}
	case <-ctx.Done():
		return "", fmt.Errorf("Timeout during container run")
	}

	logStream, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", fmt.Errorf("Failed to get container logs: %w", err)
	}

	logs, err := io.ReadAll(logStream)
	if err != nil {
		return "", fmt.Errorf("Failed to read log stream: %w", err)
	}

	return string(logs), nil
}
