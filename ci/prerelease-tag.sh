#!/bin/bash
# Script to create a prerelease tag

echo "Checking version consistency..."
file_version=v$(cat VERSION)
repo_version=$(git describe --abbrev=0 --tags --exclude='*-pre*' 2>/dev/null || echo "0.0.0")
if [ "$file_version" != "$repo_version" ]; then
    echo "Error: Repository version ($repo_version) doesn't match VERSION file ($file_version)"
    exit 1
fi

git fetch origin main
local_commit=$(git rev-parse HEAD)
remote_commit=$(git rev-parse origin/main)
if [ "$local_commit" != "$remote_commit" ]; then
    echo "Error: Local main ($local_commit) is not in sync with remote main ($remote_commit)"
    exit 1
fi

echo "Configuring git user..."
git config --global user.email "auto-runner-releaser@gitlab.com"
git config --global user.name "Auto Runner Releaser"

PRERELEASE_TAG=$(./ci/version)
echo "Creating tag: $PRERELEASE_TAG"
git tag -a $PRERELEASE_TAG -m "Prerelease $PRERELEASE_TAG"
git push https://oauth2:$GRIT_RELEASE_GITLAB_TOKEN@gitlab.com/gitlab-org/ci-cd/runner-tools/grit.git $PRERELEASE_TAG