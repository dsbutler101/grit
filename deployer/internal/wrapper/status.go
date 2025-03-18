package wrapper

import (
	"gitlab.com/gitlab-org/gitlab-runner/helpers/runner_wrapper/api"
)

type Status struct {
	Status        api.Status
	FailureReason string
}

func (s Status) IsRunning() bool {
	return s.Status == api.StatusRunning || s.Status == api.StatusInShutdown
}
