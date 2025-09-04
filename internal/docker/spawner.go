package docker

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/joshL1215/RemoteCodeSandbox/internal/models"
	"github.com/joshL1215/RemoteCodeSandbox/internal/wrapper"
)

// Spawns a code execution container
func SpawnJob(cli *client.Client, language string, code string, cases []models.Case) (string, error) {

	ctx := context.Background()

	// Selects image for language
	images := map[string]string{
		"python": "python:3.12-alpine",
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
		Image: image,
		Cmd:
	}

}
