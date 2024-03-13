# 使用官方 Golang 镜像作为构建环境
FROM golang:1.21-alpine as builder
LABEL authors="nite07"

# 设置工作目录
WORKDIR /app

# 复制源代码到工作目录
COPY . .
RUN go mod download

# 获取参数
ARG version

# 使用 -ldflags 参数进行编译
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X sub2clash/config.Version=${version}" -o sub2clash main.go

FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从 builder 镜像中复制出编译好的二进制文件
COPY --from=builder /app/sub2clash /app/sub2clash

# 设置容器的默认启动命令
ENTRYPOINT ["/app/sub2clash"]
