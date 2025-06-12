package runner_version_lookup

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

// Instead of a struct, we now use a slice of strings
type manifestArray []string

func loadManifest(t *testing.T) manifestArray {
	data, err := os.ReadFile("manifest.json")
	require.NoError(t, err)
	var versions manifestArray
	err = json.Unmarshal(data, &versions)
	require.NoError(t, err)
	return versions
}

// Helper function to create maps on the fly, same as in Terraform
func createMaps(versions manifestArray) (map[string]string, map[string]string) {
	skewToVersion := make(map[string]string)
	versionToSkew := make(map[string]string)
	
	for idx, version := range versions {
		skewStr := strconv.Itoa(idx)
		skewToVersion[skewStr] = version
		versionToSkew[version] = skewStr
	}
	
	return skewToVersion, versionToSkew
}

func TestRunnerVersionLookupManifest(t *testing.T) {
	versions := loadManifest(t)
	// Manifest should be an array of versions where the index is the skew
	require.Len(t, versions, 3, "Manifest should contain 3 versions")
	
	// Verify all array elements are valid version strings
	for idx, version := range versions {
		require.NotEmpty(t, version, "Version at index %d should not be empty", idx)
	}
	
	// Create the maps for testing
	skewToVersion, versionToSkew := createMaps(versions)
	
	// Verify bi-directional mapping works correctly
	for skew, version := range skewToVersion {
		reverseSkew, exists := versionToSkew[version]
		require.True(t, exists, "Version %q should map back to a skew", version)
		require.Equal(t, skew, reverseSkew, "Skew %q must map to a version that maps back to it", skew)
	}
}

func TestRunnerVersionLookupModuleSuccess(t *testing.T) {
	versions := loadManifest(t)
	skewToVersion, versionToSkew := createMaps(versions)
	
	name := test_tools.JobName(t)

	for skew, version := range skewToVersion {
		t.Run(skew, func(t *testing.T) {
			i, err := strconv.Atoi(skew)
			require.NoError(t, err)
			moduleVars := map[string]any{
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{},
					"min_support": "experimental",
				},
				"skew": float64(i),
			}
			expectedOutputs := map[string]any{
				"skew":           float64(i),
				"runner_version": version,
			}
			test_tools.ApplyAndAssertOutputs(t, moduleVars, expectedOutputs)
		})
	}
	for version, skew := range versionToSkew {
		t.Run(version, func(t *testing.T) {
			i, err := strconv.Atoi(skew)
			require.NoError(t, err)
			moduleVars := map[string]any{
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{},
					"min_support": "experimental",
				},
				"runner_version": version,
			}
			expectedOutputs := map[string]any{
				"skew":           float64(i),
				"runner_version": version,
			}
			test_tools.ApplyAndAssertOutputs(t, moduleVars, expectedOutputs)
		})
	}
}

func TestRunnerVersionLookupModuleError(t *testing.T) {
	versions := loadManifest(t)
	skewToVersion, _ := createMaps(versions)
	skew := float64(0)
	version := skewToVersion["0"]
	
	name := test_tools.JobName(t)

	testCases := map[string]struct {
		name       string
		moduleVars map[string]any
	}{
		"both skew and version provided": {
			moduleVars: map[string]any{
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{},
					"min_support": "experimental",
				},
				"skew":           skew,
				"runner_version": version,
			},
		},
		"neither skew nor version provided": {
			moduleVars: map[string]any{
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{},
					"min_support": "experimental",
				},
			},
		},
		"skew out of bounds": {
			moduleVars: map[string]any{
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{},
					"min_support": "experimental",
				},
				"skew": float64(-1),
			},
		},
		"version in wrong format": {
			moduleVars: map[string]any{
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{},
					"min_support": "experimental",
				},
				"runner_version": "v1.2.3", // should not have 'v'
			},
		},
		"unsupported version without allow flag": {
			moduleVars: map[string]any{
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{},
					"min_support": "experimental",
				},
				"runner_version": "16.8.0", // Not in the manifest
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			test_tools.ApplyAndAssertError(t, tc.moduleVars, true)
		})
	}
}

func TestRunnerVersionLookupUnsupportedVersionsWithAllowFlag(t *testing.T) {
	// Test that unsupported versions work when allow_unsupported_versions = true
	unsupportedVersion := "16.8.0" // A version not in the manifest
	
	name := test_tools.JobName(t)
	
	moduleVars := map[string]any{
		"metadata": map[string]interface{}{
			"name":        name,
			"labels":      map[string]string{},
			"min_support": "experimental",
		},
		"runner_version":             unsupportedVersion,
		"allow_unsupported_versions": true,
	}
	
	expectedOutputs := map[string]any{
		"skew":           float64(-1), // Skew should be -1 for unsupported versions
		"runner_version": unsupportedVersion,
	}
	
	test_tools.ApplyAndAssertOutputs(t, moduleVars, expectedOutputs)
}