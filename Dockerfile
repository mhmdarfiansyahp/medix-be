FROM golang:1.26.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /medix-be ./cmd/api

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /medix-be /app/medix-be
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/config /app/config

EXPOSE 8080

CMD ["./medix-be"]
