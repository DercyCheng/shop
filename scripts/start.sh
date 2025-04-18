#!/bin/bash

# 定义基础路径
SERVICES_ROOT="/Users/dercyc/go/src/Pro/shop"
LOG_DIR="${SERVICES_ROOT}/logs"
SRV_DIR="${SERVICES_ROOT}/shop_srv"
API_DIR="${SERVICES_ROOT}/shop_api"

# 创建日志目录
mkdir -p ${LOG_DIR}

# 显示彩色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    print_info "检查依赖..."
    
    # 检查 Go 环境
    if ! command -v go &> /dev/null; then
        print_error "未找到 Go 环境，请安装 Go 1.16+"
        exit 1
    fi
    
    # 检查必要服务
    if ! command -v mysql &> /dev/null; then
        print_warn "未找到 MySQL 客户端，请确保 MySQL 服务已启动"
    fi
    
    if ! command -v redis-cli &> /dev/null; then
        print_warn "未找到 Redis 客户端，请确保 Redis 服务已启动"
    fi
    
    print_info "依赖检查完成"
}

# 启动服务层服务
start_srv_services() {
    print_info "启动服务层(SRV)服务..."
    
    # 用户服务
    cd ${SRV_DIR}/user_srv
    nohup go run main.go -p 50051 > ${LOG_DIR}/user_srv.log 2>&1 &
    print_info "用户服务(user_srv)已启动，端口: 50051，日志: ${LOG_DIR}/user_srv.log"
    
    # 商品服务
    cd ${SRV_DIR}/goods_srv
    nohup go run main.go -p 50052 > ${LOG_DIR}/goods_srv.log 2>&1 &
    print_info "商品服务(goods_srv)已启动，端口: 50052，日志: ${LOG_DIR}/goods_srv.log"
    
    # 库存服务
    cd ${SRV_DIR}/inventory_srv
    nohup go run main.go -p 50053 > ${LOG_DIR}/inventory_srv.log 2>&1 &
    print_info "库存服务(inventory_srv)已启动，端口: 50053，日志: ${LOG_DIR}/inventory_srv.log"
    
    # 订单服务
    cd ${SRV_DIR}/order_srv
    nohup go run main.go -p 50054 > ${LOG_DIR}/order_srv.log 2>&1 &
    print_info "订单服务(order_srv)已启动，端口: 50054，日志: ${LOG_DIR}/order_srv.log"
    
    # 用户操作服务
    cd ${SRV_DIR}/userop_srv
    nohup go run main.go -p 50055 > ${LOG_DIR}/userop_srv.log 2>&1 &
    print_info "用户操作服务(userop_srv)已启动，端口: 50055，日志: ${LOG_DIR}/userop_srv.log"
    
    # 等待服务启动
    print_info "等待服务层启动完成..."
    sleep 5
}

# 启动API层服务
start_api_services() {
    print_info "启动API层(Web)服务..."
    
    # 用户API
    cd ${API_DIR}/user_web
    nohup go run main.go > ${LOG_DIR}/user_web.log 2>&1 &
    print_info "用户API(user_web)已启动，端口: 8021，日志: ${LOG_DIR}/user_web.log"
    
    # 商品API
    cd ${API_DIR}/goods_web
    nohup go run main.go > ${LOG_DIR}/goods_web.log 2>&1 &
    print_info "商品API(goods_web)已启动，端口: 8023，日志: ${LOG_DIR}/goods_web.log"
    
    # 订单API
    cd ${API_DIR}/order_web
    nohup go run main.go > ${LOG_DIR}/order_web.log 2>&1 &
    print_info "订单API(order_web)已启动，端口: 8024，日志: ${LOG_DIR}/order_web.log"
    
    # OSS服务API
    cd ${API_DIR}/oss_web
    nohup go run main.go > ${LOG_DIR}/oss_web.log 2>&1 &
    print_info "OSS服务API(oss_web)已启动，端口: 8029，日志: ${LOG_DIR}/oss_web.log"
    
    # 用户操作API
    cd ${API_DIR}/userop_web
    nohup go run main.go > ${LOG_DIR}/userop_web.log 2>&1 &
    print_info "用户操作API(userop_web)已启动，端口: 8027，日志: ${LOG_DIR}/userop_web.log"
}

# 停止所有服务
stop_services() {
    print_info "停止所有服务..."
    
    # 查找并停止所有相关进程
    pkill -f "go run main.go"
    
    print_info "所有服务已停止"
}

# 检查服务状态
check_status() {
    print_info "检查服务状态..."
    
    # 检查服务进程
    ps aux | grep -E "go run main.go" | grep -v grep
    
    print_info "服务状态检查完成"
}

# 使用Docker Compose启动系统
start_with_docker() {
    print_info "使用Docker Compose启动系统..."
    
    # 检查Docker是否安装
    if ! command -v docker &> /dev/null; then
        print_error "未找到Docker，请安装Docker和Docker Compose"
        exit 1
    fi
    
    # 启动服务
    docker-compose up -d
    
    print_info "Docker容器已启动，使用 'docker-compose ps' 查看状态"
}

# 使用说明
show_usage() {
    echo "Shop微服务系统管理脚本"
    echo "用法: $0 [命令]"
    echo "命令:"
    echo "  start       启动所有服务（本地模式）"
    echo "  stop        停止所有服务"
    echo "  restart     重启所有服务"
    echo "  status      检查服务状态"
    echo "  start_srv   仅启动服务层(SRV)服务"
    echo "  start_api   仅启动API层(Web)服务"
    echo "  docker      使用Docker Compose启动系统"
    echo "  help        显示此帮助信息"
}

# 主程序
case "$1" in
    start)
        check_dependencies
        start_srv_services
        start_api_services
        print_info "所有服务已启动"
        ;;
    stop)
        stop_services
        ;;
    restart)
        stop_services
        sleep 2
        check_dependencies
        start_srv_services
        start_api_services
        print_info "所有服务已重启"
        ;;
    status)
        check_status
        ;;
    start_srv)
        check_dependencies
        start_srv_services
        ;;
    start_api)
        check_dependencies
        start_api_services
        ;;
    docker)
        start_with_docker
        ;;
    help|*)
        show_usage
        ;;
esac

exit 0