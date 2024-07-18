.PHONY: lint-docs
lint-docs:
	scripts/lint-docs.sh

.PHONY: terraform-fmt-check
terraform-fmt-check:
	terraform fmt -check -recursive

.PHONY: lint-terraform
lint-terraform:
	go test -tags lint .