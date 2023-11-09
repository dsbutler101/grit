//go:build e2e

package prod

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"
)

const (
	asgOutputKey              = "autoscaling_group_names"
	gitlabTokenVar            = "GITLAB_TOKEN"
	gritEndToEndTestProjectID = 52010278
	region                    = "us-east-1"
	runnerTokenVar            = "RUNNER_TOKEN"
	testKey                   = "AWS_PROD_TESTENDTOEND"
)

func TestEndToEnd(t *testing.T) {

	gitlabToken := os.Getenv(gitlabTokenVar)
	require.NotEmpty(t, gitlabToken)
	client, err := gitlab.NewClient(gitlabToken)
	require.NoError(t, err)

	runnerToken := os.Getenv(runnerTokenVar)
	require.NotEmpty(t, runnerToken)

	// Create runner stack
	options := &terraform.Options{
		TerraformBinary: "terraform",
		TerraformDir:    ".",
		Vars: map[string]interface{}{
			"manager_service":       "ec2",
			"fleeting_service":      "ec2",
			"gitlab_project_id":     gritEndToEndTestProjectID,
			"gitlab_runner_tags":    []string{t.Name()},
			"fleeting_os":           "linux",
			"ami":                   "ami-0735db9b38fcbdb39",
			"instance_type":         "t2.medium",
			"aws_vpc_cidr":          "10.0.0.0/24",
			"capacity_per_instance": 1,
			"scale_min":             1,
			"scale_max":             1,
			"executor":              "docker-autoscaler",
			"min_maturity":          "alpha",
			"runner_token":          runnerToken,
		},
	}
	_, err = terraform.InitAndApplyE(t, options)
	defer func() {
		terraform.Destroy(t, options)
	}()
	require.NoError(t, err)

	// TODO: implement check for runner stack health
	time.Sleep(5 * time.Minute)

	// Run a job
	main := "main"
	key := testKey
	uniqueValue := strconv.Itoa(int(rand.Uint32()))
	pipeline, _, err := client.Pipelines.CreatePipeline(gritEndToEndTestProjectID, &gitlab.CreatePipelineOptions{
		Ref: &main,
		Variables: &[]*gitlab.PipelineVariableOptions{{
			Key:   &key,
			Value: &uniqueValue,
		}},
	})
	require.NoError(t, err)
	pipelineID := pipeline.ID

	// TODO: poll for job completion
	time.Sleep(time.Minute)
	jobs, _, err := client.Jobs.ListPipelineJobs(gritEndToEndTestProjectID, pipelineID, &gitlab.ListJobsOptions{})
	require.Len(t, jobs, 1)
	job := jobs[0]
	require.Equal(t, "success", job.Status)
	logReader, _, err := client.Jobs.GetTraceFile(gritEndToEndTestProjectID, job.ID)
	require.NoError(t, err)
	logBytes, err := io.ReadAll(logReader)
	require.NoError(t, err)
	log := string(logBytes)

	// Assert the job printed our unique value
	require.Contains(t, log, uniqueValue, fmt.Sprintf("looking for %v. found:\n%v", uniqueValue, log))

	// Assert the job ran on our stack
	asg, err := terraform.OutputE(t, options, asgOutputKey)
	require.NoError(t, err)
	autoscalingClient, err := aws.NewAsgClientE(t, region)
	require.NoError(t, err)
	groups, err := autoscalingClient.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{&asg},
	})
	if len(groups.AutoScalingGroups) != 1 {
		t.Fatalf("expected 1 asg. found %v", len(groups.AutoScalingGroups))
	}
	jobRanOnAsg := false
	for _, instance := range groups.AutoScalingGroups[0].Instances {
		id := instance.InstanceId
		if id == nil {
			continue
		}
		if strings.Contains(log, *id) {
			jobRanOnAsg = true
		}
	}
	require.True(t, jobRanOnAsg)
}
