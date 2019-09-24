API_DIR = $(ROOT_DIR)/api/protobuf-spec
AV_LIST = $(ROOT_DIR)/core/multiav

protobuf-generate-api:		## Generates protocol buffers definitions files. 
	protoc -I $(API_DIR)/ \
		--go_out=plugins=grpc:$(ROOT_DIR)/core/multiav/$(AV_VENDOR)/proto/ \
		$(API_DIR)/multiav.$(AV_VENDOR).proto
	cd $(ROOT_DIR)/core/multiav/$(AV_VENDOR)/proto \
		&& mv multiav.$(AV_VENDOR).pb.go $(AV_VENDOR).pb.go

protobuf-generate-api-all:	## Generates protocol buffers definitions files for all AVs.
	for AV_VENDOR in $(shell ls $(AV_LIST)) ; do \
		 protoc -I $(API_DIR)/ --go_out=plugins=grpc:$(ROOT_DIR)/core/multiav/$$AV_VENDOR/proto/ $(API_DIR)/multiav.$$AV_VENDOR.proto ; \
		 cd $(ROOT_DIR)/core/multiav/$$AV_VENDOR/proto && mv multiav.$$AV_VENDOR.pb.go $$AV_VENDOR.pb.go ; \
	done
protobuf-install-compiler: 	## Install protobuf compiler
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
	unzip protoc-3.7.1-linux-x86_64.zip -d protoc3
	sudo mv protoc3/bin/* /usr/local/bin/
	sudo mv protoc3/include/* /usr/local/include/
	rm protoc-3.7.1-linux-x86_64.zip
	rm -r proto3/

protobuf-protoc-gen-go:	## Install protoc plugin for Go
	go get -u github.com/golang/protobuf/protoc-gen-go


