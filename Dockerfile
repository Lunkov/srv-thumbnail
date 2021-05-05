# Two-stage build:
#    first  FROM prepares a binary file in full environment ~780MB
#    second FROM takes only binary file ~20MB
FROM golang:1.15 AS builder

# File Author / Maintainer
MAINTAINER DIGITAL COUNTRY

RUN apt-get update && apt-get install -y net-tools dnsutils  ca-certificates libproj-dev protobuf-compiler && apt-get clean -y
WORKDIR /root

ENV GO111MODULE=auto
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOROOT=/usr/local/go
ENV GOBIN=/root/go
ENV GOPATH $HOME/go
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin:$GOROOT:$GOPATH:$GOBIN
ENV CGO_CFLAGS="-g -O2"
ENV CGO_CPPFLAGS=""
ENV CGO_CXXFLAGS="-g -O2"
ENV CGO_FFLAGS="-g -O2"
ENV CGO_LDFLAGS="-g -O2"
#ENV GCCGO="gccgo"
#ENV CC="clang"
#ENV CXX="clang++"
ENV GOGCCFLAGS="-fPIC -m64 -pthread -fno-caret-diagnostics -Qunused-arguments -fmessage-length=0 -gno-record-gcc-switches -fno-common"

RUN cat /etc/*-release
RUN pwd
RUN ls -l /home/
RUN ls -l /root/

RUN mkdir /app && mkdir /app/templates && mkdir /app/etc
ADD *.go /app/
# ADD go.* /app/
WORKDIR /app/
RUN cd /app

RUN ls -l
RUN set
RUN go version
RUN go get -d
RUN go build -v -o ./web-service ./...

RUN ls -l /go/src/

RUN rm -rf /root/.ssh/

#########
# second stage to obtain a very small image
FROM alpine:latest
# File Author / Maintainer
MAINTAINER DIGITAL COUNTRY

RUN mkdir /app && mkdir /app/static && mkdir /app/etc

VOLUME /app/etc
VOLUME /app/templates
VOLUME /app/storage

WORKDIR /app

COPY --from=builder /app/web-service /app/web-service
RUN chmod +x /app/web-service

# libproj-dev 
RUN apk update && \
    apk add -u ca-certificates && \
    rm -rf /var/lib/apt/lists/*

ADD ./docker/nsswitch.conf /etc/nsswitch.conf
#ADD ./web-service /app/web-service

# Run the command on container startup
EXPOSE 3000
CMD ["/app/web-service"]
