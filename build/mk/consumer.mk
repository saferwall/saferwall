consumer-build:		## Build consumer docker container.
	@echo "${GREEN} [*] =============== Build Consumer =============== ${RESET}"
	sudo make docker-build IMG=consumer \
		DOCKER_FILE=build/docker/Dockerfile.consumer DOCKER_DIR=./pkg/consumer ;
	@EXIT_CODE=$$?
	@if test $$EXIT_CODE ! 0; then \
		sudo make docker-build IMG=consumer \
			DOCKER_FILE=build/docker/Dockerfile.consumer DOCKER_DIR=./pkg/consumer ; \
	fi

goconsumer-build:		## Build goconsumer docker container.
	@echo "${GREEN} [*] =============== Build goConsumer =============== ${RESET}"
	sudo make docker-build IMG=goconsumer \
		DOCKER_FILE=build/docker/Dockerfile.goconsumer DOCKER_DIR=. ;
	@EXIT_CODE=$$?
	@if test $$EXIT_CODE ! 0; then \
		sudo make docker-build IMG=goconsumer \
			DOCKER_FILE=build/docker/Dockerfile.goconsumer DOCKER_DIR=. ; \
	fi

consumer-release:	## Release consumer docker container.
	@echo "${GREEN} [*] =============== Build and Release Consumer =============== ${RESET}"
	sudo make docker-release IMG=consumer \
		VERSION=$(SAFERWALL_VER) DOCKER_FILE=build/docker/Dockerfile.consumer DOCKER_DIR=./pkg/consumer ;
	@EXIT_CODE=$$?
	@if test $$EXIT_CODE ! 0; then \
		sudo make docker-release IMG=consumer VERSION=$(SAFERWALL_VER) \
			DOCKER_FILE=build/docker/Dockerfile.consumer DOCKER_DIR=./pkg/consumer ; \
	fi

goconsumer-release:	## Release goconsumer docker container.
	@echo "${GREEN} [*] =============== Build and Release goConsumer =============== ${RESET}"
	sudo make docker-release IMG=goconsumer \
		VERSION=$(SAFERWALL_VER) DOCKER_FILE=build/docker/Dockerfile.goconsumer DOCKER_DIR=. ;
	@EXIT_CODE=$$?
	@if test $$EXIT_CODE ! 0; then \
		sudo make docker-release IMG=goconsumer VERSION=$(SAFERWALL_VER) \
			DOCKER_FILE=build/docker/Dockerfile.goconsumer DOCKER_DIR=. ; \
	fi
