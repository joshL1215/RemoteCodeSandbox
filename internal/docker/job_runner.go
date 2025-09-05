package docker

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
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
		return "", fmt.Errorf("language not supported: %s", language)
	}

	// Temp directory creation
	jobDir, err := wrapper.WrapJobDir(language, code, cases)
	if err != nil {
		return "", fmt.Errorf("encountered an error while wrapping code: %w", err)
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
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	containerID := cont.ID

	if err := cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to run container: %w", err)
	}

	resultCh, errCh := cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)

	// Watches over channels for container results
	select {
	case <-resultCh:
		// Got a result so move onto logging
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("container failed during run: %w", err)
		}
	case <-ctx.Done():
		return "", fmt.Errorf("timeout during container run: %w", err)
	}

	logStream, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: false,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}

	var stdoutBuffer, stderrBuffer bytes.Buffer
	_, err = stdcopy.StdCopy(&stdoutBuffer, &stderrBuffer, logStream)
	if err != nil {
		return "", fmt.Errorf("failed to read Docker logs: %w", err)
	}

	return stdoutBuffer.String(), nil
}
