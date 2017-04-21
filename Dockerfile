FROM golang:alpine

MAINTAINER sine "sinerwr@gmail.com"

RUN apk --update add git
RUN go-wrapper download github.com/SiCo-DevOps/H
RUN apk del git

WORKDIR $GOPATH/src/github.com/SiCo-DevOps/H

RUN go-wrapper install
RUN rm -rf $GOPATH/src

ADD config.sample.json $GOPAH/bin/config.json

EXPOSE 2048

CMD H