.PHONY: test
test:
	go test -v ./...

.PHONY: test-e2e
test-e2e:
	go test -tags e2e ./...