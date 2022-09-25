FROM golang:1.19.1-alpine

WORKDIR /app

ENV PORT=3000

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build

EXPOSE 3000

CMD ["./url-shortener"]
