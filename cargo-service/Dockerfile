FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cargo-service/cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/app /usr/local/bin/cargo-service
COPY --from=builder /build/cargo-service/migrations /app/migrations
RUN chmod +x /usr/local/bin/cargo-service
EXPOSE 8001
CMD ["cargo-service"]