FROM golang:alpine

MAINTAINER sine "sinerwr@gmail.com"

RUN apk --update add git
RUN go-wrapper download github.com/SiCo-DevOps/H
RUN apk del git

WORKDIR $GOPATH/src/github.com/SiCo-DevOps/H

RUN go-wrapper install

WORKDIR $GOPATH/bin/

RUN rm -rf $GOPATH/src

ADD config.sample.json $GOPATH/bin/config.json

EXPOSE 2048

VOLUME $GOPATH/bin/config.json

CMD H