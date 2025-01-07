package common

import (
	"fmt"
	"os"
)

var projectEnvVars = []string{
	"GOOGLE_PROJECT",
	"GOOGLE_CLOUD_PROJECT",
	"GOOGLE_CLOUD_PROJECT_ID",
	"GCLOUD_PROJECT",
	"CLOUDSDK_CORE_PROJECT",
}

var regionEnvVars = []string{
	"GOOGLE_REGION",
	"GCLOUD_REGION",
	"CLOUDSDK_COMPUTE_REGION",
}

func getFirstNonEmptyEnvVar(envVarNames []string) (string, error) {
	value := getFirstNonEmptyEnvVarOrEmptyString(envVarNames)
	if value == "" {
		return "", fmt.Errorf(
			"all of the following env vars %s are empty. At least one must be non-empty",
			envVarNames,
		)
	}

	return value, nil
}

func getFirstNonEmptyEnvVarOrEmptyString(envVarNames []string) string {
	for _, name := range envVarNames {
		if value := os.Getenv(name); value != "" {
			return value
		}
	}

	return ""
}

func getGoogleProjectIDFromEnvVar() (string, error) {
	return getFirstNonEmptyEnvVar(projectEnvVars)
}

func getGoogleRegionFromEnvVar() (string, error) {
	return getFirstNonEmptyEnvVar(regionEnvVars)
}
