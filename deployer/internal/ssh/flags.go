package ssh

import (
	"errors"
	"fmt"
)

var (
	ErrDialingFlagsInconsistentUsage   = errors.New("dialing flags inconsistent usage")
	ErrProxyJumpFlagsInconsistentUsage = errors.New("proxy jump flags inconsistent usage")
)

type Flags struct {
	Username string
	KeyFile  string

	ProxyJump ProxyJumpFlags

	ProxyCommand string

	Command string
}

func (f Flags) Validate() error {
	err := f.validateSSHDialingConsistency()
	if err != nil {
		return err
	}

	err = f.ProxyJump.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (f Flags) validateSSHDialingConsistency() error {
	if !f.ProxyJump.IsUsed() && f.ProxyCommand == "" && f.Command == "" {
		return nil
	}

	if f.ProxyJump.IsUsed() && f.ProxyCommand == "" && f.Command == "" {
		return nil
	}

	if !f.ProxyJump.IsUsed() && f.ProxyCommand != "" && f.Command == "" {
		return nil
	}

	if !f.ProxyJump.IsUsed() && f.ProxyCommand == "" && f.Command != "" {
		return nil
	}

	return fmt.Errorf("%w: specify either ProxyJump, ProxyCommand or Command; mixing is not supported", ErrDialingFlagsInconsistentUsage)
}

type ProxyJumpFlags struct {
	Address  string
	Username string
	KeyFile  string
}

func (f ProxyJumpFlags) Validate() error {
	if f.Address == "" && f.Username == "" && f.KeyFile == "" {
		return nil
	}

	if f.IsUsed() {
		return nil
	}

	return fmt.Errorf("%w: specify all ProxyJump flags or none of them", ErrProxyJumpFlagsInconsistentUsage)
}

func (f ProxyJumpFlags) IsUsed() bool {
	return f.Address != "" && f.Username != "" && f.KeyFile != ""
}
