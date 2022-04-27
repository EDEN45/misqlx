# syntax=docker/dockerfile:1
FROM golang:1.18

COPY . /go/src/app

WORKDIR /go/src/app/example

RUN go build -o migrate migrate.go

CMD ["./migrate"]