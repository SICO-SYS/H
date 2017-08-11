FROM golang:alpine

MAINTAINER sine "sinerwr@gmail.com"

RUN apk --update add git && \
    go-wrapper download github.com/SiCo-Ops/H && \
    apk del git && \
    cd $GOPATH/src/github.com/SiCo-Ops/H && \
    go-wrapper install && \
    rm -rf $GOPATH/src

EXPOSE 2048

VOLUME $GOPATH/bin/config.json

CMD ["$GOPATH/bin/H"]