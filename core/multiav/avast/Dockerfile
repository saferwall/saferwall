################################
# STEP 1 build executable binary
################################

FROM golang:1.12-alpine AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR $GOPATH/src/saferwall/avast/
COPY . .

# Fetch dependencies.
RUN go get -d -v 

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/avastscanner .


############################
# STEP 2 build a small image
############################

FROM saferwall/avast:0.0.1 AS final
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="gRPC server over linux version of Avast"

# Update VPS
RUN /var/lib/avast/Setup/avast.vpsupdate

# Copy our static executable.
COPY --from=builder /go/bin/avastscanner /bin/avastscanner

# Create an app user so our program doesn't run as root.
RUN groupadd -r saferwall && useradd --no-log-init -r -g saferwall saferwall

# Update permissions
RUN usermod -u 101 saferwall
RUN groupmod -g 102 saferwall

# Switch to our user
USER saferwall

ENTRYPOINT ["/bin/avastscanner"]