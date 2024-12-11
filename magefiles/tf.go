package main

import (
	"github.com/magefile/mage/mg"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/common"
)

type Terraform mg.Namespace

// TerraformInitAndApply runs terraform init and apply against the main.tf provided. The tfstate is then stored on GitLab.
func (Terraform) TerraformInitAndApply() error {
	return common.TerraformInitAndApply()
}

// TerraformInitAndDestroy runs terraform init and destroy against the main.tf provided. The tfstate is retrieved from the GitLab.
func (Terraform) TerraformInitAndDestroy() error {
	return common.TerraformInitAndDestroy()
}
