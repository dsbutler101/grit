package common

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func initTerraform() (*tfexec.Terraform, *JobEnv, error) {
	je, err := getJobEnv()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieved some environment variables: %w", err)
	}

	tfPath, err := getAbsPathOfExec("terraform")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get absolute path of terraform: %w", err)
	}

	tf, err := tfexec.NewTerraform(je.GritE2EDir, tfPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create terraform exec instance: %w", err)
	}

	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)

	return tf, je, err
}

func TerraformInitAndApply() error {
	tf, je, err := initTerraform()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), je.JobTimeout)
	defer cancel()

	if err := tf.Init(
		ctx,
	); err != nil {
		return fmt.Errorf("failed to initialize terraform: %w", err)
	}

	if err := tf.Apply(
		ctx,
		tfexec.Var(fmt.Sprintf("gitlab_runner_token=%s", je.RunnerToken)),
		tfexec.Var(fmt.Sprintf("name=%s", je.Name)),
		tfexec.Var(fmt.Sprintf("google_region=%s", je.Region)),
		tfexec.Var(fmt.Sprintf("google_zone=%s", je.Zone)),
		tfexec.Var(fmt.Sprintf("google_project=%s", je.ProjectID)),
	); err != nil {
		return fmt.Errorf("failed to apply terraform: %w", err)
	}

	return nil
}

func TerraformInitAndDestroy() error {
	tf, je, err := initTerraform()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), je.JobTimeout)
	defer cancel()

	if err := tf.Init(
		ctx,
	); err != nil {
		return fmt.Errorf("failed to initialize terraform: %w", err)
	}

	if err := tf.Destroy(
		ctx,
		tfexec.Var(fmt.Sprintf("gitlab_runner_token=%s", je.RunnerToken)),
		tfexec.Var(fmt.Sprintf("name=%s", je.Name)),
		tfexec.Var(fmt.Sprintf("google_region=%s", je.Region)),
		tfexec.Var(fmt.Sprintf("google_zone=%s", je.Zone)),
		tfexec.Var(fmt.Sprintf("google_project=%s", je.ProjectID)),
	); err != nil {
		return fmt.Errorf("failed to destroy terraform infra: %w", err)
	}

	// Delete the GitLab State
	if _, err := deleteGitLabRemoteState(je); err != nil {
		return fmt.Errorf("failed to delete GitLab remote state: %w", err)
	}

	return nil
}

func deleteGitLabRemoteState(je *JobEnv) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", je.TerraformHTTPAddress, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Private-Token", je.TerraformHTTPPassword)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func TerraformOutput() (map[string]string, error) {
	tf, je, err := initTerraform()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), je.JobTimeout)
	defer cancel()

	state := make(map[string]string)

	tfState, err := tf.Output(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get terraform output: %w", err)
	}

	for k, v := range tfState {
		state[k] = string(v.Value)
	}

	return state, nil
}
