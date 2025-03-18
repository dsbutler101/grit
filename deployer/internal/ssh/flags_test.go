package ssh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlags_Validate(t *testing.T) {
	testProxyJumpFlags := ProxyJumpFlags{
		Address:  "test-address",
		Username: "test-username",
		KeyFile:  "test-key-file",
	}

	defaultAssertError := func(t *testing.T, err error) {
		assert.ErrorIs(t, err, ErrDialingFlagsInconsistentUsage)
	}

	tests := map[string]struct {
		flags       Flags
		assertError func(t *testing.T, err error)
	}{
		"no dialing options used": {},
		"proxy jump specified": {
			flags: Flags{
				ProxyJump: testProxyJumpFlags,
			},
		},
		"proxy command specified": {
			flags: Flags{
				ProxyCommand: "test-proxy-command",
			},
		},
		"ssh command specified": {
			flags: Flags{
				Command: "test-ssh-command",
			},
		},
		"both proxy jump and proxy command specified": {
			flags: Flags{
				ProxyJump:    testProxyJumpFlags,
				ProxyCommand: "test-proxy-command",
			},
			assertError: defaultAssertError,
		},
		"both proxy jump and ssh command specified": {
			flags: Flags{
				ProxyJump: testProxyJumpFlags,
				Command:   "test-ssh-command",
			},
			assertError: defaultAssertError,
		},
		"both proxy command and ssh command specified": {
			flags: Flags{
				ProxyCommand: "test-proxy-command",
				Command:      "test-ssh-command",
			},
			assertError: defaultAssertError,
		},
		"invalid proxy jump configuration": {
			flags: Flags{
				ProxyJump: ProxyJumpFlags{
					Address: "address",
				},
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrProxyJumpFlagsInconsistentUsage)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			err := tt.flags.Validate()

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestProxyJumpFlags_Validate(t *testing.T) {
	defaultAssertError := func(t *testing.T, err error) {
		assert.ErrorIs(t, err, ErrProxyJumpFlagsInconsistentUsage)
	}

	tests := map[string]struct {
		flags       ProxyJumpFlags
		assertError func(t *testing.T, err error)
	}{
		"only address specified": {
			flags: ProxyJumpFlags{
				Address: "test-address",
			},
			assertError: defaultAssertError,
		},
		"only username specified": {
			flags: ProxyJumpFlags{
				Username: "test-username",
			},
			assertError: defaultAssertError,
		},
		"only key file specified": {
			flags: ProxyJumpFlags{
				KeyFile: "test-key-file",
			},
			assertError: defaultAssertError,
		},
		"all but address specified": {
			flags: ProxyJumpFlags{
				Username: "test-username",
				KeyFile:  "test-key-file",
			},
			assertError: defaultAssertError,
		},
		"all but username specified": {
			flags: ProxyJumpFlags{
				Address: "test-address",
				KeyFile: "test-key-file",
			},
			assertError: defaultAssertError,
		},
		"all but key file specified": {
			flags: ProxyJumpFlags{
				Address:  "test-address",
				Username: "test-username",
			},
			assertError: defaultAssertError,
		},
		"all specified": {
			flags: ProxyJumpFlags{
				Address:  "test-address",
				Username: "test-username",
				KeyFile:  "test-key-file",
			},
		},
		"none specified": {
			flags: ProxyJumpFlags{},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			err := tt.flags.Validate()

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
