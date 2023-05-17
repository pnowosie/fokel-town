# STEP 1 use a temporary image to build a static binary

FROM golang:1.19 AS builder

# Pull build dependencies
WORKDIR /app
COPY . .
RUN go mod download

# Run tests
RUN go test -v ./...

# Build static image.
RUN GIT_SHA=$(git rev-parse --short HEAD) && \
    CGO_ENABLED=0 GOARCH=amd64 GOOS=linux \
    go build -a \
    -ldflags "-extldflags '-static' -w -s -X main.appSha=$GIT_SHA" \
    -o /go/bin/merkle-service \
    ./cmd/api

# STEP 2 worker image with application binary only

FROM alpine:3.16
COPY --from=builder /go/bin/merkle-service /go/bin/merkle-service

EXPOSE 4000
ENTRYPOINT ["/go/bin/merkle-service", "-host", "0.0.0.0"]