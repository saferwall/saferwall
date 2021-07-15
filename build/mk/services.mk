svc-build:		## Build a microservice docker container
	@echo "${GREEN} [*] =============== Build $(SVC) Microservice =============== ${RESET}"
	sudo make docker-build IMG=$(SVC) \
		DOCKER_FILE=build/docker/Dockerfile.$(SVC) DOCKER_DIR=. ;
	@EXIT_CODE=$$?
	@if test $$EXIT_CODE ! 0; then \
		sudo make docker-build IMG=$(SVC) \
			DOCKER_FILE=build/docker/Dockerfile.$(SVC) DOCKER_DIR=. ; \
	fi

svc-release:	## Build and release a microservice docker container
	@echo "${GREEN} [*] =============== Build and Release $(SVC) Microservice =============== ${RESET}"
	sudo make docker-release IMG=$(SVC) \
		VERSION=$(SAFERWALL_VER) \
		DOCKER_FILE=build/docker/Dockerfile.$(SVC) DOCKER_DIR=. ;
	@EXIT_CODE=$$?
	@if test $$EXIT_CODE ! 0; then \
	sudo make docker-release IMG=$(SVC) \
		VERSION=$(SAFERWALL_VER) \
		DOCKER_FILE=build/docker/Dockerfile.$(SVC) DOCKER_DIR=. ;
	fi
