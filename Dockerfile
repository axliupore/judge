FROM golang:latest AS build

WORKDIR /judge

COPY . .

ENV GOPROXY=https://goproxy.cn,direct

RUN go mod download -x

WORKDIR /judge/cmd
RUN go build -o judge

WORKDIR /judge/go-judge
RUN go mod download -x

RUN go generate ./cmd/go-judge/version \
    && CGO_ENABLE=0 go build -v -tags grpcnotrace,nomsgpack -o go-judge ./cmd/go-judge

FROM debian:latest

RUN sed -i 's@deb.debian.org@mirrors.tuna.tsinghua.edu.cn@g' /etc/apt/sources.list.d/debian.sources

RUN apt update && \
    apt upgrade -y && \
    apt install -y --fix-missing g++ golang-go python3 openjdk-17-jre-headless openjdk-17-jdk && \
    apt clean && rm -rf /var/lib/apt/lists/*

RUN apt update && \
    apt upgrade -y && \
    apt install -y curl gnupg && \
    curl -sL https://deb.nodesource.com/setup_18.x | bash - && \
    apt install -y nodejs && \
    apt clean && rm -rf /var/lib/apt/lists/*

RUN npm config set registry https://registry.npmmirror.com

RUN npm install -g typescript

WORKDIR /judge

COPY --from=build /judge/cmd/judge  /judge/go-judge/go-judge  /judge/go-judge/mount.yaml /judge/run.sh /judge/

EXPOSE 6048

RUN chmod +x /judge/run.sh

ENTRYPOINT ["/judge/run.sh"]