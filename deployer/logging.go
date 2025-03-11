package deployer

import (
	"os"
	"strconv"
)

const (
	logAddSourcesEnv   = "LOG_ADD_SOURCES"
	logCustomFormatEnv = "LOG_CUSTOM_FORMAT"
)

// These will be disabled during compilation time.
// Extended logging can be still enabled in the runtime
// by using a dedicated environment variables.
var (
	addSources      = "true"
	customLogFormat = "true"

	addSourcesB      bool
	customLogFormatB bool
)

func init() {
	if ok, err := strconv.ParseBool(envWithDefault(logAddSourcesEnv, addSources)); err == nil {
		addSourcesB = ok
	}

	if ok, err := strconv.ParseBool(envWithDefault(logCustomFormatEnv, customLogFormat)); err == nil {
		customLogFormatB = ok
	}
}

func envWithDefault(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		v = def
	}

	return v
}

func AddSources() bool {
	return addSourcesB
}

func CustomLogFormat() bool {
	return customLogFormatB
}
