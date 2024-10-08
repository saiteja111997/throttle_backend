FROM golang:1.22.5-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy


