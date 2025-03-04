TF_MODULES := $(shell find . -name "main.tf" -exec dirname {} \; | sort -u)

# inits all the modules but uses root .terraform as a cache dir
# useful for CI and before running validate locally
.PHONY: terraform-init
terraform-init: $(TF_MODULES:%=%-terraform-init)
%-terraform-init: MODULE=$*
%-terraform-init:
	cd $(MODULE) && terraform init -backend=false

# clean up local terraform caches and lock files
.PHONY: clean
clean:
	find . -type d -name '.terraform' -exec rm -r {} \+
	find . -name '*.terraform.lock.hcl' -exec rm {} \+
