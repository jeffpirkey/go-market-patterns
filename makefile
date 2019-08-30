# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
BINARY_NAME=go-market-patterns
COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html
compute?=3

all: test run

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

run: build
	./$(BINARY_NAME) -log-level=INFO

run-memory: build
	./$(BINARY_NAME) -log-level=INFO -trunc-load=true -data-file=data/stocks-small.zip -company-file=data/nyse-symb-name.csv -compute=3

run-mongo: build
	./$(BINARY_NAME) -log-level=INFO -db-connect=mongodb://localhost:27017

trunc-load: build
	./$(BINARY_NAME) -log-level=INFO -trunc-load=true -data-file=data/stocks.zip -company-file=data/nyse-symb-name.csv

train: build
	./$(BINARY_NAME) -log-level=INFO -compute=$(compute)

test:
	$(GOTEST)

test-mongo:
	$(GOTEST) -db-connect=mongodb://localhost:27017 -mongo-db-name=test

cover:
	$(GOTEST) -coverprofile $(COVERAGE_OUT)
	$(GOTOOL) cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_OUT)
	rm -f $(COVERAGE_HTML)

