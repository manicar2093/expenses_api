FROM golang:1.18.3-alpine3.16 as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go mod download
RUN go build -o server cmd/api/*.go

FROM alpine:latest

WORKDIR /api

RUN apk add tzdata

COPY --from=builder /app/server /server

EXPOSE 8000

CMD [ "/server" ]
