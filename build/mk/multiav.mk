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
	sudo make docker-build ARGS=$(DOCKER_BUILD_ARGS) IMG=$(AV_VENDOR) VERSION=0.0.1 DOCKER_FILE=build/docker/Dockerfile.$(AV_VENDOR) DOCKER_DIR=build/data

multiav-release-av:	multiav-build-av ## release an AV inside a docker contrainer.
	sudo make docker-release ARGS=$(DOCKER_BUILD_ARGS) IMG=$(AV_VENDOR) VERSION=0.0.1 DOCKER_FILE=build/docker/Dockerfile.$(AV_VENDOR) DOCKER_DIR=build/data


multiav-build-av-go: multiav-build-av-go ## Build the AV with the gRPC server
	sudo make docker-build IMG=go$(AV_VENDOR) VERSION=0.0.1 \
	 DOCKER_FILE=build/docker/Dockerfile.go$(AV_VENDOR) \
	 DOCKER_DIR=pkg/grpc/multiav/$(AV_VENDOR)/server


multiav-release-av-go: ## Release the AV with the gRPC server
	sudo make docker-release IMG=go$(AV_VENDOR) VERSION=0.0.1 \
	 DOCKER_FILE=build/docker/Dockerfile.go$(AV_VENDOR) \
	 DOCKER_DIR=pkg/grpc/multiav/$(AV_VENDOR)/server
