package test_tools

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

const (
	JobIdVar                  = "CI_JOB_ID"
	CommitSHAVar              = "CI_COMMIT_SHA"
	GitlabTokenVar            = "GITLAB_TOKEN"
	GritEndToEndTestProjectID = 52010278
	Region                    = "us-east-1"
	RunnerTokenVar            = "RUNNER_TOKEN"
)

func JobName(_ *testing.T) string {
	jobId := os.Getenv(JobIdVar)
	sha := os.Getenv(CommitSHAVar)

	id := fmt.Sprintf("%s:%s:%d", jobId, sha, time.Now().Unix())
	hash := sha1.Sum([]byte(id))

	return fmt.Sprintf("u-%x", hash[:5])
}

func PlanAndAssert(t *testing.T, moduleVars map[string]any, expectedModules []string) {
	tempDir := t.TempDir()
	tempFilePath := filepath.Join(tempDir, "plan.json")

	options := &terraform.Options{
		TerraformBinary: "terraform",
		TerraformDir:    ".",
		PlanFilePath:    tempFilePath,
		Vars:            moduleVars,
	}

	plan := terraform.InitAndPlanAndShowWithStruct(t, options)

	assert.Len(t, plan.ResourcePlannedValuesMap, len(expectedModules))

	for _, v := range expectedModules {
		terraform.AssertResourceChangesMapKeyExists(t, plan, v)
	}

	// For easy troubleshooting, output keys that should be present
	fmt.Println(t.Name(), "ResourceChangesMap keys:")
	for k, _ := range plan.ResourceChangesMap {
		fmt.Printf("\"%v\",\n", k)
	}
}

func PlanAndAssertError(t *testing.T, moduleVars map[string]any, wantErr bool) {
	tempDir := t.TempDir()
	tempFilePath := filepath.Join(tempDir, "plan.json")

	options := &terraform.Options{
		TerraformBinary: "terraform",
		TerraformDir:    ".",
		PlanFilePath:    tempFilePath,
		Vars:            moduleVars,
	}

	_, err := terraform.InitAndPlanE(t, options)
	if wantErr {
		require.Error(t, err)
	} else {
		require.NoError(t, err)
	}
}
