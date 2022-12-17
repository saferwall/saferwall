dc-pull: ## Make docker-compose retrieves the lastest images of all services.
	docker-compose pull

dc-up: ##  Start docker-compose (args: SVC, name of the service to exclude)
	@echo "${GREEN} [*] =============== Docker Compose UP =============== ${RESET}"
	docker compose ps --services | grep -v '${SVC}\|ml\|comodo\|clamav' | xargs docker-compose up
