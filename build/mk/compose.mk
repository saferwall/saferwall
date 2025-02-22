compose/pull: ## Make docker-compose retrieves the latest images of all services.
	docker compose pull

compose/up: ##  Start docker-compose (args: SVC: name of the service to exclude)
	@echo "${GREEN} [*] =============== Docker Compose UP =============== ${RESET}"
	docker compose config --services | grep -v '${SVC}' | xargs docker compose up
