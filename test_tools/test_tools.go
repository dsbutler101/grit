package test_tools

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
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

func Plan(t *testing.T, moduleVars map[string]any) *terraform.PlanStruct {
	tempDir := t.TempDir()
	tempFilePath := filepath.Join(tempDir, "plan.json")

	options := terraform.Options{
		TerraformBinary: "terraform",
		TerraformDir:    ".",
		PlanFilePath:    tempFilePath,
		Vars:            moduleVars,
	}

	optionsDiscardLogger := options
	optionsDiscardLogger.Logger = logger.Discard

	// We separate InitAndPlan from ShowWithStruct to use different options.
	//
	// In InitAndPlan we want to have a logger that will print terraform output
	// to the test output. That is very useful when debugging.
	_ = terraform.InitAndPlan(t, &options)

	// In ShowWithStruct we use the options set that explicitly disable logger
	// by using logger.Discard. This is due the fact that to parse the plan output
	// for further analysis, ShowWithStruct needs to call `terraform show -json ...`
	// which will print the JSON representation of the plan file to the output. And
	// that may include sensitive values like tokens, passwords etc.
	// As many of these may be created dynamically, CI/CD variables masking will not
	// help in case of GitLab CI/CD execution.
	//
	// Therefore - we need to disable logging in this specific case.
	return terraform.ShowWithStruct(t, &optionsDiscardLogger)
}

func PlanAndAssert(t *testing.T, moduleVars map[string]any, expectedModules []string) {
	plan := Plan(t, moduleVars)
	AssertWithPlan(t, plan, expectedModules)
}

func AssertWithPlan(t *testing.T, plan *terraform.PlanStruct, expectedModules []string) {
	assert.Len(t, plan.ResourcePlannedValuesMap, len(expectedModules))

	for _, v := range expectedModules {
		terraform.AssertResourceChangesMapKeyExists(t, plan, v)
	}

	// For easy troubleshooting, output keys that should be present
	fmt.Println(t.Name(), "ResourceChangesMap keys:")
	for k := range plan.ResourceChangesMap {
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

func AssertProviderConfigExists(t *testing.T, plan *terraform.PlanStruct, name string) {
	t.Helper()
	assert.Contains(t, plan.RawPlan.Config.ProviderConfigs, name, `Expected provider config "%s" to exist, but does not`, name)
}
