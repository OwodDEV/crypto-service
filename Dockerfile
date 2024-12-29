FROM golang:1.23.1-alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 go build ./cmd/crypto-service

CMD ["/app/crypto-service"]
