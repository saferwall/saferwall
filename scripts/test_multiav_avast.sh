#!/bin/bash

if [ -f "/tmp/saferwall/circleci/avast" ]; then
	make go-test GOPKG=github.com/saferwall/saferwall/pkg/multiav/avast
fi
