FROM golang:1.26-alpine AS builder
WORKDIR /app

# Cache go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN go build -o /app/proxy ./proxy/main.go

# Runtime stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/proxy /app/proxy
EXPOSE 8080
CMD ["/app/proxy"]
