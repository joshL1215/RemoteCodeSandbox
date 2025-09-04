package wrapper

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joshL1215/RemoteCodeSandbox/internal/models"
)

// Packs together runnable code submission file, entrypoint main file, and case json into a tmpdir
// Returns directory of tmpdir
func WrapJobDir(code string, language string, cases []models.Case) (string, error) {

	jobDir, err := os.MkdirTemp("", "job-*")
	if err != nil {
		return "", fmt.Errorf("Failed to create tempdir: %w", err)
	}

	exts := map[string]string{
		"python": ".py",
		"node":   ".js",
	}

	ext := exts[language]

	// Writes code submission
	if err := os.WriteFile(filepath.Join(jobDir, "submission"+ext), []byte(code), 0644); err != nil {
		return "", fmt.Errorf("Failed to write code file: %w", err)
	}

	// Writes a runner script to loop through cases
	if err := os.WriteFile(filepath.Join(jobDir, "main"+ext), []byte(PythonEntry), 0644); err != nil { // TODO: Dynamically choose entry code
		return "", fmt.Errorf("Failed to write entry file")
	}

	// Adds json of cases to tmpdir
	casesData, err := json.MarshalIndent(cases, "", "  ")
	if err != nil {
		return "", fmt.Errorf("Failed to marshal cases: %w", err)
	}

	// write JSON to file
	if err := os.WriteFile(filepath.Join(jobDir, "cases.json"), casesData, 0644); err != nil {
		return "", fmt.Errorf("Failed to write cases.json file: %w", err)
	}

	return jobDir, nil
}
