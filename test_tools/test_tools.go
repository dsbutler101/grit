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
	JobIdVar     = "CI_JOB_ID"
	CommitSHAVar = "CI_COMMIT_SHA"
)

// JobName returns a unique name for a test
// that can be used for terraform resources.
// It should allow integration tests against the same
// cloud provider to run in parallel without conflicts.
func JobName(t *testing.T) string {
	// These CI vars can be empty for local dev
	jobId := os.Getenv(JobIdVar)
	sha := os.Getenv(CommitSHAVar)
	testName := t.Name()

	id := fmt.Sprintf("%s:%s:%s:%d", jobId, sha, testName, time.Now().Unix())
	hash := sha1.Sum([]byte(id))

	return fmt.Sprintf("u-%x", hash[:8])
}

func Plan(t *testing.T, moduleVars map[string]any) *terraform.PlanStruct {
	plan, err := PlanE(t, moduleVars)
	require.NoError(t, err)
	return plan
}

func PlanE(t *testing.T, moduleVars map[string]any) (*terraform.PlanStruct, error) {
	tempDir := t.TempDir()
	tempFilePath := filepath.Join(tempDir, "plan.json")

	options := terraform.Options{
		TerraformBinary: "terraform",
		TerraformDir:    ".",
		PlanFilePath:    tempFilePath,
		Vars:            moduleVars,
	}

	// We separate InitAndPlan from ShowWithStruct to use different options.
	//
	// In InitAndPlan we want to have a logger that will print terraform output
	// to the test output. That is very useful when debugging.
	_, err := terraform.InitAndPlanE(t, &options)
	if err != nil {
		return nil, err
	}

	optionsDiscardLogger := options
	optionsDiscardLogger.Logger = logger.Discard

	// In ShowWithStruct we use the options set that explicitly disable logger
	// by using logger.Discard. This is due the fact that to parse the plan output
	// for further analysis, ShowWithStruct needs to call `terraform show -json ...`
	// which will print the JSON representation of the plan file to the output. And
	// that may include sensitive values like tokens, passwords etc.
	// As many of these may be created dynamically, CI/CD variables masking will not
	// help in case of GitLab CI/CD execution.
	//
	// Therefore - we need to disable logging in this specific case.
	return terraform.ShowWithStruct(t, &optionsDiscardLogger), nil
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
	_, err := PlanE(t, moduleVars)

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
