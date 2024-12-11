//go:build e2e

package e2e

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

const (
	groupOutputKey = "instance_group_name"
	testKey        = "GOOGLE_PROD_GKE_WINDOWS_TESTENDTOEND"
)

type k8sCredentials struct {
	host          string
	token         string
	caCertificate []byte
}

func TestEndToEnd(t *testing.T) {
	name := os.Getenv("JOB_NAME")
	require.NotEmpty(t, name)

	gitlabToken := os.Getenv(test_tools.GitlabTokenVar)
	require.NotEmpty(t, gitlabToken)
	client, err := gitlab.NewClient(gitlabToken)
	require.NoError(t, err)

	tfHTTPAddress := os.Getenv(test_tools.TerraformHTTPAddress)
	require.NotEmpty(t, tfHTTPAddress)
	tfUsername := os.Getenv(test_tools.TerraformHTTPUsername)
	require.NotEmpty(t, tfUsername)
	tfPassword := os.Getenv(test_tools.TerraformHTTPPassword)
	require.NotEmpty(t, tfPassword)
	tfLockAddress := os.Getenv(test_tools.TerraformHTTPLockAddress)
	require.NotEmpty(t, tfLockAddress)
	tfUnlockAddress := os.Getenv(test_tools.TerraformHTTPUnlockAddress)
	require.NotEmpty(t, tfUnlockAddress)

	opts := &terraform.Options{
		TerraformBinary: "terraform",
		TerraformDir:    ".",
		BackendConfig: map[string]interface{}{
			"address":        tfHTTPAddress,
			"username":       tfUsername,
			"password":       tfPassword,
			"lock_address":   tfLockAddress,
			"unlock_address": tfUnlockAddress,
		},
	}
	terraform.Init(t, opts)
	output, err := terraform.OutputAllE(t, opts)
	require.NoError(t, err)

	cred := k8sCredentials{
		host:          output["host"].(string),
		token:         output["access_token"].(string),
		caCertificate: []byte(output["ca_certificate"].(string)),
	}

	// TODO: implement check for runner stack health
	requireRunnerManagerRunning(t, &cred, name)

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

	var job *gitlab.Job
	jobs, _, err := client.Jobs.ListPipelineJobs(test_tools.GritEndToEndTestProjectID, pipeline.ID, &gitlab.ListJobsOptions{})
	require.NoError(t, err, fmt.Sprintf("failed to list jobs in the pipeline %d", pipeline.ID))
	require.Len(t, jobs, 1, fmt.Sprintf("More than one job found in the pipeline %d", pipeline.ID))

	jobID := jobs[0].ID
	for {
		job, _, err = client.Jobs.GetJob(test_tools.GritEndToEndTestProjectID, jobID)
		require.NoError(t, err, fmt.Sprintf("failed to get job %d in the pipeline %d", jobID, pipeline.ID))

		if job.Status != "created" && job.Status != "pending" && job.Status != "running" {
			break
		}

		fmt.Println("Waiting for job. Current status:", job.Status)
		time.Sleep(10 * time.Second)
	}

	require.Equal(t, "success", job.Status)
	logReader, _, err := client.Jobs.GetTraceFile(test_tools.GritEndToEndTestProjectID, job.ID)
	require.NoError(t, err)
	logBytes, err := io.ReadAll(logReader)
	require.NoError(t, err)
	log := string(logBytes)

	// Assert the job printed our unique value
	require.Contains(t, log, uniqueValue, fmt.Sprintf("looking for %v. found:\n%v", uniqueValue, log))
}

func requireRunnerManagerRunning(t *testing.T, cred *k8sCredentials, instancePrefix string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cred.caCertificate)

	clientConfig := &rest.Config{
		Host:        cred.host,
		BearerToken: cred.token,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	client, err := kubernetes.NewForConfig(clientConfig)
	require.NoError(t, err, "Could not connect to the GKE Cluster")

	found, err := waitForPod(ctx, client, "gitlab-runner-system", instancePrefix)
	require.True(
		t,
		found,
		fmt.Sprintf("Could not find a pod name with prefix: %s-runner. Failed with error: %v", instancePrefix, err),
	)
}

func waitForPod(
	ctx context.Context,
	client kubernetes.Interface,
	ns string,
	prefix string,
) (bool, error) {
	watcher, err := client.CoreV1().Pods(ns).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to start watch on Pods: %w", err)
	}
	defer watcher.Stop()

	for {
		select {
		case <-ctx.Done():
			return false, fmt.Errorf("timed out waiting for pod with prefix %s", prefix)
		case event := <-watcher.ResultChan():
			switch event.Type {
			case watch.Added:
				pod, ok := event.Object.(*api.Pod)
				if !ok {
					continue
				}

				if strings.HasPrefix(pod.Name, prefix) {
					return true, nil
				}
			}
		}
	}
}
