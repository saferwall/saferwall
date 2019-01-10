#!/bin/bash

export PATH=$PATH:/usr/local/go/bin
export GOPATH=/go/


if [ -f "/tmp/saferwall/circleci/avast" ]  && [ "$1" = "avast" ]; then
	make go-test GOPKG=github.com/saferwall/saferwall/pkg/multiav/avast
elif [ -f "/tmp/saferwall/circleci/avira" ]  && [ "$1" = "avira" ]; then
	make go-test GOPKG=github.com/saferwall/saferwall/pkg/multiav/avira
elif [ -f "/tmp/saferwall/circleci/bitdefender" ] && [ "$1" = "bitdefender" ]; then
	make go-test GOPKG=github.com/saferwall/saferwall/pkg/multiav/bitdefender
elif [ -f "/tmp/saferwall/circleci/clamav" ] && [ "$1" = "clamav" ]; then
	make go-test GOPKG=github.com/saferwall/saferwall/pkg/multiav/clamav
elif [ -f "/tmp/saferwall/circleci/comodo" ] && [ "$1" = "comodo" ]; then
	make go-test GOPKG=github.com/saferwall/saferwall/pkg/multiav/comodo
elif [ -f "/tmp/saferwall/circleci/eset" ] && [ "$1" = "eset" ]; then
	make go-test GOPKG=github.com/saferwall/saferwall/pkg/multiav/eset
fi

