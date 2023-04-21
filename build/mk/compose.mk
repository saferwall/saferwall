dc-pull: ## Make docker-compose retrieves the lastest images of all services.
	docker-compose pull

dc-up: ##  Start docker-compose (args: SVC: name of the service to exclude)
	@echo "${GREEN} [*] =============== Docker Compose UP =============== ${RESET}"
	docker compose config --services | grep -v '${SVC}\|ml' | xargs docker compose up
