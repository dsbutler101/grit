package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/magefile/mage/sh"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/magefiles/utils"
)

const (
	gotestsumVersion        = "v1.12.0"
	gocoverCoberturaVersion = "v1.3.0"

	coverageOutFileName = "coverage.txt"
	coverageXMLFileName = "coverage.xml"
)

var (
	coverageTotalRx = regexp.MustCompile(`total:\t+\(statements\)\t+([0-9]+\.[0-9]+%)`)
)

func UnitForPath(path string) error {
	return utils.OnWD(path, func() error {
		return Unit()
	})
}

func Unit() error {
	var testsRunErr error

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd: %w", err)
	}

	runArgs := []string{
		"run",
		fmt.Sprintf("gotest.tools/gotestsum@%s", gotestsumVersion),
	}

	if utils.IsCI() {
		runArgs = append(runArgs, "--junitfile=junit.xml")
	}

	coverageOutFilePath := filepath.Join(wd, coverageOutFileName)
	testArgs := []string{
		"--",
		"-count=1",
		fmt.Sprintf("-coverprofile=%s", coverageOutFilePath),
		"-covermode=count",
		"./...",
	}

	args := append(runArgs, testArgs...)

	err = sh.RunV("go", args...)
	if err != nil {
		testsRunErr = fmt.Errorf("runing tests: %w", err)
	}

	err = generateCoverageReportFile(wd)
	if err != nil {
		err = fmt.Errorf("running gocover-cobertura: %w", err)
	}

	if testsRunErr != nil {
		if err != nil {
			fmt.Println(err)
		}

		return testsRunErr
	}

	return err
}

func generateCoverageReportFile(wd string) error {
	coverageOutFilePath := filepath.Join(wd, coverageOutFileName)
	in, err := os.Open(coverageOutFilePath)
	if err != nil {
		return fmt.Errorf("opening coverage output file %s: %w", coverageOutFilePath, err)
	}

	defer in.Close()

	coverageOutXMLFilePath := filepath.Join(wd, coverageXMLFileName)
	out, err := os.Create(coverageOutXMLFilePath)
	if err != nil {
		return fmt.Errorf("creating coverage output XML file %s: %w", coverageOutXMLFilePath, err)
	}

	defer out.Close()

	cmd := exec.Command(
		"go", "run",
		fmt.Sprintf("github.com/boumenot/gocover-cobertura@%s", gocoverCoberturaVersion),
		"-ignore-gen-files",
	)
	cmd.Stderr = os.Stderr
	cmd.Stdin = in
	cmd.Stdout = out

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("running gocover-cobertura: %w", err)
	}

	output, err := sh.Output("go", "tool", "cover", fmt.Sprintf("-func=%s", coverageOutFilePath))
	if err != nil {
		return fmt.Errorf("generating goverage summary: %w", err)
	}

	coverageTotal := coverageTotalRx.FindStringSubmatch(output)
	fmt.Printf("\nCoverage TOTAL: %s\n", coverageTotal[1])

	return nil
}
