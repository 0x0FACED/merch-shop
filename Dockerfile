FROM golang:1.23.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o avito-shop cmd/app/main.go 

FROM alpine:latest

RUN apk add --no-cache \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    dumb-init \
    bash

WORKDIR /app

COPY --from=builder /app/avito-shop ./
COPY .env ./.env

ENTRYPOINT ["dumb-init", "--"]
CMD ["./avito-shop"]
