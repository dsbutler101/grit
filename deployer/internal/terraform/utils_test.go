package terraform

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadRequiredStringTo(t *testing.T) {
	testStringKey := "string-key"
	testString := "string-value"
	testEmptyStringKey := "empty-string-key"
	testIntKey := "int-key"
	testRMMap := map[string]any{
		testStringKey:      testString,
		testEmptyStringKey: "",
		testIntKey:         0,
	}

	tests := map[string]struct {
		key         string
		assertError func(t *testing.T, err error)
	}{
		"missing required entry": {
			key: "unknown-key",
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errMissingRequiredEntry)
			},
		},
		"empty required entry": {
			key: testEmptyStringKey,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errRequiredEntryIsEmpty)
			},
		},
		"entry is not a string": {
			key: testIntKey,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errIsNotAString)
			},
		},
		"success": {
			key: testStringKey,
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			var target string
			err := readRequiredStringTo(testRMMap, tt.key, &target)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testString, target)
		})
	}
}

func TestReadOptionalStringTo(t *testing.T) {
	testStringKey := "string-key"
	testString := "string-value"
	testEmptyStringKey := "empty-string-key"
	testIntKey := "int-key"
	testRMMap := map[string]any{
		testStringKey:      testString,
		testEmptyStringKey: "",
		testIntKey:         0,
	}

	tests := map[string]struct {
		key         string
		assertError func(t *testing.T, err error)
		expectedVal string
	}{
		"missing optional entry": {
			key:         "unknown-key",
			expectedVal: "",
		},
		"empty optional entry": {
			key:         testEmptyStringKey,
			expectedVal: "",
		},
		"entry is not a string": {
			key: testIntKey,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errIsNotAString)
			},
		},
		"success": {
			key:         testStringKey,
			expectedVal: testString,
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			var target string
			err := readOptionalStringTo(testRMMap, tt.key, &target)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedVal, target)
		})
	}
}

func TestAssertWorkdir(t *testing.T) {
	testPath := t.TempDir()
	testRelativePath := "path-1"

	tests := map[string]struct {
		path         string
		changeWD     func(t *testing.T)
		expectedPath string
	}{
		"wd is an absolute path": {
			path:         testPath,
			expectedPath: testPath,
		},
		"wd is a relative path": {
			path: testRelativePath,
			changeWD: func(t *testing.T) {
				curDir, err := os.Getwd()

				require.NoError(t, err)
				require.NoError(t, os.Chdir(testPath))

				t.Cleanup(func() {
					require.NoError(t, os.Chdir(curDir))
				})
			},
			expectedPath: filepath.Join(testPath, testRelativePath),
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			if tt.changeWD != nil {
				tt.changeWD(t)
			}

			path, err := assertWorkdir(tt.path)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPath, path)
		})
	}
}
