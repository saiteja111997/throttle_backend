FROM golang:1.22.5-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

RUN go build -o main ./main.go

RUN chmod +x main

EXPOSE 8080

CMD ["./main"]