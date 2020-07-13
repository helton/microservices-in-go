FROM golang:1.14

WORKDIR /go/src
COPY ./src .
RUN GOOS=linux go build -ldflags="-s -w"
CMD ["./product"]