package wrapper

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/gitlab-org/gitlab-runner/helpers/runner_wrapper/api"
)

func TestStatus_IsRunning(t *testing.T) {
	tests := map[api.Status]bool{
		api.StatusUnknown:    false,
		api.StatusRunning:    true,
		api.StatusInShutdown: true,
		api.StatusStopped:    false,
	}

	for s, r := range tests {
		status := Status{
			Status: s,
		}

		assert.Equal(t, r, status.IsRunning())
	}
}
