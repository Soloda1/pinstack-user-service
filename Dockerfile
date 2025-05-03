FROM golang:1.24.2-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o user-service cmd/server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/user-service .
COPY --from=builder /app/config ./config

EXPOSE 50051

CMD ["./user-service"] 