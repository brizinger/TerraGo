FROM golang:1.14

WORKDIR /go/src/app
COPY /bin/main .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["./main"]