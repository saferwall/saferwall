################################
# STEP 1 build executable binary
################################

FROM golang:1.12-alpine AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR $GOPATH/src/saferwall/eset/
COPY . .

# Fetch dependencies.
RUN go get -d -v 

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/esetscanner .


############################
# STEP 2 build a small image
############################

FROM saferwall/eset:0.0.1 AS final
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="ESET File Server Security for Linux in a docker container"

# Vars
ENV ESET_CONFIG_DIR 	/etc/opt/eset
ENV ESET_INSTALL_DIR 	/opt/eset

# Update the VDF
RUN /opt/eset/esets/sbin/esets_update

# Create an app user so our program doesn't run as root.
RUN addgroup -S saferwall && adduser -S -G saferwall saferwall  --shell /bin/sh

# Copy our binary
COPY --from=builder /go/bin/esetscanner /bin/esetscanner

# Configure permissions
RUN chown saferwall:saferwall /bin/esetscanner \
    && chmod -R o+rw $ESET_CONFIG_DIR \ 
    && chmod -R o+rw $ESET_INSTALL_DIR

# Switch to our user
USER saferwall

ENTRYPOINT ["/bin/esetscanner"]
