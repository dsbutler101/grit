.PHONY: test
test:
	go test -count=1 -v ./...

.PHONY: test-e2e
test-e2e:
	go test -count=1 -tags e2e -v ./...
