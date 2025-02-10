package common

import (
	"context"
	"fmt"
	"log"
	"os"

	backoff "github.com/cenkalti/backoff/v5"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func WaitForRunners(runnerTag string, maxRetries uint) error {
	log.SetOutput(os.Stdout)
	log.Println("Waiting for runner tagged '" + runnerTag + "' to be ready ...")

	env, err := getE2ETestEnv()
	if err != nil {
		return fmt.Errorf("failed to retrieve some environment variables: %w", err)
	}

	ctx := context.Background()
	runnerStatusChecker := func() (bool, error) {
		err := checkRunnersStatus(ctx, runnerTag, env)
		return err == nil, err
	}

	_, err = backoff.Retry(
		ctx,
		runnerStatusChecker,
		backoff.WithMaxTries(maxRetries),
	)

	return err
}

func checkRunnersStatus(ctx context.Context, runnerTag string, env *E2ETestEnv) error {

	runnerRetrieval := func() ([]*gitlab.Runner, error) {
		return getProjectRunners(env.GitlabToken, env.GitLabProjectID, "online", runnerTag)
	}
	runners, err := backoff.Retry(
		ctx,
		runnerRetrieval,
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithMaxTries(10),
	)
	if err != nil {
		return err
	}

	switch {
	case len(runners) != 1:
		return fmt.Errorf("no online runners found")
	default:
		return nil
	}
}

func getProjectRunners(gitlabToken, projectId string, status string, runnerTags ...string) ([]*gitlab.Runner, error) {
	glab, err := gitlab.NewClient(gitlabToken)
	if err != nil {
		return nil, fmt.Errorf("initializing glab client: %w", err)
	}

	runners, _, err := glab.Runners.ListProjectRunners(
		projectId,
		&gitlab.ListProjectRunnersOptions{
			ListOptions: gitlab.ListOptions{
				Page:    1,
				PerPage: 100,
			},
			Status:  &status,
			TagList: &runnerTags,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("retrieving project runners: %w", err)
	}

	return runners, nil
}
