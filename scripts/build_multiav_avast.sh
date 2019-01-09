#!/bin/bash

if [ -f "/tmp/avast.circleci" ]; then
	sudo --preserve-env make install-avast
fi
