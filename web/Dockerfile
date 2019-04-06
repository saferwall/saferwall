FROM ubuntu:xenial

WORKDIR /backend

# Copy our static executable.
COPY server server
COPY ./config/app.toml ./config/app.toml
COPY ./app/schema ./app/schema

# Run the hello binary.
ENTRYPOINT ["./server"]