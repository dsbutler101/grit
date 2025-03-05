package terraform

import (
	"errors"
)

var (
	ErrNoTargetFlagSet    = errors.New("one of target or target-state-file must be specified")
	ErrBothTargetFlagsSet = errors.New("only one of target or target-state-file can be specified at once")
)

type Flags struct {
	Target          string
	TargetStateFile string
}

func (f *Flags) Validate() error {
	if f.Target == "" && f.TargetStateFile == "" {
		return ErrNoTargetFlagSet
	}

	if f.Target != "" && f.TargetStateFile != "" {
		return ErrBothTargetFlagsSet
	}

	return nil
}
