#######################
# Version definitions #
#######################

GO_VERSION ?= 1.23.6

#####################
# CI Registry setup #
#####################

CI_REGISTRY ?= registry.gitlab.com
CI_PROJECT_PATH ?= gitlab-org/ci-cd/runner-tools/grit
CI_REGISTRY_IMAGE ?= $(CI_REGISTRY)/$(CI_PROJECT_PATH)
CI_IMAGE ?= $(CI_REGISTRY_IMAGE)/ci:latest

.PHONY: build-ci-image
build-ci-image: IMAGE ?= $(CI_IMAGE)
build-ci-image:
ifdef CI_REGISTRY_USER
	# Logging into $(CI_REGISTRY)
	@docker login --username $(CI_REGISTRY_USER) --password $(CI_REGISTRY_PASSWORD) $(CI_REGISTRY)
	docker pull $(IMAGE) || echo "Remote image $(IMAGE) not available - will not use cache"
endif
	docker build \
		--cache-from $(IMAGE) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		-t $(IMAGE) \
		-f ./dockerfiles/ci/Dockerfile \
		.
ifdef CI_REGISTRY_USER
	# Pushing $(IMAGE)
	@docker push $(IMAGE)
	# Logging out from $(CI_REGISTRY)
	@docker logout $(CI_REGISTRY)
endif
