# GRIT Release Process

This document outlines the processes for creating new GRIT releases including major and minor version updates and patch releases.

## Table of Contents

- [Overview](#overview)
- [Pre-release Tags](#pre-release-tags)
- [Standard Release Process](#standard-release-process)
- [Patch Release Process](#patch-release-process)
- [Troubleshooting](#troubleshooting)

## Overview

GRIT releases are managed through the [Releases](https://gitlab.com/gitlab-org/ci-cd/runner-tools/releases) project,
in combination with configurations in GRIT itself to help facilitate the release process.
It involves updating configuration files in Releases and running the associated GitLab CI/CD jobs to create and publish new releases.

Following the release steps automatically:

- Generates a new `CHANGELOG`.
- Tags and cuts a release on the `main` branch.
- Creates a version branch and merge it into `main`.
- Verifies that MRs that are flagged with a `type::` label. A job fails if any MRs are tagged incorrectly.

## Pre-release Tags

All merges to the `main` branch of GRIT are automatically tagged with a prerelease tag without creating an official release. This action provides versioned references to development builds.

The pre-release tagging follows the following conventions:

- If the current latest release is `vX.Y.Z`, the pre-release tag is for the next minor version (`vX.Y+1.0-pre.N`).
- `N` represents the number of commits after the previous release. For example, if the latest release is `v1.0.0`
  and there have been 10 commits to `main` since that release, the pre-release tag is `v1.1.0-pre.10`

These pre-release tags allow testing development versions and provide a clear trail of code evolution between official releases.

## Standard Release Process

Follow these steps to create a new major or minor release of GRIT:

1. **Update the Release Configuration**

   1. Navigate to the [`config.yml`](https://gitlab.com/gitlab-org/ci-cd/runner-tools/releases/-/blob/main/config.yml) file in the Releases repository.
   1. Update the `releases` section with the new GRIT version information:

      ```yaml
      releases:
        - name: grit
          variant: canonical
          version: vX.Y.0  # Replace X.Y with the new version
          app_version: vA.B.0  # Update with the corresponding GitLab milestone version
      ```

      Example:

      ```yaml
      releases:
        - name: grit
          variant: canonical
          version: v0.12.0
          app_version: v17.10.0
      ```

1. **Create and Merge the MR**

   Create a Merge Request with your changes to the `config.yml` file. Once approved and merged, proceed to the next step.

1. **Run the Releases Pipeline**

   Navigate to the CI/CD Pipelines section of the Releases repository and monitor the pipeline triggered by your merge.

1. **Verify the Dry Run stages and manually run actual changes**

   The pipeline runs several stages in "dry run" mode first. For each stage:

   - Monitor the logs for any issues.
   - Confirm the dry run was successful before running the next job.
   - Check that the correct version is applied.

1. **Verify Release Completion**

   Once all jobs for GRIT are successfully completed:

   - The new version is tagged in the repository.
   - Release artifacts are published.
   - Verifies that the release is visible in the [GRIT releases page](https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit/-/releases).
   - Confirms that the `CHANGELOG` has been properly updated with all the changes since the last release.

## Patch Release Process

For creating patch releases (e.g., v1.0.1, v1.0.2):

1. **Cherry-pick Required Commits**

   - Identify the commits that need to be included in the patch.
   - Cherry-pick these commits into the stable branch for the current version.

   For example, if patching version v1.0.x:

   ```shell
   git checkout 1-0-stable
   git cherry-pick <commit-hash>
   git push
   ```

   > **Important:** Ensure commits are cherry-picked in the correct order to maintain Git history integrity.

1. **Update the Release Configuration**

   1. Navigate to the [config](https://gitlab.com/gitlab-org/ci-cd/runner-tools/releases/-/blob/main/config.yml) file in the Releases repository.
   1. Add a new entry for the patch version:

      ```yaml
      releases:
        - name: grit
          variant: canonical
          version: vX.Y.Z  # Replace X.Y.Z with the new patch version
          app_version: vA.B.C  # Update with the corresponding GitLab milestone version
      ```

   > **Important:** Check existing GRIT releases to determine the next patch version. The version will not auto-increment.

1. **Create and Merge the MR**

   Create a Merge Request with your changes to the `config.yml` file. Once approved and merged, proceed with the release pipeline.

1. **Follow the Pipeline Process**

   Follow steps 3 to 6 from the standard release process. The pipeline automatically detects the stable branch (for example, `1-0-stable`) and apply the changes for the patch release.

## Troubleshooting

If you encounter issues during the release process:

- Check the pipeline logs for specific error messages.
- Ensure the version numbers are correctly formatted.
- Confirm that the stable branch exists (for patch releases).
- Verify that all required merge request labels are properly applied.
