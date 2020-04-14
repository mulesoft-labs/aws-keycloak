GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=aws-keycloak

.PHONY: all
all: test build

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

.PHONY: vendor
vendor:
	$(GOMOD) vendor

.PHONY: tidy
tidy:
	$(GOMOD) tidy

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
