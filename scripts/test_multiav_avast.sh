#!/bin/bash

if [ -f "/tmp/avast.circleci" ]; then
	sudo --preserve-env make go-test GOPKG=github.com/saferwall/saferwall/pkg/multiav/avast
fi
