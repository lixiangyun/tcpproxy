FROM golang:latest
MAINTAINER lixiangyun linimbus@126.com

WORKDIR /gopath/
ENV GOPATH=/gopath/
ENV GOOS=linux
ENV CGO_ENABLED=0

RUN go get -u -v github.com/lixiangyun/tcpproxy
WORKDIR /gopath/src/github.com/lixiangyun/tcpproxy
RUN go build .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /opt/
COPY --from=0 /gopath/src/github.com/lixiangyun/tcpproxy/tcpproxy ./tcpproxy

RUN chmod +x *

ENTRYPOINT ["./tcpproxy"]
