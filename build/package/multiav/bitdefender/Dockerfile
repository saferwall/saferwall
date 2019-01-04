FROM ubuntu:xenial
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.1"
LABEL description="Bitdefender Scanner for Unices/Unix in a docker container"

# Install dependencies
RUN apt-get update && apt-get install make -y

# Use make inside the container to install Avast
COPY Makefile ./
RUN make install

# Copy the server binary
COPY server /bin/server

# Add EICAR test file
ADD http://www.eicar.org/download/eicar.com.txt ./

# Expose our gRPC port
EXPOSE 50051

# Entry point
ENTRYPOINT [/bin/server]
