# 第一阶段：构建阶段
FROM golang:1.23.7-alpine AS builder

# 设置工作目录
WORKDIR /app

#ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0

COPY . .

# 接收版本信息作为构建参数
ARG VERSION
ARG BUILD_DATE
ARG COMMIT_ID

# 下载依赖项并构建，注入变量
RUN go mod download
# 用 shell 形式的 RUN 命令，保证环境变量被替换
RUN go build -ldflags="-X main.version=${VERSION} -X main.buildDate=${BUILD_DATE} -X main.commitID=${COMMIT_ID}" -o main main.go

FROM alpine

WORKDIR /app

ARG TIMEZONE=Asia/Shanghai

RUN set -ex \
#    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk add ca-certificates \
    && apk add --no-cache tzdata\
    && rm -rf /var/lib/apk/lists/*

COPY --from=builder /app/main .

VOLUME /data

EXPOSE 8080

CMD ["./main"]


