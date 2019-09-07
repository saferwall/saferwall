################################
# STEP 1 build executable binary
################################

FROM golang:1.13-alpine AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR $GOPATH/src/saferwall/web/
COPY . .

# Fetch dependencies.
RUN go get ./...

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/server .


############################
# STEP 2 build a small image
############################

FROM debian:stretch-slim as final
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="Saferwall backend"

WORKDIR /backend

# Copy the app
COPY --from=builder /go/bin/server .
COPY ./config/app.toml ./config/app.toml
COPY ./app/schema ./app/schema

# Run the server.
ENTRYPOINT ["/backend/server"]