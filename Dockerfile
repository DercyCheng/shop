FROM golang:1.18-alpine AS builder

WORKDIR /app

# 设置Go环境
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 先复制go.mod和go.sum文件来利用Docker缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 构建应用，ARG允许构建时传入参数
ARG SERVICE_PATH
ARG SERVICE_NAME
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/${SERVICE_NAME} ${SERVICE_PATH}

# 使用更小的镜像
FROM alpine:latest

# 安装必要的SSL证书
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从builder阶段复制编译好的二进制文件
COPY --from=builder /app/bin /app/bin
# 复制配置文件
COPY --from=builder /app/config-pro.yaml /app/config-pro.yaml

# 暴露端口（可在docker-compose或k8s中覆盖）
EXPOSE 8000

# 设置时区
ENV TZ=Asia/Shanghai

# 设置入口点
ARG SERVICE_NAME
ENTRYPOINT ["/app/bin/${SERVICE_NAME}"]