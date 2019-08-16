################################
# STEP 1 build executable binary
################################

FROM golang:1.12-alpine AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR $GOPATH/src/saferwall/kaspersky/
COPY . .

# Fetch dependencies.
RUN go get -d -v 

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/kasperskyscanner .


############################
# STEP 2 build a small image
############################

FROM saferwall/kaspersky:0.0.1 AS final
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="Kaspersky Anti-Virus for Linux File Servers in a docker container"

# Vars
ENV KASPERSKY_BIN 			/opt/kaspersky/kesl/bin/kesl-control
ENV KASPERSKY_SETUP 		/opt/kaspersky/kesl/bin/kesl-setup.pl
ENV KASPERSKY_INSTALL_DIR   /opt/kaspersky
# ENV KASPERSKY_VAR_DIR       /var/opt/kaspersky
# ENV KASPERSKY_LOG_DIR       /var/log/kaspersky

# Required packages
RUN apt-get update && apt-get install -y sudo

# Update VDF
RUN service kesl-supervisor start \
    && $KASPERSKY_BIN --start-task 6 \
    && sleep 2m \
    && $KASPERSKY_BIN --app-info

# Create an app user so our program doesn't run as root.
RUN groupadd -r saferwall && useradd --no-log-init -r -g saferwall saferwall

# Copy our binary
COPY --from=builder /go/bin/kasperskyscanner /bin/kasperskyscanner

# Update permissions
RUN usermod -aG sudo saferwall
RUN echo 'saferwall    ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
RUN usermod -u 101 saferwall
RUN groupmod -g 102 saferwall
# RUN chown -R saferwall:saferwall $KASPERSKY_INSTALL_DIR
# RUN chown -R saferwall:saferwall $KASPERSKY_VAR_DIR

# Switch to our user
USER saferwall

ENTRYPOINT ["/bin/kasperskyscanner"]
