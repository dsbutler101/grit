.PHONY: test
test:
ifdef CI
	go run gotest.tools/gotestsum@latest --junitfile=junit.xml -- -count=1 ./...
else
	go run gotest.tools/gotestsum@latest -- -count=1 ./...
endif
