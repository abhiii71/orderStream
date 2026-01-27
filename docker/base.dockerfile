FROM golang:1.24-alpine3.20

RUN apk --no-cache add gcc g++ make ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
