ui-docker-run:		## Run the docker container
	sudo docker run -it -p 80:80 --name ui saferwall/ui

ui-build: ## Build frontend in a docker container
	@echo "${GREEN} [*] =============== Build Frontend =============== ${RESET}"
	sudo make docker-build IMG=ui \
		DOCKER_FILE=build/docker/Dockerfile.frontend DOCKER_DIR=ui/ ;
	@EXIT_CODE=$$?
	@if test $$EXIT_CODE ! 0; then \
		sudo make docker-release IMG=ui \
		DOCKER_FILE=build/docker/Dockerfile.frontend DOCKER_DIR=ui/ ; \
	fi

ui-release: ## Build and release frontend in a docker container.
	@echo "${GREEN} [*] =============== Build and Release Frontend =============== ${RESET}"
	sudo make docker-release IMG=ui VERSION=$(SAFERWALL_VER) \
	 DOCKER_FILE=build/docker/Dockerfile.frontend DOCKER_DIR=ui/ ;
	@EXIT_CODE=$$?
	@if test $$EXIT_CODE ! 0; then \
		sudo make docker-release IMG=ui VERSION=$(SAFERWALL_VER) \
		 DOCKER_FILE=build/docker/Dockerfile.frontend DOCKER_DIR=ui/ ; \
	fi