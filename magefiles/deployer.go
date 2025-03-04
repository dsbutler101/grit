//go:build mage

package main

import (
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

// GoGenerate runs go generate within the deployer submodule
func (Deployer) GoGenerate() error {
	return deployer.GoGenerate()
}
