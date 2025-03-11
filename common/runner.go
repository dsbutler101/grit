package common

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	backoff "github.com/cenkalti/backoff/v5"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func WaitForRunners(ctx context.Context, runnerTag string) error {
	log.SetOutput(os.Stdout)
	log.Printf("Waiting for runner tagged %q to be ready ...", runnerTag)

	env, err := getE2ETestEnv()
	if err != nil {
		return fmt.Errorf("failed to retrieve some environment variables: %w", err)
	}

	glab, err := gitlab.NewClient(env.GitlabToken)
	if err != nil {
		return fmt.Errorf("initializing glab client: %w", err)
	}

	_, err = backoff.Retry(
		ctx,
		func() (bool, error) {
			runner, err := checkRunnersStatus(glab, runnerTag, env)
			if err == nil {
				runner.Token = "REDACTED"
				log.Printf("success: %#v", runner)
			}
			return err == nil, err
		},
		backoff.WithBackOff(&backoff.ExponentialBackOff{
			InitialInterval: time.Second,
			Multiplier:      2,
			MaxInterval:     time.Minute,
		}),
		backoff.WithNotify(func(err error, d time.Duration) {
			log.Printf("error: %s", err)
			log.Printf("trying again in %s", d)
		}),
	)

	return err
}

func checkRunnersStatus(glab *gitlab.Client, runnerTag string, env *E2ETestEnv) (*gitlab.Runner, error) {
	runners, err := getProjectRunners(glab, env.GitLabProjectID, "online", runnerTag)
	if err != nil {
		return nil, err
	}

	if len(runners) != 1 {
		return nil, fmt.Errorf("no online runners found with tag '%s'", runnerTag)
	}

	return runners[0], nil
}

func getProjectRunners(glab *gitlab.Client, projectId string, status string, runnerTags ...string) ([]*gitlab.Runner, error) {
	runners, _, err := glab.Runners.ListProjectRunners(
		projectId,
		&gitlab.ListProjectRunnersOptions{
			ListOptions: gitlab.ListOptions{
				Page:    1,
				PerPage: 1,
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
