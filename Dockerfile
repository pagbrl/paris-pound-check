## BUILDER PART
FROM golang:alpine AS builder

COPY . $GOPATH/src/paris-pound-check
WORKDIR $GOPATH/src/paris-pound-check
RUN ls -la $GOPATH/src/paris-pound-check

RUN apk update && apk add --no-cache git

RUN adduser -D -g '' poundcheck

RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w" -o /go/bin/paris-pound-check

## RUNNER PART
FROM scratch

# We copy the user entry from the builder
COPY --from=builder /etc/passwd /etc/passwd

# We also need the ca-certificates for x509
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# We copy the binary from the builder
COPY --from=builder /go/bin/paris-pound-check /go/bin/paris-pound-check

USER poundcheck

# Run the binary.
ENTRYPOINT ["/go/bin/paris-pound-check"]
