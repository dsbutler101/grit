export VERSION := $(shell ./ci/version)
REVISION := $(shell git rev-parse --short=8 HEAD || echo unknown)
BRANCH := $(shell git show-ref | grep "$(REVISION)" | grep -v HEAD | awk '{print $$2}' | sed 's|refs/remotes/origin/||' | sed 's|refs/heads/||' | sort | head -n 1)
LATEST_STABLE_TAG := $(shell git -c versionsort.prereleaseSuffix="-pre" tag -l "v*.*.*" | sort -rV | awk '!/rc/' | head -n 1)

.PHONY: version
version:
	@echo Current version: $(VERSION)
	@echo Current revision: $(REVISION)
	@echo Current branch: $(BRANCH)
	@echo Latest stable: $(LATEST_STABLE_TAG)
