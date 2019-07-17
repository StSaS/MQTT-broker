FROM golang:1.8

COPY . /go/src/server

RUN go get server
RUN go install server

CMD ["server"]
