FROM golang:1.19.1-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o urlShortener

# Create multi-stage build to reduce size
FROM alpine:3.16.2

WORKDIR /app

COPY --from=builder /app .

CMD ["./urlShortener"]
