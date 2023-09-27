FROM golang:1.20.3-alpine AS builder

COPY . /tools/hbit

RUN apk add build-base libpcap-dev git \
  && go env -w GOPROXY=https://goproxy.cn \
  && go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest \
  && go install -v github.com/projectdiscovery/shuffledns/cmd/shuffledns@latest \
  && go install -v github.com/projectdiscovery/mapcidr/cmd/mapcidr@latest \
  && go install github.com/projectdiscovery/alterx/cmd/alterx@latest \
  && go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest \
  && cd /tools \
  && git clone https://github.com/alwaystest18/cdnChecker.git \
  && cd cdnChecker \
  && go install \
  && go build cdnChecker.go \
  && cd /tools \
  && git clone https://github.com/alwaystest18/hostCollision.git \
  && cd hostCollision \
  && go install \
  && go build hostCollision.go \
  && cd /tools \
  && git clone https://github.com/alwaystest18/dnsVerifier.git \
  && cd dnsVerifier \
  && go install \
  && go build dnsVerifier.go \
  && cd /tools/hbit \
  && go install \
  && go build hbit.go


FROM centos:7

COPY --from=builder /tools /tools
COPY --from=builder /go/bin/* /tools/bin/
COPY ./dict /tools/dict


RUN yum install -y libpcap-devel git gcc make \
  && git clone --branch=master \
               --depth=1 \
               https://github.com/blechschmidt/massdns.git \
  && cd massdns \
  && make \
  && mv bin/massdns /usr/bin/massdns \
  && rm -rf /massdns \
  && cd /opt \
  && curl https://musl.libc.org/releases/musl-1.2.2.tar.gz -o musl-1.2.2.tar.gz \
  && tar -xvf musl-1.2.2.tar.gz \
  && cd musl-1.2.2 \
  && ./configure \
  && make \
  && make install \
  && chmod +x /tools/hbit/naabu_bin \
  && /tools/dnsVerifier/dnsVerifier -r /tools/dnsVerifier/resolvers_all.txt -o /tools/dnsVerifier/resolvers.txt \
  && /tools/dnsVerifier/dnsVerifier -r /tools/cdnChecker/resolvers.txt -o /tools/cdnChecker/resolvers_cn.txt
