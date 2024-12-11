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
}

func getJobEnv() (*JobEnv, error) {
	var err error
	je := &JobEnv{}

	je.Name = os.Getenv(JobName)
	if strings.Trim(je.Name, " ") == "" {
		return nil, fmt.Errorf("job name is empty")
	}

	je.GritE2EDir = os.Getenv(GritE2eDirectory)
	if strings.Trim(je.GritE2EDir, " ") == "" {
		return nil, fmt.Errorf("directory containing the main.tf is empty")
	}

	je.GitlabToken = os.Getenv(GitlabTokenVar)
	if strings.Trim(je.GitlabToken, " ") == "" {
		return nil, fmt.Errorf("gitlab token is empty")
	}

	je.RunnerToken = os.Getenv(RunnerTokenPowerShellVar)
	if strings.Trim(je.RunnerToken, " ") == "" {
		return nil, fmt.Errorf("failed to retried runner token. Runner token %s is empty", RunnerTokenPowerShellVar)
	}

	je.ProjectID, err = getGoogleProjectIDFromEnvVar()
	if strings.Trim(je.ProjectID, " ") == "" {
		return nil, fmt.Errorf("failed to retried GCP Project ID. GCP Project ID is empty: %w", err)
	}

	je.Region, err = getGoogleRegionFromEnvVar()
	if strings.Trim(je.Region, " ") == "" {
		return nil, fmt.Errorf("failed to retried GCP Region. GCP Region is empty: %w", err)
	}

	je.Zone = os.Getenv("GOOGLE_ZONE")
	if strings.Trim(je.Zone, " ") == "" {
		return nil, fmt.Errorf("failed to retried GCP Zone. GCP Zone is empty")
	}

	je.JobTimeout = getJobTimeout(os.Getenv(CIJobTimeout))

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
