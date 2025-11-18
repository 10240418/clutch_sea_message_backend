# Build code
FROM golang:1.24-alpine AS build-stage

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# Run release
FROM alpine:3.14 AS release-stage

# 安装必要的工具（用于健康检查）
RUN apk add --no-cache ca-certificates wget

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=build-stage /app/main /app/main

# 创建日志目录
RUN mkdir -p /app/logs

EXPOSE 8080

ENTRYPOINT ["/app/main"]

# docker build --platform linux/amd64,linux/arm64 -t idreamsky/hisense-vmi-server .