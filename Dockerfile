FROM ubuntu:14.04.4
MAINTAINER showwin <showwin.czy@gmail.com>

# set timezone JST
RUN /bin/cp -p  /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

# install tools
RUN apt-get update
RUN apt-get install -y git wget g++ gcc libc6-dev make pwgen nginx vim-tiny

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

# install MySQL
CMD "debconf-set-selections <<< 'mysql-server mysql-server/root_password password mypassword'"
CMD "debconf-set-selections <<< 'mysql-server mysql-server/root_password_again password mypassword'"
RUN apt-get -y install mysql-server

# Nginx setting for SSL
RUN NGPASS=$(pwgen 16 1) && \
    mkdir -p /usr/local/nginx && \
    mkdir -p /usr/local/nginx/conf && \
    cd /usr/local/nginx/conf && \
    openssl genrsa -passout pass:$NGPASS -des3 -out server.key 1024 && \
    openssl req -new -key server.key -out server.csr -passin pass:$NGPASS -subj "/C=AU/ST=Some-State/O=Internet Widgits Pty Ltd" && \
    cp server.key server.key.org && \
    openssl rsa -in server.key.org -out server.key -passin pass:$NGPASS && \
    openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt

# add Go app
ADD . /go/src/development/showwin/Gizix

# update Nginx SSL setting
RUN cp /go/src/development/showwin/Gizix/docker/nginx_conf /etc/nginx/sites-enabled/gizix

EXPOSE 443

CMD ["/go/src/development/showwin/Gizix/docker/start-server.sh"]
