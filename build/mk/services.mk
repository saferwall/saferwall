svc-build:		## Build a microservice docker container
	@echo "${GREEN} [*] =============== Build $(SVC) Microservice =============== ${RESET}"
	$(eval BUILD_ARGS := --build-arg GITHUB_USER=$(GITHUB_USER) --build-arg GITHUB_TOKEN=$(GITHUB_TOKEN))
	make docker-build ARGS="$(BUILD_ARGS)" IMG=$(SVC) \
		DOCKER_FILE=build/docker/Dockerfile.$(SVC) DOCKER_DIR=. ;
	@EXIT_CODE=$$?
	@if test $$EXIT_CODE ! 0; then \
		make docker-build IMG=$(SVC) \
			DOCKER_FILE=build/docker/Dockerfile.$(SVC) DOCKER_DIR=. ; \
	fi

svc-release:	## Build and release a microservice docker container
	@echo "${GREEN} [*] =============== Build and Release $(SVC) Microservice =============== ${RESET}"
	$(eval BUILD_ARGS := --build-arg GITHUB_USER=$(GITHUB_USER) --build-arg GITHUB_TOKEN=$(GITHUB_TOKEN))
	make docker-release ARGS="$(BUILD_ARGS)" IMG=$(SVC) \
		VERSION=$(SAFERWALL_VER) DOCKER_FILE=build/docker/Dockerfile.$(SVC) DOCKER_DIR=. ;
	@EXIT_CODE=$$?
	@if test $$EXIT_CODE ! 0; then \
		make docker-release IMG=$(SVC) VERSION=$(SAFERWALL_VER) \
			DOCKER_FILE=build/docker/Dockerfile.$(SVC) DOCKER_DIR=. ; \
	fi
