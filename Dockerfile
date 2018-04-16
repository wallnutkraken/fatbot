FROM golang:1.10

RUN apt update && apt install git -y

RUN mkdir -p /go/src/github.com/wallnutkraken/fatbot
WORKDIR /go/src/github.com/wallnutkraken/fatbot
COPY . /go/src/github.com/wallnutkraken/fatbot

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV GOBIN $GOPATH/bin

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure --vendor-only

RUN GOMAXPROCS=$(grep -c ^processor /proc/cpuinfo) go install github.com/wallnutkraken/fatbot/cmd/fatbot
CMD ["fatbot"]
