#!/bin/bash

if [ -f "/tmp/avast.circleci" ]; then
	make install-avast
fi
