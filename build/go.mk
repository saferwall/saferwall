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
	go test $(GOPKG)
