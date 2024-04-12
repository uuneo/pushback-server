# 第一阶段：构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

COPY . .

# 下载依赖项
RUN go mod download \
    && go build -o /app/main main.go


FROM alpine:latest

WORKDIR /app

ARG TIMEZONE=Asia/Shanghai

RUN set -ex \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk add ca-certificates \
    && apk add --no-cache tzdata\
    && rm -rf /var/lib/apk/lists/*

COPY --from=builder /app/main /app/main

VOLUME /data

EXPOSE 8080

CMD ["./main"]


