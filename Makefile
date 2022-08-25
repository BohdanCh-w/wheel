all: fmt lint test

# Check code format
fmt:
	@echo "Checking formatting of go code..."
	@result=$$(gofmt -d -l -e . 2>&1); \
		if [ "$$result" ]; then \
			echo "$$result"; \
			echo "gofmt failed!"; \
			exit 1; \
		fi

# Check linters
lint:
	@echo "Check linters..."
	@PATH=$(PATH):$(PWD)/bin golangci-lint version
	@PATH=$(PATH):$(PWD)/bin golangci-lint run -c .golangci.yml ./...

# Run tests
test:
	@echo "Run tests..."
	@go test -v -race -coverprofile coverage/test.out -covermode=atomic ./...;
	@go tool cover -func coverage/test.out
	@go tool cover -html coverage/test.out -o coverage/gocover.html
	@PATH=$(PATH):$(PWD)/bin gocover-cobertura < coverage/test.out > coverage/cobertura.xml;

install-gocover:
	@echo "Download gocover-cobertura..."
	@GOBIN=$(PWD)/bin go install github.com/boumenot/gocover-cobertura@v1.2.0

install-linters:
	@echo "Download golangci-lint ..."
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.49.0
