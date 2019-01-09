#!/bin/bash

if [ -f "/tmp/saferwall/circleci/avast" ]; then
	make install-avast
fi
