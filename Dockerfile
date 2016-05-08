FROM ubuntu:14.04
MAINTAINER showwin <showwin.czy@gmail.com>

# set timezone JST
RUN /bin/cp -p  /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

# install tools
RUN apt-get update
RUN apt-get install -y git wget g++ gcc libc6-dev make vim-tiny

# install Go
RUN cd /tmp && \
    wget https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz && \
    tar zxvf go1.6.linux-amd64.tar.gz && \
    rm go1.6.linux-amd64.tar.gz && \
    mv go /usr/local/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/development/showwin/Gizix

# add Go app
ADD . /go/src/development/showwin/Gizix

EXPOSE 8080

CMD ["/go/src/development/showwin/Gizix/docker/start-server.sh"]
