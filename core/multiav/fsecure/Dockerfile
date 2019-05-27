################################
# STEP 1 build executable binary
################################

FROM golang:1.12-alpine AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR $GOPATH/src/saferwall/fsecure/
COPY . .

# Fetch dependencies.
RUN go get -d -v 

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/fsecurescanner .


############################
# STEP 2 build a small image
############################

FROM saferwall/fsecure:0.0.1 AS final
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="FSecure Linux Security in a docker container"

# Vars
ENV FSECURE_INSTALL_DIR /opt/f-secure
ENV FSECURE_CONFIG_DIR /etc/opt/f-secure/
ENV FSECURE_VERSION     11.10.68
ENV FSECURE_UPDATE      http://download.f-secure.com/latest/fsdbupdate9.run
ENV FSECURE_URL         https://download.f-secure.com/corpro/ls/trial/fsls-$FSECURE_VERSION-rtm.tar.gz
ENV FSECURE_TMP         /tmp/fsecure

# Copy our static executable.
COPY --from=builder /go/bin/fsecurescanner /bin/fsecurescanner

# Update VDF
RUN wget $FSECURE_UPDATE -P $FSECURE_TMP \
	&& mv $FSECURE_TMP/fsdbupdate9.run $FSECURE_INSTALL_DIR \
	&& /etc/init.d/fsaua start \
	&& /etc/init.d/fsupdate start \
	&& $FSECURE_INSTALL_DIR/fsav/bin/dbupdate $FSECURE_INSTALL_DIR/fsdbupdate9.run ; exit 0 \
	&& /opt/f-secure/fsav/bin/fsav --version

# Create an app user so our program doesn't run as root.
RUN groupadd -r saferwall && useradd --no-log-init -r -g saferwall saferwall

# Update permissions
RUN usermod -u 103 messagebus
RUN usermod -u 101 saferwall
RUN groupmod -g 102 saferwall
RUN chown -R saferwall:saferwall $FSECURE_INSTALL_DIR
RUN chown -R saferwall:saferwall $FSECURE_CONFIG_DIR

# Switch to our user
USER saferwall

ENTRYPOINT ["/bin/fsecurescanner"]
