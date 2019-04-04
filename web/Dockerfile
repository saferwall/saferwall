############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

COPY . $GOPATH/src/saferwall/web/
WORKDIR $GOPATH/src/saferwall/web/

# Fetch dependencies.
# Using go get.
RUN go get -d -v

# Build the binary.
RUN go build -ldflags "-s -w" -o /go/bin/server


############################
# STEP 2 build a small image
############################
FROM alpine:latest

WORKDIR /backend

# Copy our static executable.
COPY --from=builder /go/bin/server server
COPY --from=builder /go/src/saferwall/web/config/app.toml ./config/app.toml
COPY --from=builder /go/src/saferwall/web/app/schema ./app/schema

# Run the hello binary.
ENTRYPOINT ["./server"]