FROM docker.io/library/golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates

ENV GOPROXY=direct
ENV GOSUMDB=off
ENV GONOSUMDB=*
ENV GOINSECURE=*
ENV GIT_SSL_NO_VERIFY=true

# 1. Copy EVERYTHING first
COPY . .

# 2. Cleanup and Fix Dependencies
# Remove existing go.sum to force fresh resolution
RUN rm -f go.sum
# Run tidy to generate correct go.mod and go.sum based on source
RUN go mod tidy
# Download dependencies
RUN go mod download

# 3. Build
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
