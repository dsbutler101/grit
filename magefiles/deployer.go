//go:build mage

package main

import (
	"context"
	"strings"

	"github.com/magefile/mage/mg"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/magefiles/deployer"
)

type Deployer mg.Namespace

// Compile compiles the Deployer binary for the current platform
func (Deployer) Compile() error {
	return deployer.Compile()
}

// CompileFor compiles the Deployer binary for the provided platforms
func (Deployer) CompileFor(platforms string) error {
	var platformDefs []string
	if platforms != "" {
		platformDefs = strings.SplitN(platforms, " ", -1)
	}

	return deployer.Compile(platformDefs...)
}

// Upload uploads compiled Deployer binaries to GitLab generic packages repository
// Works only from within CI/CD job - authentication hardcoded to use
// CI_JOB_TOKEN.
func (Deployer) Upload(ctx context.Context) error {
	return deployer.Upload(ctx)
}

// GoGenerate runs go generate within the deployer submodule
func (Deployer) GoGenerate() error {
	return deployer.GoGenerate()
}
