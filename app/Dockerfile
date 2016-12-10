FROM golang:latest

ADD . /go/src/github.com/mariusmilea/piab/app

WORKDIR /go/src/github.com/mariusmilea/piab/app

RUN go get .

RUN go install .

ENTRYPOINT /go/bin/app --alertmanager-server=alertmanager --alertmanager-port=9093

EXPOSE 12345
