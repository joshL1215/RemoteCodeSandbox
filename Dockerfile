# Build
FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/api

# Run
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 5001
CMD ["./app"]
