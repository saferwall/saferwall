#!/bin/bash

BASE_DIR=$(dirname "$0")

# Install vault to pull secrets
VAULT_ZIP = vault_1.0.1_linux_amd64.zip
wget https://releases.hashicorp.com/vault/1.0.1/vault_1.0.1_linux_amd64.zip
unzip -o $VAULT_ZIP -d /usr/bin
rm -f $VAULT_ZIP

# Pull .env file
vault read -field=.env secret/.env | base64 -d > $BASE_DIR/../.env

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
