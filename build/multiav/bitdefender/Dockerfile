FROM ubuntu:xenial
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.1"
LABEL description="Bitdefender Scanner for Unices/Unix in a docker container"

# Install dependencies
RUN apt-get update && apt-get install make -y

# Use make inside the container to install Bitdefender
COPY Makefile ./
RUN make install-bitdefender

# Copy the server binary
COPY server /bin/server

# Add the EICAR Anti-Virus Test File
ADD http://www.eicar.org/download/eicar.com.txt /samples/eicar
RUN chmod -R o+r /samples

# Clean up
RUN apt-get purge -y --auto-remove make \
    && rm -rf /tmp/*

# Expose our gRPC port
EXPOSE 50051

# Entry point
ENTRYPOINT ["/bin/server"]
