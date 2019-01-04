FROM ubuntu:xenial
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.1"
LABEL description="Avast for Linux in a docker container"

# Install dependencies
RUN apt-get update && apt-get install make -y

# Use the make file itself to install Avira
COPY Makefile ./
RUN make install


# Copy the server binary
COPY server /bin/server

# Expose our gRPC port
EXPOSE 50051

# Entry point
ENTRYPOINT service avast start && /bin/server
