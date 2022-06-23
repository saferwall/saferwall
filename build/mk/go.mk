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

GO_VERSION = 1.17.11
go-setup:	## Download and install go
	curl -O https://dl.google.com/go/go$(GO_VERSION).linux-amd64.tar.gz
	tar -xvf go$(GO_VERSION).linux-amd64.tar.gz
	sudo rm -rf /usr/local/go
	sudo mv go /usr/local/
	mkdir -p ~/go
	rm go$(GO_VERSION).linux-amd64.tar.gz
	export PATH=$(PATH):/usr/local/go/bin && go version
