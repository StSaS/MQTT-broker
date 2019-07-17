FROM golang:1.8

COPY . /go/src/client

RUN go get client
RUN go install client

CMD ["client"]
