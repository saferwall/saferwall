# Needed SHELL since I'm using zsh
SHELL := /bin/bash

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

# Retrieve the root directory of the project.
ROOT_DIR	:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# Define standard colors
BLACK        := $(shell tput -Txterm setaf 0)
RED          := $(shell tput -Txterm setaf 1)
GREEN        := $(shell tput -Txterm setaf 2)
YELLOW       := $(shell tput -Txterm setaf 3)
LIGHTPURPLE  := $(shell tput -Txterm setaf 4)
PURPLE       := $(shell tput -Txterm setaf 5)
BLUE         := $(shell tput -Txterm setaf 6)
WHITE        := $(shell tput -Txterm setaf 7)

RESET := $(shell tput -Txterm sgr0)

# Our config file.
include $(ROOT_DIR)/.env

# Include our internals makefiles.
include build/mk/aws.mk
include build/mk/docker.mk
include build/mk/minikube.mk
include build/mk/kind.mk
include build/mk/vault.mk
include build/mk/multiav.mk
include build/mk/go.mk
include build/mk/jekyll.mk
include build/mk/nodejs.mk
include build/mk/yarn.mk
include build/mk/nsq.mk
include build/mk/trid.mk
include build/mk/couchbase.mk
include build/mk/k8s.mk
include build/mk/vbox.mk
include build/mk/minio.mk
include build/mk/die.mk
include build/mk/packer.mk
include build/mk/kernel.mk
include build/mk/kvm.mk
include build/mk/helm.mk
include build/mk/kops.mk
include build/mk/consumer.mk
include build/mk/backend.mk
include build/mk/frontend.mk
include build/mk/website.mk
include build/mk/protobuf.mk
include build/mk/yara.mk
include build/mk/certbot.mk
include build/mk/saferwall.mk
include build/mk/elastic.mk
include build/mk/vagrant.mk
