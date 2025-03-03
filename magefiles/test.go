//go:build mage

package main

import (
	"github.com/magefile/mage/mg"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/magefiles/tests"
)

type Test mg.Namespace

// Unit runs Go unit tests on a given path
func (Test) Unit(path string) error {
	return tests.UnitForPath(path)
}
