#!/bin/bash

if [ -f "/tmp/saferwall/circleci/avast" ] && [ "$1" = "avast" ]; then
	make install-avast
elif [ -f "/tmp/saferwall/circleci/avira" ] && [ "$1" = "avira" ]; then
	make install-avira
fi
