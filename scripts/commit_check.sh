#!/bin/bash

#  Exit immediately if a command exits with a non-zero status.
set -e

# latest commit
LATEST_COMMIT=$(git rev-parse HEAD)

# latest commit where a folder was changed
PKG_AVAST_COMMIT=$(git log -1 --format=format:%H --full-diff pkg/multiav/avast/)
PKG_AVIRA_COMMIT=$(git log -1 --format=format:%H --full-diff pkg/multiav/avira/)
PKG_BITDEFENDER_COMMIT=$(git log -1 --format=format:%H --full-diff pkg/multiav/bitdefender/)
PKG_CLAMAV_COMMIT=$(git log -1 --format=format:%H --full-diff pkg/multiav/clamav/)
PKG_COMODO_COMMIT=$(git log -1 --format=format:%H --full-diff pkg/multiav/comodo/)
PKG_ESET_COMMIT=$(git log -1 --format=format:%H --full-diff pkg/multiav/eset/)
PKG_KASPERSKY_COMMIT=$(git log -1 --format=format:%H --full-diff pkg/multiav/kaspersky/)
PKG_FSECURE_COMMIT=$(git log -1 --format=format:%H --full-diff pkg/multiav/fsecure/)
PKG_MCAFEE_COMMIT=$(git log -1 --format=format:%H --full-diff pkg/multiav/mcafee/)
PKG_SOPHOS_COMMIT=$(git log -1 --format=format:%H --full-diff pkg/multiav/sophos/)
PKG_SYMANTEC_COMMIT=$(git log -1 --format=format:%H --full-diff pkg/multiav/symantec/)

PKG_CRYPTO_COMMIT=$(git log -1 --format=format:%H --full-diff pkg/crypto/)

# create a directory to store the state of changed files
mkdir -p /tmp/saferwall/circleci
touch /tmp/saferwall/circleci/empty

if [ $PKG_AVAST_COMMIT = $LATEST_COMMIT ]; then
	echo "files in pkg/multiav/avast has changed"
	touch /tmp/saferwall/circleci/avast
fi

if [ $PKG_AVIRA_COMMIT = $LATEST_COMMIT ]; then
	echo "files in pkg/multiav/avira has changed"
	touch /tmp/saferwall/circleci/avira
fi

if [ $PKG_BITDEFENDER_COMMIT = $LATEST_COMMIT ]; then
	echo "files in pkg/multiav/bitdefender has changed"
	touch /tmp/saferwall/circleci/bitdefender
fi

if [ $PKG_CLAMAV_COMMIT = $LATEST_COMMIT ]; then
	echo "files in pkg/multiav/clamav has changed"
	touch /tmp/saferwall/circleci/clamav
fi

if [ $PKG_COMODO_COMMIT = $LATEST_COMMIT ]; then
	echo "files in pkg/multiav/comodo has changed"
	touch /tmp/saferwall/circleci/comodo
fi

if [ $PKG_ESET_COMMIT = $LATEST_COMMIT ]; then
	echo "files in pkg/multiav/eset has changed"
	touch /tmp/saferwall/circleci/eset
fi

if [ $PKG_FSECURE_COMMIT = $LATEST_COMMIT ]; then
	echo "files in pkg/multiav/fsecure has changed"
	touch /tmp/saferwall/circleci/fsecure
fi

if [ $PKG_KASPERSKY_COMMIT = $LATEST_COMMIT ]; then
	echo "files in pkg/multiav/kaspersky has changed"
	touch /tmp/saferwall/circleci/kaspersky
fi

if [ $PKG_MCAFEE_COMMIT = $LATEST_COMMIT ]; then
	echo "files in pkg/multiav/mcafee has changed"
	touch /tmp/saferwall/circleci/mcafee
fi

if [ $PKG_SOPHOS_COMMIT = $LATEST_COMMIT ]; then
	echo "files in pkg/multiav/sophos has changed"
	touch /tmp/saferwall/circleci/sophos
fi

if [ $PKG_SYMANTEC_COMMIT = $LATEST_COMMIT ]; then
	echo "files in pkg/multiav/symantec has changed"
	touch /tmp/saferwall/circleci/symantec
fi

if [ $PKG_CRYPTO_COMMIT = $LATEST_COMMIT ]; then
	echo "files in pkg/crypto has changed"
	touch /tmp/saferwall/circleci/crypto
fi
