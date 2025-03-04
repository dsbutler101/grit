package deployer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/magefile/mage/sh"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/magefiles/mockery"
)

var (
	mockFileRx = regexp.MustCompile(`^mock_.*\.go$`)
)

func GoGenerate() error {
	err := mockery.Install()
	if err != nil {
		return fmt.Errorf("installing mockery: %w", err)
	}

	fmt.Println()
	fmt.Println("Generating mock files...")

	return onDeployerWD(func() error {
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if !mockFileRx.MatchString(info.Name()) {
				return nil
			}

			fmt.Printf("... deleting %s\n", path)
			return sh.Rm(path)
		})
		if err != nil {
			return fmt.Errorf("walking directory: %w", err)
		}

		fmt.Println("... generating new mock files with go:generate")

		err = sh.RunV("go", "generate", "./...")
		if err != nil {
			return fmt.Errorf("generating mock files: %w", err)
		}

		fmt.Println("... DONE")

		return nil
	})
}
