package ssh

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestConfig_SSHClientConfig(t *testing.T) {
	testAddress := "test-address:1234"
	testUsername := "test-username"

	tests := map[string]struct {
		config             Config
		assertError        func(t *testing.T, err error)
		assertSSHCliConfig func(t *testing.T, sshCfg *ssh.ClientConfig)
	}{
		"invalid key": {
			config: Config{
				KeyPemBytes: nil,
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errParsingPrivateKey)
			},
		},
		"valid configuration": {
			config: Config{
				Address:     testAddress,
				Username:    testUsername,
				KeyPemBytes: []byte(testSSHUserPrivateKeyPem),
			},
			assertSSHCliConfig: func(t *testing.T, sshCfg *ssh.ClientConfig) {
				assert.Equal(t, testUsername, sshCfg.User)
				assert.GreaterOrEqual(t, len(sshCfg.Auth), 1)
				assert.Equal(t, sshClientConnectionTimeout, sshCfg.Timeout)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			sshCliConf, err := tt.config.SSHClientConfig()

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, tt.assertSSHCliConfig, "assertSSHCliConfig must be defined in test definition")
			tt.assertSSHCliConfig(t, sshCliConf)
		})
	}
}

func TestConfig_Network(t *testing.T) {
	var cfg Config

	assert.Equal(t, sshNetwork, cfg.Network())
}
