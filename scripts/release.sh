#!/usr/bin/env bash

# release.sh - Script to create a new stable release for GRIT
# Usage: ./release.sh <tag_name> <artifacts_list_file>

set -euo pipefail

# Check for required arguments
if [ $# -lt 2 ]; then
  echo "Usage: $0 <tag_name> <artifacts_list_file>"
  echo "Example: $0 v1.2.3 artifacts_list.json"
  exit 1
fi

TAG_NAME=$1
ARTIFACTS_LIST_FILE=$2

# Verify tag name is provided and valid
if [[ ! $TAG_NAME =~ ^v[0-9]+\.[0-9]+\.[0-9]+ ]]; then
  echo "Error: Tag name should follow semantic versioning pattern (e.g., v1.2.3)"
  exit 1
fi

# Verify artifacts list file exists
if [ ! -f "$ARTIFACTS_LIST_FILE" ]; then
  echo "Error: Artifacts list file '$ARTIFACTS_LIST_FILE' not found"
  exit 1
fi

# Set variables
CHANGELOG_URL="$CI_PROJECT_URL/-/blob/$TAG_NAME/CHANGELOG.md"
PROJECT_ID=$CI_PROJECT_ID

echo "Releasing new version $TAG_NAME"

# Check for required dependencies
if ! command -v jq &> /dev/null; then
  echo "Error: jq is required but not installed. Please install it before running this script."
  exit 1
fi

if ! command -v glab &> /dev/null; then
  echo "Error: glab is required but not installed."
  echo "You can get it from: https://gitlab.com/gitlab-org/cli"
  exit 1
fi

# Create the description string
DESCRIPTION="See [the changelog](${CHANGELOG_URL}) :rocket:

GitLab Runner Infrastructure Toolkit (GRIT) documentation can be found at $CI_PROJECT_URL/-/blob/main/README.md"

# Create a JSON array of assets
ASSETS_ARGS=$(jq -c '[.[] | {"name": .file_name, "url": .web_url, "direct_asset_path": "/binaries/\(.file_name)"}]' < "$ARTIFACTS_LIST_FILE")

# Run the release command directly
glab release create "$TAG_NAME" \
  --name "$TAG_NAME" \
  --notes "$DESCRIPTION" \
  --assets-links "$ASSETS_ARGS"

echo "Release $TAG_NAME created successfully!"