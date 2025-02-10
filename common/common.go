package common

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	JobIdVar                  = "CI_JOB_ID"
	CommitSHAVar              = "CI_COMMIT_SHA"
	JobName                   = "JOB_NAME"
	GitlabTokenVar            = "GITLAB_TOKEN"
	GitlabTokenVarE2E         = "GITLAB_TOKEN_E2E"
	GritEndToEndTestProjectID = 52010278
	GitLabProjectID           = "CI_PROJECT_ID"
	RunnerTokenVar            = "RUNNER_TOKEN"

	TerraformHTTPAddress  = "TF_HTTP_ADDRESS"
	TerraformHTTPUsername = "TF_HTTP_USERNAME"
	TerraformHTTPPassword = "TF_HTTP_PASSWORD"
)

// E2ETestEnv is the environment for the end-to-end tests
// variables shared between tests for Go automation.
// Terraform variables are set with TF_VAR_ prefix.
type E2ETestEnv struct {
	Name            string
	GitlabToken     string
	GitLabProjectID string
	ProjectID       string

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

	je.Name, err = envVar(JobName)
	if err != nil {
		return nil, fmt.Errorf("setting up job: %w", err)
	}

	je.GitlabToken, err = envVar(GitlabTokenVarE2E)
	if err != nil {
		return nil, fmt.Errorf("setting up job: %w", err)
	}

	je.GitLabProjectID, err = envVar(GitLabProjectID)
	if err != nil {
		return nil, fmt.Errorf("setting up job: %w", err)
	}

	je.ProjectID, err = getGoogleProjectIDFromEnvVar()
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
