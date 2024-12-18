package common

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	JobIdVar                  = "CI_JOB_ID"
	CommitSHAVar              = "CI_COMMIT_SHA"
	JobName                   = "JOB_NAME"
	GitlabTokenVar            = "GITLAB_TOKEN"
	GritE2eDirectory          = "E2E_DIR"
	GritEndToEndTestProjectID = 52010278
	Region                    = "us-east-1"
	RunnerTokenVar            = "RUNNER_TOKEN"
	RunnerTokenPowerShellVar  = "RUNNER_TOKEN_POWERSHELL"

	TerraformHTTPAddress       = "TF_HTTP_ADDRESS"
	TerraformHTTPUsername      = "TF_HTTP_USERNAME"
	TerraformHTTPPassword      = "TF_HTTP_PASSWORD"
	TerraformHTTPLockAddress   = "TF_HTTP_LOCK_ADDRESS"
	TerraformHTTPUnlockAddress = "TF_HTTP_UNLOCK_ADDRESS"

	CIJobTimeout = "CI_JOB_TIMEOUT"
)

type JobEnv struct {
	Name        string
	GritE2EDir  string
	GitlabToken string
	RunnerToken string
	ProjectID   string
	Region      string
	Zone        string
	JobTimeout  time.Duration

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

func getJobEnv() (*JobEnv, error) {
	var err error
	je := &JobEnv{}

	je.Name, err = envVar(JobName)
	if err != nil {
		return nil, fmt.Errorf("setting up job: %w", err)
	}

	je.GritE2EDir, err = envVar(GritE2eDirectory)
	if err != nil {
		return nil, fmt.Errorf("setting up job: %w", err)
	}

	je.GitlabToken, err = envVar(GitlabTokenVar)
	if err != nil {
		return nil, fmt.Errorf("setting up job: %w", err)
	}

	je.RunnerToken, err = envVar(RunnerTokenPowerShellVar)
	if err != nil {
		return nil, fmt.Errorf("setting up job: %w", err)
	}

	je.ProjectID, err = getGoogleProjectIDFromEnvVar()
	if err != nil {
		return nil, fmt.Errorf("setting up job: %w", err)
	}

	je.Region, err = getGoogleRegionFromEnvVar()
	if strings.Trim(je.Region, " ") == "" {
		return nil, fmt.Errorf("setting up job: %w", err)
	}

	je.Zone, err = envVar("GOOGLE_ZONE")
	if err != nil {
		return nil, fmt.Errorf("setting up job: %w", err)
	}

	je.JobTimeout = getJobTimeout(os.Getenv(CIJobTimeout))

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

func getJobTimeout(timeout string) time.Duration {
	duration, err := time.ParseDuration(timeout)
	if err == nil {
		return duration
	}

	seconds, err := strconv.Atoi(timeout)
	if err != nil {
		return 1 * time.Hour
	}

	return time.Duration(seconds) * time.Second
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
