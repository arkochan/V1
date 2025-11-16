# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY .env .
EXPOSE 8080
CMD ["./server"]
