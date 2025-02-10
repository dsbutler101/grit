package common

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func initTerraform(dir string) (*tfexec.Terraform, error) {
	tfPath, err := getAbsPathOfExec("terraform")
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path of terraform: %w", err)
	}

	tf, err := tfexec.NewTerraform(dir, tfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create terraform exec instance: %w", err)
	}

	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)

	return tf, err
}

func TerraformInitAndApply(ctx context.Context, dir string) error {
	tf, err := initTerraform(dir)
	if err != nil {
		return err
	}

	if err := tf.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize terraform: %w", err)
	}

	if err := tf.Apply(ctx); err != nil {
		return fmt.Errorf("failed to apply terraform: %w", err)
	}

	return nil
}

func TerraformInitAndDestroy(ctx context.Context, dir string) error {
	tf, err := initTerraform(dir)
	if err != nil {
		return err
	}

	if err := tf.Init(
		ctx,
	); err != nil {
		return fmt.Errorf("failed to initialize terraform: %w", err)
	}

	if err := tf.Destroy(ctx); err != nil {
		return fmt.Errorf("failed to destroy terraform infra: %w", err)
	}

	// Delete the GitLab State
	if _, err := deleteGitLabRemoteState(); err != nil {
		return fmt.Errorf("failed to delete GitLab remote state: %w", err)
	}

	return nil
}

func deleteGitLabRemoteState() ([]byte, error) {
	env, err := getE2ETestEnv()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", env.TerraformHTTPAddress, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Private-Token", env.TerraformHTTPPassword)

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
