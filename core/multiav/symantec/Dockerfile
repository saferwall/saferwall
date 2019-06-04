################################
# STEP 1 build executable binary
################################

FROM golang:1.12-alpine AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR $GOPATH/src/saferwall/symantec/
COPY . .

# Fetch dependencies.
RUN go get -d -v 

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/symantecscanner .


############################
# STEP 2 build a small image
############################

FROM saferwall/symantec:0.0.1
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="Symantec Endpoint Protection Linux Client in a docker container"

# Vars
ENV SYMANTEC_SAV	        /opt/Symantec/symantec_antivirus/sav
ENV SYMANTEC_INSTALL_DIR    /opt/Symantec
ENV SYMANTEC_VAR_DIR        /var/symantec

# Update the VDF
RUN /etc/init.d/symcfgd start \
    && /etc/init.d/rtvscand start \ 
    && /etc/init.d/smcd start \
	&& $SYMANTEC_SAV liveupdate --update \
	&& $SYMANTEC_SAV info --defs

# Copy our static executable.
COPY --from=builder /go/bin/symantecscanner /bin/symantecscanner

# Create an app user so our program doesn't run as root.
RUN groupadd -r saferwall && useradd --no-log-init -r -g saferwall saferwall

# Install sudo
RUN apt-get update && apt-get install -y sudo

# Update permissions
RUN usermod -aG sudo saferwall
RUN echo 'saferwall    ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
RUN usermod -u 101 saferwall
RUN groupmod -g 102 saferwall
RUN chown -R saferwall:saferwall $SYMANTEC_INSTALL_DIR 
RUN chown -R saferwall:saferwall $SYMANTEC_VAR_DIR

# Switch to our user
USER saferwall

ENTRYPOINT ["/bin/symantecscanner"]