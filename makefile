# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
BINARY_NAME=market-patterns
COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

run-load:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME) -load=true

test:
	$(GOTEST)

cover:
	$(GOTEST) -coverprofile $(COVERAGE_OUT)
	$(GOTOOL) cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_OUT)
	rm -f $(COVERAGE_HTML)

