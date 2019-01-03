# Needed SHELL since I'm using zsh
SHELL := /bin/bash
.PHONY: help

# Add the following 'help' target to your Makefile
# And add help text after each target name starting with '\#\#'

help: ## This help message
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"


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