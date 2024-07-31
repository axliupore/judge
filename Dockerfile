FROM golang:latest AS build

WORKDIR /judge

COPY . .

ENV GOPROXY=https://proxy.golang.org,direct

RUN go mod download -x

RUN go build -o main

WORKDIR /judge/go-judge
RUN go mod download -x

RUN go generate ./cmd/go-judge/version \
    && CGO_ENABLE=0 go build -v -tags grpcnotrace,nomsgpack -o go-judge ./cmd/go-judge

RUN echo '#!/bin/sh\n./main &\n./go-judge' > /judge/start.sh && chmod +x /judge/start.sh

ENTRYPOINT ["/judge/start.sh"]
