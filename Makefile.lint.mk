.PHONY: lint-docs
lint-docs:
	scripts/lint-docs.sh

.PHONY: terraform-fmt-check
terraform-fmt-check:
	terraform fmt -check -recursive

.PHONY: lint-go
lint-go:
	golangci-lint run
