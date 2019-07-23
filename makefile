# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
BINARY_NAME=market-patterns
COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html

all: test run

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

run: build
	./$(BINARY_NAME) -yaml-config=app-config.yaml

trunc-load: build
	./$(BINARY_NAME) -start-http-server=false -data-file=data/stocks.zip -company-file=data/nyse-symb-name.csv

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

