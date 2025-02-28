ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
# have a central plugin cache so we don't duplicate in each module
# each module will get symlinks to this
export TF_PLUGIN_CACHE_DIR=$(ROOT_DIR).terraform

TF_MODULES := $(shell find . -name "main.tf" -exec dirname {} \; | sort -u)

.terraform:
	mkdir -p $@

# inits all the modules but uses root .terraform as a cache dir
# useful for CI and before running validate locally
.PHONY: terraform-init
terraform-init: .terraform
terraform-init: $(TF_MODULES:%=%-terraform-init)
%-terraform-init: MODULE=$*
%-terraform-init:
	cd $(MODULE) && terraform init -backend=false

# clean up local terraform caches and lock files
.PHONY: clean
clean:
	find . -type d -name '.terraform' -exec rm -r {} \+
	find . -name '*.terraform.lock.hcl' -exec rm {} \+
