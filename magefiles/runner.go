//go:build mage
// +build mage

package main

import (
	"context"

	"github.com/magefile/mage/mg"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/common"
)

type Runner mg.Namespace

// WaitForRunners waits for the runners to be online and fails if Runner doesn't come online.
func (Runner) WaitForRunners(ctx context.Context, runnerTag string) error {
	return common.WaitForRunners(ctx, runnerTag)
}
