FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates

# Bypass Go Proxy/SumDB issues for this specific repo context
ENV GOPROXY=direct
ENV GOSUMDB=off

# 1. Copy EVERYTHING first
# This includes source code, go.mod, and potentially broken go.sum
COPY . .

# 2. Cleanup and Fix Dependencies
# Remove existing go.sum to force fresh resolution
RUN rm -f go.sum
# Run tidy to generate correct go.mod and go.sum based on the source code present
RUN go mod tidy
# Download dependencies
RUN go mod download

# 3. Build using the now-correct dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
