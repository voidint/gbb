FROM centos:centos6

MAINTAINER "voidint <voidint@126.com>"

ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ENV GOLANG_VERSION 1.7.4
ENV GOLANG_DOWNLOAD_URL https://storage.googleapis.com/golang/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 47fda42e46b4c3ec93fa5d4d4cc6a748aa3f9411a2a2b7e08e3a6d80d753ec8b

RUN yum install -y git-core \
    && yum clean all \
    && curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz \
    && mkdir -p /go/{src,bin} \
    && go get -u github.com/constabulary/gb/... \
    && go get -u github.com/voidint/gbb \
    && cd $GOPATH/src/github.com/voidint/gbb \
    && gbb --debug \
    && mv ./gbb $GOPATH/bin/gbb \
    && cd $GOPATH \
    && rm -rf $GOPATH/src


WORKDIR $GOPATH



