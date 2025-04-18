FROM golang:1.16-alpine AS builder

WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache gcc musl-dev

# 复制go.mod和go.sum文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=1 GOOS=linux go build -a -o user_web main.go

# 使用小体积的基础镜像
FROM alpine:latest

WORKDIR /app

# 添加必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 从builder阶段复制编译好的应用
COPY --from=builder /app/user_web .
COPY ./config_*.yaml ./
COPY ./index.html ./

# 暴露API服务端口
EXPOSE 8021

# 设置环境变量
ENV CONFIG_NAME="config_home.yaml"

# 启动应用
CMD ["./user_web"]