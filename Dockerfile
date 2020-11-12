FROM golang:1.15.4-alpine3.12

RUN mkdir /app

ADD . /app

WORKDIR /app

# Install GCC for SqlLite
RUN apk add build-base

RUN go build -o ./build/downloader ./cmd/downloader/main.go
RUN go build -o ./build/server ./cmd/server/main.go

EXPOSE 8080

CMD ["/app/build/server"]