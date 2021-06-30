#!/bin/sh

gpg --quiet --batch --yes --decrypt --passphrase="$SECRETS_PASSPHRASE" \
    --output ./.github/secrets/secrets.tar.gz ./.github/secrets/secrets.tar.gz.gpg \
    && tar zxf ./.github/secrets/secrets.tar.gz -C ./build/data
gpg --quiet --batch --yes --decrypt --passphrase="$SECRETS_PASSPHRASE" \
    --output ./.env ./.github/secrets/.env.gpg