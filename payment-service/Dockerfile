FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./payment-service/cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/app /usr/local/bin/payment-service
COPY --from=builder /build/payment-service/migrations /app/migrations
RUN chmod +x /usr/local/bin/payment-service
EXPOSE 8007
CMD ["payment-service"]