package deployer

import (
	"runtime"

	ver "go.maczukin.dev/libs/version"
)

var (
	version      = "dev"
	revision     = "HEAD"
	gitReference = "HEAD"
	builtAt      = "now"

	vi = ver.Info{
		Version:      version,
		Revision:     revision,
		GitReference: gitReference,
		BuiltAt:      builtAt,
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
	}
)

func VersionInfo() ver.Info {
	return vi
}
