API_DIR = $(ROOT_DIR)/api/protobuf-spec
SERVICES_DIR = $(ROOT_DIR)/services

protobuf-generate-api:		## Generate go code from protobuf spec.
	protoc -I $(API_DIR)/ --go_out=$(SERVICES_DIR)/proto/ \
		$(API_DIR)/message.proto

protobuf-install-compiler: 	## Install protobuf compiler
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
	unzip protoc-3.7.1-linux-x86_64.zip -d protoc3
	sudo mv protoc3/bin/* /usr/local/bin/
	sudo mv protoc3/include/* /usr/local/include/
	rm protoc-3.7.1-linux-x86_64.zip
	rm -r protoc3/

protobuf-protoc-gen-go:	## Install protoc plugin for Go
	go get -u github.com/golang/protobuf/protoc-gen-go
