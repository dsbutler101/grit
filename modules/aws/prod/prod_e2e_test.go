//go:build e2e

package aws_prod

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	terratest_aws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

const (
	asgOutputKey = "autoscaling_group_names"
	testKey      = "AWS_PROD_TESTENDTOEND"
)

func TestEndToEnd(t *testing.T) {

	jobId := os.Getenv(test_tools.JobIdVar)
	require.NotEmpty(t, jobId)
	gitlabToken := os.Getenv(test_tools.GitlabTokenVar)
	require.NotEmpty(t, gitlabToken)
	client, err := gitlab.NewClient(gitlabToken)
	require.NoError(t, err)

	runnerToken := os.Getenv(test_tools.RunnerTokenVar)
	require.NotEmpty(t, runnerToken)

	// Create runner stack
	options := &terraform.Options{
		TerraformBinary: "terraform",
		TerraformDir:    ".",
		Vars: map[string]interface{}{
			"manager_service":       "ec2",
			"fleeting_service":      "ec2",
			"gitlab_project_id":     test_tools.GritEndToEndTestProjectID,
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
			"name":                  "e2e-" + jobId,
		},
	}
	_, err = terraform.InitAndApplyE(t, options)
	defer func() {
		terraform.Destroy(t, options)
	}()
	require.NoError(t, err)

	// TODO: implement check for runner stack health
	instanceName := "e2e-" + jobId + "_runner-manager"
	requireRunnerManagerRunning(t, instanceName)

	// Run a job
	main := "main"
	key := testKey
	uniqueValue := strconv.Itoa(int(rand.Uint32()))
	pipeline, _, err := client.Pipelines.CreatePipeline(test_tools.GritEndToEndTestProjectID, &gitlab.CreatePipelineOptions{
		Ref: &main,
		Variables: &[]*gitlab.PipelineVariableOptions{{
			Key:   &key,
			Value: &uniqueValue,
		}},
	})
	require.NoError(t, err)
	pipelineID := pipeline.ID

	var job *gitlab.Job
	jobs, _, err := client.Jobs.ListPipelineJobs(test_tools.GritEndToEndTestProjectID, pipelineID, &gitlab.ListJobsOptions{})
	require.NoError(t, err)
	require.Len(t, jobs, 1)
	jobID := jobs[0].ID

	// poll every second for 15 minutes for job completion
	for i := 0; i < 60*15; i++ {
		job, _, err = client.Jobs.GetJob(test_tools.GritEndToEndTestProjectID, jobID)
		require.NoError(t, err)

		if job.Status != "created" && job.Status != "pending" && job.Status != "running" {
			break
		}

		fmt.Println("Waiting for job. Current status:", job.Status)
		time.Sleep(time.Second)
	}

	require.Equal(t, "success", job.Status)
	logReader, _, err := client.Jobs.GetTraceFile(test_tools.GritEndToEndTestProjectID, job.ID)
	require.NoError(t, err)
	logBytes, err := io.ReadAll(logReader)
	require.NoError(t, err)
	log := string(logBytes)

	// Assert the job printed our unique value
	require.Contains(t, log, uniqueValue, fmt.Sprintf("looking for %v. found:\n%v", uniqueValue, log))

	// Assert the job ran on our stack
	asg, err := terraform.OutputE(t, options, asgOutputKey)
	require.NoError(t, err)
	autoscalingClient, err := terratest_aws.NewAsgClientE(t, test_tools.Region)
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

func requireRunnerManagerRunning(t *testing.T, instanceName string) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	svc := ec2.New(sess)

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(instanceName)},
			},
		},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	require.NotNil(t, result)
	require.Len(t, result.Reservations, 1)
	require.Len(t, result.Reservations[0].Instances, 1)
	require.Equal(t, "running", *result.Reservations[0].Instances[0].State.Name)
}
