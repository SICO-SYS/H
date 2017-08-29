FROM golang:alpine

MAINTAINER sine "sinerwr@gmail.com"

RUN apk --update add git && \
    go-wrapper download github.com/SiCo-Ops/H && \
    apk del git && \
    cd $GOPATH/src/github.com/SiCo-Ops/H && \
    cp config.sample.json $GOPATH/bin/config.json && \
    go-wrapper install && \
    cd $GOPATH/bin &&\
    rm -rf $GOPATH/src

EXPOSE 2048

CMD ["./H"]