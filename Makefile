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
	make env-up
	./integration/run.sh
	make env-down
clean:
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
build-docker:
	make build-linux
	docker build -t gcr.io/rvlt/helloworld:latest .
	make clean
env-up:
	./integration/up.sh
env-down:
	./integration/down.sh