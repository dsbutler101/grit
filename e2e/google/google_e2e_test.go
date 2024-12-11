//go:build e2e

package e2e

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	terratest_gcp "github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/common"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

const (
	groupOutputKey = "instance_group_name"
	testKey        = "GOOGLE_PROD_TESTENDTOEND"
)

func TestEndToEnd(t *testing.T) {
	name := test_tools.JobName(t)

	jobId := os.Getenv(common.JobIdVar)
	require.NotEmpty(t, jobId)
	gitlabToken := os.Getenv(common.GitlabTokenVar)
	require.NotEmpty(t, gitlabToken)
	client, err := gitlab.NewClient(gitlabToken)
	require.NoError(t, err)

	runnerToken := os.Getenv(common.RunnerTokenVar)
	require.NotEmpty(t, runnerToken)

	projectID := terratest_gcp.GetGoogleProjectIDFromEnvVar(t)
	require.NotEmpty(t, projectID)

	region := terratest_gcp.GetGoogleRegionFromEnvVar(t)
	require.NotEmpty(t, region)

	zone := os.Getenv("GOOGLE_ZONE")
	require.NotEmpty(t, zone)

	// Create runner stack
	options := &terraform.Options{
		TerraformBinary: "terraform",
		TerraformDir:    ".",
		Vars: map[string]interface{}{
			"runner_token":   runnerToken,
			"name":           name,
			"job_id":         jobId,
			"google_region":  region,
			"google_zone":    zone,
			"google_project": projectID,
		},
	}
	_, err = terraform.InitAndApplyE(t, options)
	defer func() {
		terraform.Destroy(t, options)
	}()
	require.NoError(t, err)

	// TODO: implement check for runner stack health
	requireRunnerManagerRunning(t, projectID, zone, name)

	groupName, err := terraform.OutputE(t, options, groupOutputKey)
	require.NoError(t, err)
	require.NotEmpty(t, groupName)

	// Run a job
	main := "main"
	key := testKey
	uniqueValue := strconv.Itoa(int(rand.Uint32()))
	pipeline, _, err := client.Pipelines.CreatePipeline(common.GritEndToEndTestProjectID, &gitlab.CreatePipelineOptions{
		Ref: &main,
		Variables: &[]*gitlab.PipelineVariableOptions{{
			Key:   &key,
			Value: &uniqueValue,
		}},
	})
	require.NoError(t, err)
	pipelineID := pipeline.ID

	var job *gitlab.Job
	jobs, _, err := client.Jobs.ListPipelineJobs(common.GritEndToEndTestProjectID, pipelineID, &gitlab.ListJobsOptions{})
	require.NoError(t, err)
	require.Len(t, jobs, 1)
	jobID := jobs[0].ID

	// poll every second for 15 minutes for job completion
	for i := 0; i < 60*15; i++ {
		job, _, err = client.Jobs.GetJob(common.GritEndToEndTestProjectID, jobID)
		require.NoError(t, err)

		if job.Status != "created" && job.Status != "pending" && job.Status != "running" {
			break
		}

		fmt.Println("Waiting for job. Current status:", job.Status)
		time.Sleep(time.Second)
	}

	require.Equal(t, "success", job.Status)
	logReader, _, err := client.Jobs.GetTraceFile(common.GritEndToEndTestProjectID, job.ID)
	require.NoError(t, err)
	logBytes, err := io.ReadAll(logReader)
	require.NoError(t, err)
	log := string(logBytes)

	// Assert the job printed our unique value
	require.Contains(t, log, uniqueValue, fmt.Sprintf("looking for %v. found:\n%v", uniqueValue, log))
	// Assert the job ran on our stack (Matches `Instance https://www.googleapis.com/compute/v1/projects/ID/zones/ZONE/instances/u-7a1c395d97-9886c784ca0f316a connected`)
	require.Contains(t, log, fmt.Sprintf("/instances/%s-", groupName), fmt.Sprintf("looking for %v. found:\n%v", fmt.Sprintf("/instances/%s-", groupName), log))
}

func requireRunnerManagerRunning(t *testing.T, projectID string, zone string, instancePrefix string) {
	s, err := terratest_gcp.NewInstancesServiceE(t)
	require.NoError(t, err)
	list, err := s.List(projectID, zone).Do()
	require.NoError(t, err)
	found := false
	for _, instance := range list.Items {
		if instance.Name == instancePrefix+"-runner-manager" {
			require.Equal(t, "RUNNING", instance.Status)
			found = true
			break
		}
	}
	require.True(t, found)
}
