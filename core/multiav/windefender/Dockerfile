################################
# STEP 1 build executable binary
################################

FROM golang:1.12-alpine AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR $GOPATH/src/saferwall/windefender/
COPY . .

# Fetch dependencies.
RUN go get -d -v 

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/windefenderscanner .


############################
# STEP 2 build a small image
############################

FROM saferwall/windefender:0.0.1 AS final
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="gRPC server over Windows Defender port into Linux"

# Vars
ENV WINDOWS_DEFENDER_INSTALL_DIR    /opt/windows-defender
ENV WINDOWS_DEFENDER_UPDATE         https://go.microsoft.com/fwlink/?LinkID=121721&arch=x86

# Update the VDF
# Deactive it for now, as it breaks mpclient
# RUN apt-get install libreadline-dev -y && wget $WINDOWS_DEFENDER_UPDATE -O $WINDOWS_DEFENDER_INSTALL_DIR/engine/mpam-fe.exe \
#     && cd $WINDOWS_DEFENDER_INSTALL_DIR/engine \
#     && cabextract mpam-fe.exe \
#     && rm mpam-fe.exe

# Copy our static executable.
COPY --from=builder /go/bin/windefenderscanner /bin/windefenderscanner

# Create an app user so our program doesn't run as root.
RUN groupadd -r saferwall && useradd --no-log-init -r -g saferwall saferwall

# Update permissions
RUN usermod -u 101 saferwall
RUN groupmod -g 102 saferwall
RUN chown -R saferwall:saferwall $WINDOWS_DEFENDER_INSTALL_DIR

# Switch to our user
USER saferwall

ENTRYPOINT ["/bin/windefenderscanner"]