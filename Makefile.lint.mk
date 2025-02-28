.PHONY: lint-docs
lint-docs:
	scripts/lint-docs.sh

.PHONY: lint-go
lint-go:
	golangci-lint run

.PHONY: terraform-validate
terraform-validate: $(TF_MODULES:%=%-terraform-validate)
%-terraform-validate: MODULE=$*
%-terraform-validate:
	cd $(MODULE) && terraform validate

.PHONY: terraform-fmt-check
terraform-fmt-check:
	terraform fmt -check -recursive -diff

.PHONY: tflint
tflint:
	tflint --recursive --config "$$(pwd)/.tflint.hcl"
