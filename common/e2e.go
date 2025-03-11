package common

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	GitlabTokenVarE2E = "GITLAB_TOKEN"
	GitLabProjectID   = "CI_PROJECT_ID"

	TerraformHTTPAddress  = "TF_HTTP_ADDRESS"
	TerraformHTTPUsername = "TF_HTTP_USERNAME"
	TerraformHTTPPassword = "TF_HTTP_PASSWORD"
)

// E2ETestEnv contains variables shared between Go automation
// and the terraform e2e modules.
type E2ETestEnv struct {
	GitlabToken     string
	GitLabProjectID string

	// HTTP state variables used to delete remote
	// state after destroy is run.
	TerraformHTTPAddress  string
	TerraformHTTPUsername string
	TerraformHTTPPassword string
}

func envVar(name string) (string, error) {
	raw, ok := os.LookupEnv(name)
	if !ok {
		return "", fmt.Errorf("env var %s not set", name)
	}
	cleaned := strings.Trim(raw, " ")
	if cleaned == "" {
		return "", fmt.Errorf("env var %s is empty", name)
	}
	return cleaned, nil
}

func getE2ETestEnv() (*E2ETestEnv, error) {
	var err error
	je := &E2ETestEnv{}

	je.GitlabToken, err = envVar(GitlabTokenVarE2E)
	if err != nil {
		return nil, fmt.Errorf("setting up job: %w", err)
	}

	je.GitLabProjectID, err = envVar(GitLabProjectID)
	if err != nil {
		return nil, fmt.Errorf("setting up job: %w", err)
	}

	je.TerraformHTTPAddress, err = envVar(TerraformHTTPAddress)
	if err != nil {
		return nil, fmt.Errorf("setting up job: %w", err)
	}

	je.TerraformHTTPUsername, err = envVar(TerraformHTTPUsername)
	if err != nil {
		return nil, fmt.Errorf("setting up job: %w", err)
	}

	je.TerraformHTTPPassword, err = envVar(TerraformHTTPPassword)
	if err != nil {
		return nil, fmt.Errorf("setting up job: %w", err)
	}

	return je, nil
}

func getAbsPathOfExec(executable string) (string, error) {
	execPath := ""

	execPath, err := exec.LookPath(executable)
	if err != nil {
		return "", err
	}

	execPath, err = filepath.Abs(execPath)
	if err != nil {
		return "", err
	}

	return execPath, nil
}
