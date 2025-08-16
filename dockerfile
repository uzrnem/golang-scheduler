# Dockerfile
# docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v8 -t uzrnem/scheduler:0.0.1 -t uzrnem/scheduler:latest --push . --no-cache
FROM golang:1.23.11-alpine3.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o scheduler-service cmd/main.go

FROM alpine:3.22
WORKDIR /root/
COPY --from=builder /app/scheduler-service .
COPY dashboard ./dashboard

EXPOSE 8080
CMD ["./scheduler-service"]
    