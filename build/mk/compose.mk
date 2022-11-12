LIST_SERVICES = meta-svc pe-svc sandbox-svc aggregator orchestrator postprocessor avira comodo clamav ui webapis

dc-pull: ## Make docker-compose retrieves the lastest images of all services.
	docker-compose pull ${LIST_SERVICES}

dc-up: ##  Start docker-compose
