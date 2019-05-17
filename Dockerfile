FROM golang:1.8

WORKDIR /go/src/http2iso
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080

ENTRYPOINT ["http2iso"]
