//go:build mage
// +build mage

package main

import (
	"context"
	"time"

	"github.com/magefile/mage/mg"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/common"
)

type Terraform mg.Namespace

// TerraformInitAndApply runs terraform init and apply against the main.tf provided. The tfstate is then stored on GitLab.
func (Terraform) InitAndApply(dir string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return common.TerraformInitAndApply(ctx, dir)
}

// TerraformInitAndDestroy runs terraform init and destroy against the main.tf provided. The tfstate is retrieved from the GitLab.
func (Terraform) InitAndDestroy(dir string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return common.TerraformInitAndDestroy(ctx, dir)
}
