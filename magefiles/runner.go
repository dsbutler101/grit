//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/common"
)

type Runner mg.Namespace

// WaitForRunners waits for the runners to be online and fails if Runner doesn't come online.
func (Runner) WaitForRunners(runnerTag string, tries int) error {
	return common.WaitForRunners(runnerTag, uint(tries))
}
