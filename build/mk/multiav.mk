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
	cd $(ROOT_DIR)/core/multiav/$(AV_VENDOR)/proto \
		&& mv multiav.$(AV_VENDOR).pb.go $(AV_VENDOR).pb.go

protobuf-install-compiler: 	## Install protobuf compiler
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
	unzip protoc-3.7.1-linux-x86_64.zip -d protoc3
	sudo mv protoc3/bin/* /usr/local/bin/
	sudo mv protoc3/include/* /usr/local/include/
	rm protoc-3.7.1-linux-x86_64.zip
	rm -r proto3/

protobuf-protoc-gen-go:	## Install protoc plugin for Go
	go get -u github.com/golang/protobuf/protoc-gen-go

compile-multiav-server: protobuf-generate-api	## Compile gRPC server
	go build -ldflags "-s -w" -o $(ROOT_DIR)/build/mk/multiav.$(AV_VENDOR)/bin/server $(ROOT_DIR)/build/mk/multiav.$(AV_VENDOR)/server.go

multiav-build-av:	## build an AV inside a docker contrainer.
	$(eval DOCKER_BUILD_ARGS := "")
ifeq ($(AV_VENDOR),sophos)
	$(eval DOCKER_BUILD_ARGS = "--build-arg SOPHOS_URL=$(SOPHOS_URL)")
endif

ifeq ($(AV_VENDOR),symantec)
	$(eval DOCKER_BUILD_ARGS = "--build-arg SYMANTEC_URL=$(SYMANTEC_URL)")
endif

ifeq ($(AV_VENDOR),eset)
	$(eval DOCKER_BUILD_ARGS = "--build-arg ESET_USER=$(ESET_USER) --build-arg ESET_PWD=$(ESET_PWD)")
endif
	sudo make docker-release ARGS=$(DOCKER_BUILD_ARGS) IMG=$(AV_VENDOR) VERSION=0.0.1 DOCKER_FILE=build/docker/Dockerfile.$(AV_VENDOR) DOCKER_DIR=build/data

multiav-build-av-go: ## Build the AV with the gRPC server
	sudo make docker-release IMG=go$(AV_VENDOR) VERSION=0.0.1 DOCKER_FILE=internal/multiav/$(AV_VENDOR)/Dockerfile DOCKER_DIR=internal/multiav/$(AV_VENDOR)/
