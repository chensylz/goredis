# 使用最新的golang映像
FROM golang:1.19-alpine3.18 as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 设置工作目录
WORKDIR /app

# 复制应用程序代码到容器中
COPY . .

RUN go env -w GO111MODULE=on \
   && go env -w GOPROXY=https://goproxy.cn,direct \
   && go env -w CGO_ENABLED=0 \
   && go env \
   && go mod tidy \
   && go build -o main ./cmd/goredis-server/


FROM alpine:latest

# 设置工作目录
WORKDIR /app

COPY --from=0 /app/. ./

EXPOSE 6379

ENTRYPOINT ./main
