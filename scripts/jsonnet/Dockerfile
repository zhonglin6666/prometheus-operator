FROM golang:1.11-stretch

ENV JSONNET_VERSION 0.10.0

RUN apt-get update -y && apt-get install -y g++ make git jq
RUN cd /tmp && wget https://github.com/google/jsonnet/archive/v${JSONNET_VERSION}.tar.gz && \
    tar xvfz v${JSONNET_VERSION}.tar.gz && \
    cd jsonnet-${JSONNET_VERSION} && \
    make && mv jsonnet /usr/local/bin && \
    rm -rf /tmp/v${JSONNET_VERSION}.tar.gz /tmp/jsonnet-${JSONNET_VERSION}
RUN go get github.com/brancz/gojsontoyaml
RUN go get github.com/campoy/embedmd
RUN go get github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb

RUN mkdir -p /go/src/gitlab.300.cn/paas-k8s/prometheus-operator
WORKDIR /go/src/gitlab.300.cn/paas-k8s/prometheus-operator

RUN chmod -R 777 /go
