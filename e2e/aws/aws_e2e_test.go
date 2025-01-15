//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	terratest_aws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/common"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

const (
	asgOutputKey = "autoscaling_group_name"
	testKey      = "AWS_PROD_TESTENDTOEND"
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

	// Create runner stack
	options := &terraform.Options{
		TerraformBinary: "terraform",
		TerraformDir:    ".",
		Vars: map[string]interface{}{
			"runner_token": runnerToken,
			"name":         name,
			"job_id":       jobId,
		},
	}
	_, err = terraform.InitAndApplyE(t, options)
	defer func() {
		terraform.Destroy(t, options)
	}()
	require.NoError(t, err)

	// TODO: implement check for runner stack health
	instanceName := name + "_runner-manager"
	requireRunnerManagerRunning(t, instanceName)

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

	// Assert the job ran on our stack
	asg, err := terraform.OutputE(t, options, asgOutputKey)
	require.NoError(t, err)
	autoscalingClient, err := terratest_aws.NewAsgClientE(t, common.Region)
	require.NoError(t, err)
	groups, err := autoscalingClient.DescribeAutoScalingGroups(context.TODO(), &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{asg},
	})
	require.NoError(t, err)
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

func requireRunnerManagerRunning(t *testing.T, instanceName string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	require.NoError(t, err)

	svc := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []string{instanceName},
			},
		},
	}

	result, err := svc.DescribeInstances(ctx, input)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	require.NotNil(t, result)
	require.Len(t, result.Reservations, 1)
	require.Len(t, result.Reservations[0].Instances, 1)
	require.Equal(t, "running", result.Reservations[0].Instances[0].State.Name)
}
