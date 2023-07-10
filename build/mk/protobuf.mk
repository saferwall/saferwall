API_DIR = $(ROOT_DIR)/api/protobuf-spec
SERVICES_DIR = $(ROOT_DIR)/services
INTERNAL_DIR = $(ROOT_DIR)/internal

protobuf-generate-api:		## Generate go code from protobuf spec.
	protoc -I $(API_DIR)/ --go_out=$(SERVICES_DIR)/proto/ \
		$(API_DIR)/message.proto
	mv $(SERVICES_DIR)/proto/github.com/saferwall/saferwall/message.pb.go \
		$(SERVICES_DIR)/proto/message.pb.go
	rm -r $(SERVICES_DIR)/proto/github.com/

	protoc -I $(API_DIR)/ --go_out=$(INTERNAL_DIR)/agent/proto \
		--go_opt=paths=source_relative --go-grpc_out=$(INTERNAL_DIR)/agent/proto \
		--go-grpc_opt=paths=source_relative $(API_DIR)/agent.proto

PROTOC_VERSION=23.4
protobuf-install-compiler: 	## Install protobuf compiler
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-linux-x86_64.zip
	unzip protoc-$(PROTOC_VERSION)-linux-x86_64.zip -d protoc3
	sudo mv protoc3/bin/* /usr/local/bin/
	sudo mv protoc3/include/* /usr/local/include/
	rm protoc-$(PROTOC_VERSION)-linux-x86_64.zip
	rm -rf protoc3/
	protoc --version

protobuf-protoc-gen-go:	## Install protoc plugin for Go
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
