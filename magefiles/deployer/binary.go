package deployer

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/magefile/mage/sh"
	"go.maczukin.dev/libs/mageutils"
	"go.maczukin.dev/libs/version/magefiles"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/magefiles/utils"
)

const (
	submodulePath = "deployer"

	binaryName    = "deployer"
	sha256SumFile = "sha256.sum"

	outputDir = "build"
)

func Compile(platformDefs ...string) error {
	type platformDef struct {
		goos   string
		goarch string
	}

	if len(platformDefs) < 1 {
		platformDefs = []string{fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)}
	}

	var platforms []platformDef

	for _, p := range platformDefs {
		parts := strings.Split(strings.TrimSpace(p), "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid platform definition: %s", p)
		}

		platforms = append(platforms, platformDef{
			goos:   parts[0],
			goarch: parts[1],
		})
	}

	fmt.Println("Establishing version information...")
	ver, err := sh.Output("./ci/version")
	if err != nil {
		return err
	}

	rev := magefiles.Revision()

	gitReference, err := magefiles.GitReference()
	if err != nil {
		return err
	}

	fmt.Printf("... version=%s revision=%s gitReference=%s\n", ver, rev, gitReference)
	fmt.Println("...DONE")

	return onDeployerWD(func() error {
		ldFlags := []string{
			// Set version variables
			fmt.Sprintf("-X %s.version=%s", packageName(), ver),
			fmt.Sprintf("-X %s.revision=%s", packageName(), rev),
			fmt.Sprintf("-X %s.gitReference=%s", packageName(), gitReference),
			fmt.Sprintf("-X %s.builtAt=%s", packageName(), time.Now().UTC().Format(time.RFC3339)),

			// Disable extended logging by default
			fmt.Sprintf("-X %s.customLogFormat=false", packageName()),
			fmt.Sprintf("-X %s.addSources=false", packageName()),

			"-w", "-s",
		}

		var hashes []string
		for _, p := range platforms {
			hash, err := compileForPlatform(ldFlags, p.goos, p.goarch)
			if err != nil {
				return err
			}

			hashes = append(hashes, hash)
		}

		f, err := os.Create(filepath.Join(buildDirectory(), sha256SumFile))
		if err != nil {
			return err
		}
		defer f.Close()

		fmt.Println()
		fmt.Println("Compiled binaries checksums:")
		for _, hash := range hashes {
			_, _ = fmt.Fprintln(f, hash)
			fmt.Println(hash)
		}

		return nil
	})
}

func onDeployerWD(fn func() error) error {
	return utils.OnWD(submodulePath, fn)
}

func compileForPlatform(ldFlags []string, goos string, goarch string) (string, error) {
	fmt.Println()
	fmt.Println("Compiling...")
	fmt.Printf("...for GOOS=%s\n", goos)
	fmt.Printf("...for GOARCH=%s\n", goarch)
	fmt.Println("...using following go tool link flags:")
	for _, flag := range ldFlags {
		fmt.Println("   ", flag)
	}

	envs := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        goos,
		"GOARCH":      goarch,
	}

	inputDir := filepath.Join(workingDirectory(), "cmd", binaryName)
	binaryPath := filepath.Join(buildDirectory(), fmt.Sprintf("%s-%s-%s%s", binaryName, goos, goarch, binaryExtension(goos)))

	fmt.Println("...running go build")
	err := sh.RunWith(envs, "go", "build", "-o", binaryPath, "-ldflags", strings.Join(ldFlags, " "), inputDir)
	if err != nil {
		return "", fmt.Errorf("building binary: %w", err)
	}

	fmt.Println("...DONE")

	binaryHash, err := getBinaryHash(binaryPath)
	if err != nil {
		return "", err
	}

	relBinaryPath := binaryPath
	newPath, err := filepath.Rel(buildDirectory(), binaryPath)
	if err == nil {
		relBinaryPath = newPath
	}

	return fmt.Sprintf("%x  %s", binaryHash, relBinaryPath), nil
}

var packageNameOnce mageutils.Once[string]

func packageName() string {
	return packageNameOnce.MustDo(func() (string, error) {
		out, err := sh.Output("go", "list", ".")
		if err != nil {
			return "", fmt.Errorf("retrieving go package name: %w", err)
		}

		return strings.TrimSpace(out), nil
	})
}

func getBinaryHash(binaryPath string) ([32]byte, error) {
	binaryContent, err := os.ReadFile(binaryPath)
	if err != nil {
		return [32]byte{}, fmt.Errorf("reading binary: %w", err)
	}

	return sha256.Sum256(binaryContent), nil
}

var workingDirectoryOnce mageutils.Once[string]

func workingDirectory() string {
	return workingDirectoryOnce.MustDo(func() (string, error) {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("getting current working directory: %w", err)
		}

		return wd, nil
	})
}

var buildDirectoryOnce mageutils.Once[string]

func buildDirectory() string {
	return buildDirectoryOnce.Do(func() string {
		return filepath.Join(workingDirectory(), outputDir)
	})
}

func binaryExtension(goos string) string {
	if goos == "windows" {
		return ".exe"
	}

	return ""
}
