FROM ubuntu:xenial
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.1"
LABEL description="Sophos Anti-Virus for Linux in a docker container"

# Install dependencies
RUN apt-get update && apt-get install make -y

# Use make inside the container to install Kaspersky
COPY Makefile ./
RUN make install

# Copy the server binary
COPY server /bin/server

# Add the EICAR Anti-Virus Test File
ADD http://www.eicar.org/download/eicar.com.txt ./

# Expose our gRPC port
EXPOSE 50051

# Entry point
ENTRYPOINT ["/bin/server"]