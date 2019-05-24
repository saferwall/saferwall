################################
# STEP 1 build executable binary
################################

FROM golang:1.12-alpine AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR $GOPATH/src/saferwall/avira/
COPY . .

# Fetch dependencies.
RUN go get -d -v 

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/avirascanner .


############################
# STEP 2 build a small image
############################

FROM saferwall/avira:0.0.1 AS final
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="gRPC server over linux version of Avira"

# Vars
ENV AVIRA_URL  http://professional.avira-update.com/package/scancl/linux_glibc22/en/scancl-linux_glibc22.tar.gz
ENV AVIRA_FUSEBUNDLE http://install.avira-update.com/package/fusebundlegen/linux_glibc22/en/avira_fusebundlegen-linux_glibc22-en.zip
ENV AVIRA_INSTALL_DIR /opt/avira
ENV AVIRA_TMP /tmp/avira

# Update the VDF
RUN mkdir $AVIRA_TMP \ 
    && wget $AVIRA_FUSEBUNDLE -P $AVIRA_TMP \
    && unzip -o $AVIRA_TMP/avira_fusebundlegen-linux_glibc22-en.zip -d $AVIRA_TMP \
    && $AVIRA_TMP/fusebundle.bin \
    && unzip -o $AVIRA_TMP/install/fusebundle-linux_glibc22-int.zip -d $AVIRA_INSTALL_DIR

# Copy our static executable.
COPY --from=builder /go/bin/avirascanner /bin/avirascanner

# Create an app user so our program doesn't run as root.
RUN groupadd -r saferwall && useradd --no-log-init -r -g saferwall saferwall

# Update permissions
RUN usermod -u 101 saferwall
RUN groupmod -g 102 saferwall
RUN chown -R saferwall:saferwall $AVIRA_INSTALL_DIR

# Switch to our user
USER saferwall

ENTRYPOINT ["/bin/avirascanner"]