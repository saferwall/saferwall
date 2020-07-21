consumer-build:		## build consumer docker container
	@echo "${GREEN} [*] =============== Build Consumer =============== ${RESET}"
	sudo make docker-build IMG=consumer DOCKER_FILE=build/docker/Dockerfile.consumer DOCKER_DIR=. ;
	@EXIT_CODE=$$?
	@if test $$EXIT_CODE ! 0; then \
		sudo make docker-build IMG=consumer VERSION=0.0.2 DOCKER_FILE=build/docker/Dockerfile.consumer DOCKER_DIR=. ; \
	fi

consumer-release:	## release consumer docker container
	@echo "${GREEN} [*] =============== Build and Release Consumer =============== ${RESET}"
	sudo make docker-release IMG=consumer DOCKER_FILE=build/docker/Dockerfile.consumer DOCKER_DIR=. ;
	@EXIT_CODE=$$?
	@if test $$EXIT_CODE ! 0; then \
		sudo make docker-release IMG=consumer VERSION=0.0.2 DOCKER_FILE=build/docker/Dockerfile.consumer DOCKER_DIR=. ; \
	fi
