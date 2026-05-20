FROM golang:1.26-alpine AS builder
WORKDIR /app

# Cache go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download && go mod tidy

# Copy source and build
COPY . .
# Build static binary without cgo
ENV CGO_ENABLED=0
RUN go build -o /app/proxy . && chmod +x /app/proxy

# Runtime stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/proxy /app/proxy
EXPOSE 8080
CMD ["/app/proxy"]
