# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=helloworld
CMD_PATH=./cmd/$(BINARY_NAME)
    
build:
	$(GOBUILD) -o $(BINARY_NAME) $(CMD_PATH)
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) $(CMD_PATH)
test:
	$(GOTEST) -v ./...
integration-tests:
	./integration/run.sh
clean:
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
build-docker:
	make build-linux
	docker build -t rvlt/helloworld:latest .
	make clean
