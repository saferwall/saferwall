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

go-setup:	## Download and install go
	# curl -O https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz
	# tar -xvf go1.12.5.linux-amd64.tar.gz
	rm -rf /usr/local/go
	mv ./go /usr/local/
	go version
	go get -u github.com/derekparker/delve/cmd/dlv
	rm go1.12.5.linux-amd64.tar.gz