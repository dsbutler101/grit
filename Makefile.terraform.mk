ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

# meta module that references all providers
PROVIDERS_MODULE := $(ROOT_DIR)modules/internal/providers

TF_MODULES := $(shell find . -name "main.tf" -not -path "*/\.*/*" -not -path "$(PROVIDERS_MODULE)/*" -exec dirname {} \; | sort -u)

.PHONY: list-modules
list-modules:
	@echo $(TF_MODULES) | tr ' ' '\n'

.terraform:
	mkdir -p $@

# safely download all modules once
.PHONY: providers-init
providers-init: export TF_PLUGIN_CACHE_DIR=$(ROOT_DIR).terraform
providers-init: .terraform
	cd $(PROVIDERS_MODULE) && terraform init -backend=false

# inits all the modules but uses root .terraform as a cache dir
# cache is created by providers-init target
# useful for CI and before running validate locally
.PHONY: terraform-init
terraform-init: $(TF_MODULES:%=%-terraform-init)
%-terraform-init: MODULE=$*
%-terraform-init: providers-init
	cd $(MODULE) && terraform init -backend=false --plugin-dir=$(ROOT_DIR).terraform

.PHONY: terraform-validate
terraform-validate: $(TF_MODULES:%=%-terraform-validate)
%-terraform-validate: MODULE=$*
%-terraform-validate:
	cd $(MODULE) && terraform validate

# clean up local terraform caches and lock files
.PHONY: clean
clean:
	find . -type d -name '.terraform' -exec rm -r {} \+
	find . -name '*.terraform.lock.hcl' -exec rm {} \+
