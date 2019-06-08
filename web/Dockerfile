FROM debian:stretch-slim 
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="Saferwall backend"

WORKDIR /backend

# Copy our static executable.
COPY server server
COPY ./config/app.toml ./config/app.toml
COPY ./app/schema ./app/schema

# Run the hello binary.
ENTRYPOINT ["./server"]