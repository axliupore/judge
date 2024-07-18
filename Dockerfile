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

FROM debian:latest

RUN echo "deb https://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye main contrib non-free" > /etc/apt/sources.list && \
    echo "deb https://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-updates main contrib non-free" >> /etc/apt/sources.list && \
    echo "deb https://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-backports main contrib non-free" >> /etc/apt/sources.list && \
    echo "deb https://mirrors.tuna.tsinghua.edu.cn/debian-security bullseye-security main contrib non-free" >> /etc/apt/sources.list

RUN apt update && apt upgrade -y && \
    apt install -y --fix-missing g++ golang-go python3 openjdk-17-jre-headless openjdk-17-jdk && \
    apt clean && rm -rf /var/lib/apt/lists/*

RUN apt update &&  apt upgrade -y && \
    apt install -y --fix-missing nodejs npm && \
    apt clean && rm -rf /var/lib/apt/lists/*

RUN npm config set registry https://registry.npmmirror.com

RUN npm install -g typescript

WORKDIR /judge

COPY --from=build /judge/main  /judge/go-judge/go-judge  /judge/go-judge/mount.yaml /judge/config.yaml /judge/

RUN echo '#!/bin/sh\n./main &\n./go-judge' > /judge/start.sh && chmod +x /judge/start.sh


ENTRYPOINT ["/judge/start.sh"]
