package test_tools

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

const (
	JobIdVar                  = "CI_JOB_ID"
	GitlabTokenVar            = "GITLAB_TOKEN"
	GritEndToEndTestProjectID = 52010278
	Region                    = "us-east-1"
	RunnerTokenVar            = "RUNNER_TOKEN"
)

func JobName(t *testing.T) string {
	jobId := os.Getenv(JobIdVar)
	require.NotEmpty(t, jobId)
	return "unit-" + jobId
}

func PlanAndAssert(t *testing.T, moduleVars map[string]interface{}, expectedModules []string) {
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
