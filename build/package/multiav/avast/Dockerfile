FROM ubuntu:xenial
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.1"
LABEL description="Avast for Linux in a docker container"

# Add Avast repository to the source list
RUN echo 'deb http://deb.avast.com/lin/repo debian release' >> /etc/apt/sources.list

# Importing Avast PGP public key
RUN apt-key adv --fetch-keys http://files.avast.com/files/resellers/linux/avast.gpg

# Update the package list
RUN apt-get update

# Install avast
RUN apt-get install avast -y

# Setup the license
COPY license.avastlic /etc/avast/license.avastlic
RUN chown avast:avast /etc/avast/license.avastlic

# Patch update script
RUN old='^DOWNLOAD=(.*)$' && new='DOWNLOAD="curl -L -s -f"' \
    && sed -i "s|$old|$new|g" /var/lib/avast/Setup/avast.setup

# Copy the server binary
COPY server /bin/server

# Expose our gRPC port
EXPOSE 50051

# Entry point
ENTRYPOINT service avast start && /bin/server
