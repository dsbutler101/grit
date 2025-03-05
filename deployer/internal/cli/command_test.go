package cli

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestBuildRunEFromCommandExecutor(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	cmd := new(cobra.Command)
	cmd.SetContext(ctx)

	testArgs := []string{"arg-1", "arg-2"}

	ce := NewMockCommandExecutor(t)
	ce.EXPECT().Execute(ctx, cmd, testArgs).Return(assert.AnError)

	runE := BuildRunEFromCommandExecutor(ce)
	err := runE(cmd, testArgs)
	assert.ErrorIs(t, err, assert.AnError)
}
