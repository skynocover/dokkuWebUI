FROM golang:1.15.5-alpine3.12

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

WORKDIR /go/src/app/src
RUN go build -o app

CMD ["./app"]