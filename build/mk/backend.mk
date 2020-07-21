backend-build: ## Build backend in a docker container
	@echo "${GREEN} [*] =============== Build Backend =============== ${RESET}"
	sudo make docker-build IMG=backend DOCKER_FILE=build/docker/Dockerfile.backend DOCKER_DIR=. ;
	@EXIT_CODE=$$?
	@if test $$EXIT_CODE ! 0; then \
		sudo make docker-release IMG=backend DOCKER_FILE=build/docker/Dockerfile.backend DOCKER_DIR=. ; \
	fi


backend-release: ## Release backend in a docker container
	@echo "${GREEN} [*] =============== Build and Release Backend =============== ${RESET}"
	sudo make docker-release IMG=backend VERSION=0.0.2 DOCKER_FILE=build/docker/Dockerfile.backend DOCKER_DIR=. ;
	@EXIT_CODE=$$?
	@if test $$EXIT_CODE ! 0; then \
		sudo make docker-release IMG=backend VERSION=0.0.2 DOCKER_FILE=build/docker/Dockerfile.backend DOCKER_DIR=. ; \
	fi