-include build/mk/multiav.avast.mk
-include build/mk/multiav.avira.mk
-include build/mk/multiav.bitdefender.mk
-include build/mk/multiav.clamav.mk
-include build/mk/multiav.comodo.mk
-include build/mk/multiav.eset.mk
-include build/mk/multiav.fsecure.mk
-include build/mk/multiav.kaspersky.mk
-include build/mk/multiav.mcafee.mk
-include build/mk/multiav.sophos.mk
-include build/mk/multiav.symantec.mk
-include build/mk/multiav.windefender.mk

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
	go build -ldflags "-s -w" -o $(ROOT_DIR)/build/mk/multiav.$(AV_VENDOR)/bin/server $(ROOT_DIR)/build/mk/multiav.$(AV_VENDOR)/server.go 
