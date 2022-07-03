FROM golang:1.18.3-alpine3.16

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go mod download
RUN go build -o server cmd/api/*.go

EXPOSE 8000

CMD [ "/app/server" ]
