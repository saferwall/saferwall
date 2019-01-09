#!/bin/bash

export PATH=$PATH:/usr/local/go/bin

if [ -f "/tmp/saferwall/circleci/avast" ]  && [ "$1" = "avast" ]; then
	make go-test GOPKG=github.com/saferwall/saferwall/pkg/multiav/avast
elif [ -f "/tmp/saferwall/circleci/avira" ]  && [ "$1" = "avira" ]; then
	make go-test GOPKG=github.com/saferwall/saferwall/pkg/multiav/avira
fi
