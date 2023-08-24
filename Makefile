.PHONY: dev-init
dev-init:
	cd environments/dev && terraform init

.PHONY: dev-validate
dev-validate:
	cd environments/dev && terraform validate

.PHONY: dev-plan
dev-plan:
	cd environments/dev && terraform plan

.PHONY: dev-apply
dev-apply:
	cd environments/dev && terraform apply

.PHONY: dev-output
dev-output:
	cd environments/dev && terraform output
