BINARY_NAME=go-sandbox
COVERAGE_FILE=coverage.out
HTML_REPORT=coverage.html

## all: builds and tests
.PHONY: all
all: build test

## build: creates the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./...

## clean: Deletes all artefacts
.PHONY: clean
clean:
	@echo "Cleaning artefacts..."
	@go clean
	@rm -f $(BINARY_NAME) $(COVERAGE_FILE) $(HTML_REPORT)

## rebuild: Deletes all artefacts and then builds
.PHONY: rebuild
rebuild: clean build
	@echo "Force rebuilding..."

## test: just a simple test with verbose
.PHONY: test
test:
	@go test -v ./...

## testcov: Tests with coverage and shows result in the terminal
.PHONY: testcov
testcov:
	@go test -v -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -func=$(COVERAGE_FILE)

## testcov-html: Generate coverage in HTML and opens in the browser (assumes MacOS)
.PHONY: testcov-html
testcov-html: testcov
	@go tool cover -html=$(COVERAGE_FILE) -o $(HTML_REPORT)
	@echo "Coverage report generated at $(HTML_REPORT)"
	@open $(HTML_REPORT)
