//go:build lint

package variable_test

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_VariableLocations(t *testing.T) {
	defaultRegex, _ := regexp.Compile(`default\s*=`)

	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || filepath.Base(path) != "variables.tf" {
				return nil
			}

			match, _ := regexp.MatchString(".*/(dev|test|prod)/variables.tf", path)
			if !match {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				lineNum := 1
				for scanner.Scan() {
					line := scanner.Text()
					assert.Falsef(t, defaultRegex.MatchString(line), "Default value found in non-allowed directory. Terraform variable defaults are only allowed in dev, prod, and test directories. File: '%s', Line %d: '%s'\n", path, lineNum, line)
					lineNum++
				}

				if err := scanner.Err(); err != nil {
					return err
				}
			}
			return nil
		})

	assert.NoError(t, err)
}
