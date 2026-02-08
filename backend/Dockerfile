FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates

# Bypass Go Proxy/SumDB issues for this specific repo context
ENV GOPROXY=direct
ENV GOSUMDB=off

# Copy go.mod and go.sum (if exists)
COPY go.mod go.sum* ./

# Force tidy to generate/fix go.sum inside the container
RUN go mod tidy

# Download dependencies
RUN go mod download

COPY . .

# Build without -mod=mod (use what we downloaded)
# Force rebuild with correct OS/Arch
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
