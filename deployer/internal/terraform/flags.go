package terraform

import (
	"errors"
)

type Flags struct {
	Target          string
	TargetStateFile string
}

func (f *Flags) Validate() error {
	if f.Target == "" && f.TargetStateFile == "" {
		return errors.New("one of target or target-state-file must be specified")
	}

	if f.Target != "" && f.TargetStateFile != "" {
		return errors.New("only one of target or target-state-file can be specified at once")
	}

	return nil
}
