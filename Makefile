.PHONY: dev-init
dev-init:
	cd environments/dev && terraform init

.PHONY: dev-validate
dev-validate:
	cd environments/dev && terraform validate

