# Needed SHELL since I'm using zsh
SHELL := /bin/bash
.PHONY: help


# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

# Retrieve the root directory of the project
ROOT_DIR	:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# Include our env file
include .env

# Make sure variables are exported
export $(shell sed 's/=.*//' .env)

# Include our internals makefiles
include build/docker.mk
include build/vault.mk
include build/multiav.mk
include build/go.mk