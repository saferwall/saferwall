#!/bin/bash

# Install vault to pull secrets
make install-vault

# Pull .env file
vault read -field=.env secret/.env | base64 -d > ../.env

if [ -f "/tmp/saferwall/circleci/avast" ] && [ "$1" = "avast" ]; then
	make install-avast
elif [ -f "/tmp/saferwall/circleci/avira" ] && [ "$1" = "avira" ]; then
	make install-avira
elif [ -f "/tmp/saferwall/circleci/bitdefender" ] && [ "$1" = "bitdefender" ]; then
	make install-bitdefender
elif [ -f "/tmp/saferwall/circleci/clamav" ] && [ "$1" = "clamav" ]; then
	make install-clamav
elif [ -f "/tmp/saferwall/circleci/comodo" ] && [ "$1" = "comodo" ]; then
	make install-comodo
elif [ -f "/tmp/saferwall/circleci/eset" ] && [ "$1" = "eset" ]; then
	make install-eset
fi
