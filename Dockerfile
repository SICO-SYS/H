FROM golang:alpine

MAINTAINER sine "sinerwr@gmail.com"

RUN apk --update add git && \
    go-wrapper download github.com/SiCo-Ops/H && \
    apk del git && \
    cd $GOPATH/src && \
    go-wrapper install

EXPOSE 2048

VOLUME $GOPATH/bin/config.json

CMD ["$GOPATH/bin/H"]