# Needed SHELL since I'm using zsh
SHELL := /bin/bash
.PHONY: help

# Add the following 'help' target to your Makefile
# And add help text after each target name starting with '\#\#'

help: ## This help message
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"


PROTOC_ZIP = protoc-3.3.0-linux-x86_64.zip

include avast.mk
include avira.mk
include bitdefender.mk
include clamav.mk
include comodo.mk
include eset.mk
include fsecure.mk
include kaspersky.mk
include mcafee.mk
include sophos.mk
include symantec.mk
include windefender.mk

install-docker:		## install docker
	sudo apt-get install apt-transport-https ca-certificates curl software-properties-common -y
	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
	sudo apt-key fingerprint 0EBFCD88
	sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $$(lsb_release -cs) stable"
	sudo apt-get update
	sudo apt-get install docker-ce -y

install-protoc:		## install protoc - the Google protobuf compiler 
	curl -OL https://github.com/google/protobuf/releases/download/v3.3.0/$(PROTOC_ZIP)
	sudo unzip -o $(PROTOC_ZIP) -d /usr/local bin/protoc
	rm -f $(PROTOC_ZIP)