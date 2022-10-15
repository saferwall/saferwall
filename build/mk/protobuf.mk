API_DIR = $(ROOT_DIR)/api/protobuf-spec
SERVICES_DIR = $(ROOT_DIR)/services

protobuf-generate-api:		## Generate go code from protobuf spec.
	protoc -I $(API_DIR)/ --go_out=$(SERVICES_DIR)/proto/ \
		$(API_DIR)/message.proto
	mv $(SERVICES_DIR)/proto/github.com/saferwall/saferwall/message.pb.go \
		$(SERVICES_DIR)/proto/message.pb.go
	rm -r $(SERVICES_DIR)/proto/github.com/

PROTOC_VERSION=3.20.3
protobuf-install-compiler: 	## Install protobuf compiler
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-linux-x86_64.zip
	unzip protoc-$(PROTOC_VERSION)-linux-x86_64.zip -d protoc3
	sudo mv protoc3/bin/* /usr/local/bin/
	sudo mv protoc3/include/* /usr/local/include/
	rm protoc-$(PROTOC_VERSION)-linux-x86_64.zip
	rm -r protoc3/
	protoc --version

protobuf-protoc-gen-go:	## Install protoc plugin for Go
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
