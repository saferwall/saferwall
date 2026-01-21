GOFLAGS :=  -ldflags -s -w

go-compile: go-clean go-get go-build

go-clean:	## Remove object files and cached files
	go clean

go-get:		## Download and install packages and dependencies
	go get $(get)

go-install:	## Compile and install packages and dependencies
	go install $(GOFILES)

go-build:	## Compile packages and dependencies
	go build $(GOFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

go-test:	## Test packages
	go test -v $(GOPKG)

go-lint:	## Test packages
	staticcheck ./...

GO_VERSION = 1.24.8
go-setup:	## Download and install go
	curl -O https://dl.google.com/go/go$(GO_VERSION).linux-amd64.tar.gz
	tar -xvf go$(GO_VERSION).linux-amd64.tar.gz
	sudo rm -rf /usr/local/go
	sudo mv go /usr/local/
	mkdir -p ~/go
	rm go$(GO_VERSION).linux-amd64.tar.gz
	export PATH=$(PATH):/usr/local/go/bin && go version

go-install-staticcheck: ## Install staticheck.
	go install honnef.co/go/tools/cmd/staticcheck@latest

go-install-golangci-lint:	## Install golangci-lint.
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.1
