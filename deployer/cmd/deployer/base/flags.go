package base

import (
	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/services/wait"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/ssh"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/wrapper"
)

func SetupAllTFFlags(cc *cobra.Command) *terraform.Flags {
	tfFlags := SetupTFFlagsTargetOnly(cc)

	cc.PersistentFlags().StringVar(&tfFlags.TargetStateFile, "tf-target-state-file", "", "Path to the Terraform state file for the Deployment Version to use")

	return tfFlags
}

func SetupTFFlagsTargetOnly(cc *cobra.Command) *terraform.Flags {
	tfFlags := &terraform.Flags{}

	cc.PersistentFlags().StringVar(&tfFlags.Target, "tf-target", "", "Path to the Terraform code directory for the Deployment Version to use")

	return tfFlags
}

func SetupWrapperFlags(cc *cobra.Command) *wrapper.Flags {
	wrapperFlags := &wrapper.Flags{}

	cc.PersistentFlags().DurationVar(&wrapperFlags.ConnectionTimeout, "wrapper-connection-timeout", wrapper.DefaultTimeout, "How long to wait for gRPC connection to be established before reporting error")

	return wrapperFlags
}

func SetupWaitFlags(cc *cobra.Command) *wait.Flags {
	waitFlags := &wait.Flags{}

	cc.PersistentFlags().DurationVar(&waitFlags.Timeout, "timeout", wait.DefaultTimeout, "How long the wait check should run before failing")

	return waitFlags
}

func SetupSSHFlags(cc *cobra.Command) *ssh.Flags {
	sshFlags := &ssh.Flags{}

	flags := cc.PersistentFlags()
	flags.StringVar(&sshFlags.Username, "ssh-username", "", "SSH target username")
	flags.StringVar(&sshFlags.KeyFile, "ssh-key-file", "", "Path to the SSH private key file")

	flags.StringVar(&sshFlags.ProxyJump.Address, "ssh-proxy-jump-address", "", "SSH Proxy Jump server address")
	flags.StringVar(&sshFlags.ProxyJump.Username, "ssh-proxy-jump-username", "", "SSH Proxy Jump server username")
	flags.StringVar(&sshFlags.ProxyJump.KeyFile, "ssh-proxy-jump-key-file", "", "Path to the SSH Proxy Jump server private key file")

	flags.StringVar(&sshFlags.ProxyCommand, "ssh-proxy-command", "", "SSH Proxy Command")

	flags.StringVar(&sshFlags.Command, "ssh-command", "", "SSH Command")

	return sshFlags
}
