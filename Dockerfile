FROM golang:1.10-alpine3.7

LABEL Author=voidint
LABEL Email=voidint@126.com

ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

RUN echo "http://mirrors.aliyun.com/alpine/v3.7/main" > /etc/apk/repositories \
    && echo "http://mirrors.aliyun.com/alpine/v3.7/community" >> /etc/apk/repositories \
    && apk update \
    && apk --no-cache add ca-certificates openssl git \
    && go get -u -v github.com/voidint/gbb \
    && cd $GOPATH/src/github.com/voidint/gbb \
    && gbb --debug \
    && mv ./gbb $GOPATH/bin/gbb \
    && cd $GOPATH \
    && rm -rf $GOPATH/src $GOPATH/pkg

WORKDIR $GOPATH



