FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates

# Copy go.mod and go.sum (if exists)
COPY go.mod go.sum* ./

# Force tidy to generate/fix go.sum inside the container
# This is critical because local dev environment has SSL issues
RUN go mod tidy

# Download dependencies
RUN go mod download

COPY . .

# Build with -mod=mod to allow go.mod/go.sum updates if needed
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o server ./cmd/server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
