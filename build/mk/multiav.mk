-include build/mk/multiav.avast.mk
-include build/mk/multiav.avira.mk
-include build/mk/multiav.bitdefender.mk
-include build/mk/multiav.clamav.mk
-include build/mk/multiav.comodo.mk
-include build/mk/multiav.drweb.mk
-include build/mk/multiav.eset.mk
-include build/mk/multiav.fsecure.mk
-include build/mk/multiav.kaspersky.mk
-include build/mk/multiav.mcafee.mk
-include build/mk/multiav.sophos.mk
-include build/mk/multiav.symantec.mk
-include build/mk/multiav.windefender.mk


compile-multiav-server: protobuf-generate-api	## Compile gRPC server
	go build -ldflags "-s -w" -o $(ROOT_DIR)/build/mk/multiav.$(AV_VENDOR)/bin/server \
		$(ROOT_DIR)/build/mk/multiav.$(AV_VENDOR)/server.go

multiav-build-av:	## build an AV inside a docker contrainer.
	$(eval DOCKER_BUILD_ARGS := "")
ifeq ($(AV_VENDOR),sophos)
	$(eval DOCKER_BUILD_ARGS = "--build-arg SOPHOS_URL=$(SOPHOS_URL)")
endif

ifeq ($(AV_VENDOR),symantec)
	$(eval DOCKER_BUILD_ARGS = "--build-arg SYMANTEC_URL=$(SYMANTEC_URL)")
endif

ifeq ($(AV_VENDOR),eset)
	$(eval DOCKER_BUILD_ARGS = "--build-arg ESET_LICENSE_KEY=$(ESET_LICENSE_KEY)")
endif

ifeq ($(AV_VENDOR),drweb)
	$(eval DOCKER_BUILD_ARGS = "--build-arg DR_WEB_LICENSE_KEY=$(DR_WEB_LICENSE_KEY)")
endif

ifeq ($(AV_VENDOR),bitdefender)
	$(eval DOCKER_BUILD_ARGS = "--build-arg BITDEFENDER_URL=$(BITDEFENDER_URL) --build-arg BITDEFENDER_LICENSE_KEY=$(BITDEFENDER_LICENSE_KEY)")
endif

ifeq ($(AV_VENDOR),trendmicro)
	$(eval DOCKER_BUILD_ARGS = "--build-arg TREND_MICRO_LICENSE_KEY=$(TREND_MICRO_LICENSE_KEY)")
endif

	@sudo make docker-build ARGS=$(DOCKER_BUILD_ARGS) IMG=$(AV_VENDOR) \
		DOCKER_FILE=build/docker/Dockerfile.$(AV_VENDOR) DOCKER_DIR=build/data

multiav-release-av:		## Release an AV inside a docker contrainer.
	$(eval DOCKER_BUILD_ARGS := "")
ifeq ($(AV_VENDOR),sophos)
	$(eval DOCKER_BUILD_ARGS = "--build-arg SOPHOS_URL=$(SOPHOS_URL)")
endif

ifeq ($(AV_VENDOR),symantec)
	$(eval DOCKER_BUILD_ARGS = "--build-arg SYMANTEC_URL=$(SYMANTEC_URL)")
endif

ifeq ($(AV_VENDOR),eset)
	$(eval DOCKER_BUILD_ARGS = "--build-arg ESET_LICENSE_KEY=$(ESET_LICENSE_KEY)")
endif

ifeq ($(AV_VENDOR),drweb)
	$(eval DOCKER_BUILD_ARGS = "--build-arg DR_WEB_LICENSE_KEY=$(DR_WEB_LICENSE_KEY)")
endif

ifeq ($(AV_VENDOR),bitdefender)
	$(eval DOCKER_BUILD_ARGS = "--build-arg BITDEFENDER_URL=$(BITDEFENDER_URL) --build-arg BITDEFENDER_LICENSE_KEY=$(BITDEFENDER_LICENSE_KEY)")
endif

ifeq ($(AV_VENDOR),trendmicro)
	$(eval DOCKER_BUILD_ARGS = "--build-arg TREND_MICRO_LICENSE_KEY=$(TREND_MICRO_LICENSE_KEY)")
endif

	@sudo make docker-release ARGS=$(DOCKER_BUILD_ARGS) IMG=$(AV_VENDOR) \
		VERSION=$(SAFERWALL_VER) DOCKER_FILE=build/docker/Dockerfile.$(AV_VENDOR) DOCKER_DIR=build/data

multiav-build-av-go:	## Build the AV with the gRPC server
	sudo make docker-build IMG=go$(AV_VENDOR) \
	 DOCKER_FILE=build/docker/Dockerfile.go$(AV_VENDOR) \
	 DOCKER_DIR=.

multiav-release-av-go:	## Release the AV with the gRPC server
	sudo make docker-release IMG=go$(AV_VENDOR) VERSION=$(SAFERWALL_VER) \
	 DOCKER_FILE=build/docker/Dockerfile.go$(AV_VENDOR) \
	 DOCKER_DIR=.

AVs = avast avira bitdefender clamav drweb comodo eset fsecure \
		kaspersky mcafee sophos symantec trendmicro windefender
multiav-build: 	## Build all AVs.
	for av in $(AVs) ; do \
		echo "${GREEN} [*] =============== Building $$av =============== ${RESET}" ; \
		make multiav-build-av AV_VENDOR=$$av ; \
		EXIT_CODE=$$? ; \
		if test $$EXIT_CODE ! 0; then \
			make multiav-build-av AV_VENDOR=$$av ; \
		fi; \
	done

multiav-release: ## Build and release all AVs.
	for av in $(AVs) ; do \
		echo "${GREEN} [*] =============== Building $$av =============== ${RESET}" ; \
		make multiav-build-av AV_VENDOR=$$av   ; \
		EXIT_CODE=$$? ; \
		if test $$EXIT_CODE ! 0; then \
			make multiav-release-av AV_VENDOR=$$av ; \
		fi; \
	done

multiav-build-go:	## Build all AVs (go agents).
	for av in $(AVs) ; do \
		echo "${GREEN} [*] =============== Building go-$$av =============== ${RESET}" ; \
		make multiav-build-av-go AV_VENDOR=$$av   ; \
		EXIT_CODE=$$? ; \
		if test $$EXIT_CODE ! 0; then \
			make multiav-build-av-go AV_VENDOR=$$av ; \
		fi; \
	done

multiav-release-go:	## Build and release all AVs (go agents).
	for av in $(AVs) ; do \
		echo "${GREEN} [*] =============== Building go-$$av =============== ${RESET}" ; \
		make multiav-release-av-go AV_VENDOR=$$av  ; \
		EXIT_CODE=$$? ; \
		if test $$EXIT_CODE ! 0; then \
			make multiav-release-av-go AV_VENDOR=$$av ; \
		fi; \
	done
