################################
# STEP 1 build executable binary
################################

FROM golang:1.12-alpine AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR $GOPATH/src/saferwall/clamav/
COPY . .

# Fetch dependencies.
RUN go get -d -v 

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/clamavscanner .


############################
# STEP 2 build a small image
############################

FROM saferwall/clamav:0.0.1
LABEL maintainer="https://github.com/saferwall"
LABEL version="0.0.1"
LABEL description="gRPC server clamav"

# Update VPS
RUN freshclam ; clamscan -V
COPY --from=builder /go/bin/clamavscanner /bin/clamavscanner
ENTRYPOINT ["/bin/clamavscanner"]
