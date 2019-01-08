#!/bin/bash

#  Exit immediately if a command exits with a non-zero status.
set -e

# latest commit
LATEST_COMMIT=$(git rev-parse HEAD)

# latest commit where path/to/folder1 was changed
PKG_AVAST_COMMIT=$(git log -1 --format=format:%H --full-diff scripts/)

# latest commit where path/to/folder2 was changed
PKG_AVIRA_COMMIT=$(git log -1 --format=format:%H --full-diff pkg/multiav/avira)

if [ $PKG_AVAST_COMMIT = $LATEST_COMMIT ];
    then
        echo "files in pkg/multiav/avast has changed"
        touch /tmp/avast.circleci
elif [ $PKG_AVIRA_COMMIT = $LATEST_COMMIT ];
    then
        echo "files in circleci has changed"
        touch /tmp/avira.circleci
else
     echo "no folders of relevance has changed"
     exit 0;
fi
