FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o migrate ./migrations/auto.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/migrate .
COPY --from=builder /app/app .

RUN chmod +x ./migrate ./app
RUN apk add --no-cache bash