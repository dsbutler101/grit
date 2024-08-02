//go:build lint

package variable_test

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type variableLocationInvalidLine struct {
	path    string
	lineNum int
	line    string
}

type variableLocationInvalidLines []variableLocationInvalidLine

func (v variableLocationInvalidLines) String() string {
	s := "\n"
	for _, l := range v {
		s += fmt.Sprintf("%s:%d: %s\n", l.path, l.lineNum, l.line)
	}

	return s
}

func Test_VariableLocations(t *testing.T) {
	defaultRegex := regexp.MustCompile(`default\s*=`)
	ignorePathsRE := regexp.MustCompile("(^examples/.*|scenarios/.*|modules/.*/(|dev|test|prod))/variables.tf")

	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// prevent walking . dirs as they can contain third party modules (e.g. CI module cache)
			if path != "." && info.IsDir() && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}

			if info.IsDir() || filepath.Base(path) != "variables.tf" {
				return nil
			}

			if ignorePathsRE.MatchString(path) {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			var il variableLocationInvalidLines

			scanner := bufio.NewScanner(file)
			lineNum := 0
			for scanner.Scan() {
				line := scanner.Text()
				if defaultRegex.MatchString(line) {
					il = append(il, variableLocationInvalidLine{
						path:    path,
						lineNum: lineNum,
						line:    line,
					})
				}

				lineNum++
			}

			assert.Empty(t, il, "Default value found in non-allowed directory. Terraform variable defaults are only allowed in prod and test directories and in the scenarios modules")

			return scanner.Err()
		})

	assert.NoError(t, err)
}
