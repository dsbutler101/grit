package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/magefile/mage/sh"
	"go.maczukin.dev/libs/mageutils"
)

const (
	ciAPIUrlEnv     = "CI_API_V4_URL"
	ciProjectIDEnv  = "CI_PROJECT_ID"
	ciJobTokenEnv   = "CI_JOB_TOKEN"
	ciProjectURLEnv = "CI_PROJECT_URL"

	jobTokenHeader = "Job-Token"

	selectQueryParameter = "select"
	packageFileSelect    = "package_file"

	artifactsListFileName = "artifacts_list.json"
)

type packageFileSelectResponse struct {
	ID   int64 `json:"id"`
	Size int64 `json:"size"`
}

type UploadResult struct {
	FileName  string `json:"file_name"`
	WebURL    string `json:"web_url"`
	SizeBytes int64  `json:"size_bytes"`
}

func Upload(ctx context.Context) error {
	fmt.Println("Establishing version information...")
	ver, err := sh.Output("./ci/version")
	if err != nil {
		return err
	}
	fmt.Println("...DONE")

	var results []UploadResult

	err = onDeployerWD(func() error {
		binaryRx, err := regexp.Compile(fmt.Sprintf("^%s-.*-.*(.exe)?$", binaryName))
		if err != nil {
			return fmt.Errorf("compiling binaryRegexPattern: %w", err)
		}

		return filepath.Walk(buildDirectory(), func(path string, info os.FileInfo, _ error) error {
			if !binaryRx.MatchString(info.Name()) && info.Name() != sha256SumFile {
				return nil
			}

			fmt.Printf("Uploading %s...\n", info.Name())
			ur, err := upload(ctx, path, info.Name(), ver)
			if err != nil {
				return fmt.Errorf("uploading: %w", err)
			}
			fmt.Println("...DONE")

			results = append(results, ur)

			return nil
		})
	})

	if err != nil {
		return fmt.Errorf("uploading: %w", err)
	}

	artifactsListFile, err := os.Create(artifactsListFileName)
	if err != nil {
		return fmt.Errorf("creating artifacts list file: %w", err)
	}
	defer artifactsListFile.Close()

	err = json.NewEncoder(artifactsListFile).Encode(results)
	if err != nil {
		return fmt.Errorf("writing artifacts list file: %w", err)
	}

	return nil
}

func upload(ctx context.Context, path string, fileName string, version string) (UploadResult, error) {
	var ur UploadResult

	file, err := os.Open(path)
	if err != nil {
		return ur, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	u, err := uploadURL(fileName, version)
	if err != nil {
		return ur, fmt.Errorf("creating url: %w", err)
	}

	fmt.Println("- target:", u)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u, file)
	if err != nil {
		return ur, fmt.Errorf("creating upload request: %w", err)
	}

	req.Header.Set(jobTokenHeader, getCIJobToken())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ur, fmt.Errorf("sending upload request: %w", err)
	}

	fmt.Println("  status:", resp.Status)

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return ur, fmt.Errorf("reading upload response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("  response:", string(b))

		return ur, fmt.Errorf("uploading failed: status %s", resp.Status)
	}

	var r packageFileSelectResponse
	err = json.Unmarshal(b, &r)
	if err != nil {
		return ur, fmt.Errorf("parsing upload response: %w", err)
	}

	id := strconv.FormatInt(r.ID, 10)

	webURLParts := []string{
		getCIProjectURL(),
		"-",
		"package_files",
		id,
		"download",
	}

	ur.FileName = fileName
	ur.WebURL = fmt.Sprint(strings.Join(webURLParts, "/"))
	ur.SizeBytes = r.Size

	fmt.Println("  web_url:", ur.WebURL)
	fmt.Println("  size bytes:", ur.SizeBytes)

	return ur, nil
}

func uploadURL(fileName string, version string) (string, error) {
	u := fmt.Sprintf(
		"%s/projects/%s/packages/generic/%s/%s/%s",
		getCIAPIUrl(),
		getCIProjectID(),
		binaryName,
		version,
		fileName,
	)

	reqURL, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("parsing generated url: %w", err)
	}

	q := make(url.Values)
	q.Set(selectQueryParameter, packageFileSelect)
	reqURL.RawQuery = q.Encode()

	return reqURL.String(), nil
}

var ciAPIUrl mageutils.Once[string]

func getCIAPIUrl() string {
	return ciAPIUrl.MustDo(func() (string, error) {
		u := os.Getenv(ciAPIUrlEnv)
		if u == "" {
			return "", fmt.Errorf("%s not defined", ciAPIUrlEnv)
		}

		return u, nil
	})
}

var ciProjectID mageutils.Once[string]

func getCIProjectID() string {
	return ciProjectID.MustDo(func() (string, error) {
		i := os.Getenv(ciProjectIDEnv)
		if i == "" {
			return "", fmt.Errorf("%s not defined", ciProjectIDEnv)
		}

		return i, nil
	})
}

var ciJobToken mageutils.Once[string]

func getCIJobToken() string {
	return ciJobToken.MustDo(func() (string, error) {
		t := os.Getenv(ciJobTokenEnv)
		if t == "" {
			return "", fmt.Errorf("%s not defined", ciJobTokenEnv)
		}

		return t, nil
	})
}

var ciProjectURL mageutils.Once[string]

func getCIProjectURL() string {
	return ciProjectURL.MustDo(func() (string, error) {
		u := os.Getenv(ciProjectURLEnv)
		if u == "" {
			return "", fmt.Errorf("%s not defined", ciProjectURLEnv)
		}

		return u, nil
	})
}
