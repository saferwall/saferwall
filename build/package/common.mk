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
CURRENT_DIR	:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
ROOT_DIR := $(CURRENT_DIR)/../..

-include docker.mk
-include vault.mk

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

	
api:	## Generates protocol buffers definitions files. 
	protoc -I $(ROOT_DIR)/api/protobuf-spec/ \
		-I${GOPATH}/src \
		--go_out=plugins=grpc:$(ROOT_DIR)/api/protobuf-spec/ \
		$(ROOT_DIR)/api/protobuf-spec/multiav.$(AV_VENDOR).proto

compile: api	## Compile gRPC server
	go build -ldflags "-s -w" -o bin/server server/main.go 

