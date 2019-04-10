-include build/mk/multiav.avast.mk
-include build/multiav/avira/Makefile
-include build/multiav/bitdefender/Makefile
-include build/multiav/clamav/Makefile
-include build/multiav/comodo/Makefile
-include build/multiav/eset/Makefile
-include build/multiav/fsecure/Makefile
-include build/multiav/kaspersky/Makefile
-include build/multiav/mcafee/Makefile
-include build/multiav/sophos/Makefile
-include build/multiav/symantec/Makefile
-include build/multiav/windows-defender/Makefile

protobuf-generate-api:		## Generates protocol buffers definitions files. 
	protoc -I $(ROOT_DIR)/api/protobuf-spec/ \
		-I${GOPATH}/src \
		--go_out=plugins=grpc:$(ROOT_DIR)/core/multiav/$(AV_VENDOR)/proto/ \
		$(ROOT_DIR)/api/protobuf-spec/multiav.$(AV_VENDOR).proto

protobuf-install-compiler: 	## Install protobuf compiler
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
	unzip protoc-3.7.1-linux-x86_64.zip -d protoc3
	sudo mv protoc3/bin/* /usr/local/bin/
	sudo mv protoc3/include/* /usr/local/include/

protobuf-protoc-gen-go:	## Install protoc plugin for Go
	go get -u github.com/golang/protobuf/protoc-gen-go

compile-multiav-server: protobuf-generate-api	## Compile gRPC server
	go build -ldflags "-s -w" -o $(ROOT_DIR)/build/multiav/$(AV_VENDOR)/bin/server $(ROOT_DIR)/build/multiav/$(AV_VENDOR)/server.go 
