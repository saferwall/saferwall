FROM ubuntu:bionic
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="Avast for Linux in a docker container"

# Requried for apt-key
RUN apt-get update && apt-get install -y --no-install-recommends gnupg2

# Install Avast
RUN echo 'deb http://deb.avast.com/lin/repo ubuntu release' | tee --append /etc/apt/sources.list \
    && apt-key adv --fetch-keys http://files.avast.com/files/resellers/linux/avast.gpg \
    && apt-get update \
    && apt-get install -y --no-install-recommends avast \
    && rm -rf /var/lib/apt/lists/*

# Patch update script
RUN old='^DOWNLOAD=(.*)$' && new='DOWNLOAD="curl -L -s -f"' \
    && sed -i "s|$old|$new|g" /var/lib/avast/Setup/avast.setup \
    && touch /etc/avast/whitelist

# Setup the license
COPY license.avastlic /etc/avast/license.avastlic
RUN chown avast:avast /etc/avast/license.avastlic

# Add EICAR Anti-Virus Test File
ADD --chown=avast:avast http://www.eicar.org/download/eicar.com.txt eicar

#  Performs a simple test
RUN service avast start && scan eicar; rm eicar

