################################
# STEP 1 build executable binary
################################

FROM golang:1.12-alpine AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR $GOPATH/src/saferwall/comodo/
COPY . .

# Fetch dependencies.
RUN go get -d -v 

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/comodoscanner .


############################
# STEP 2 build a small image
############################

FROM saferwall/comodo:0.0.1 AS final
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="gRPC server over linux version of Comodo"

# Vars
ENV COMODO_INSTALL_DIR  /opt/COMODO
ENV COMODO_UPDATE       http://download.comodo.com/av/updates58/sigs/bases/bases.cav

# Update the VDF
ADD $COMODO_UPDATE /opt/COMODO/scanners/bases.cav

# Copy our static executable.
COPY --from=builder /go/bin/comodoscanner /bin/comodoscanner

# Create an app user so our program doesn't run as root.
RUN groupadd -r saferwall && useradd --no-log-init -r -g saferwall saferwall

# Update permissions
RUN usermod -u 101 saferwall
RUN groupmod -g 102 saferwall
RUN chown -R saferwall:saferwall $COMODO_INSTALL_DIR


# Switch to our user
USER saferwall

ENTRYPOINT ["/bin/comodoscanner"]