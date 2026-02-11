#!/bin/bash

# 思维模型平台 - Docker Compose 启动脚本

set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}  思维模型平台 - Docker 开发环境启动${NC}"
echo -e "${BLUE}============================================${NC}"
echo ""

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: Docker 未安装${NC}"
    echo "请先安装 Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

# 检查 Docker Compose 是否安装
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo -e "${RED}错误: Docker Compose 未安装${NC}"
    echo "请先安装 Docker Compose: https://docs.docker.com/compose/install/"
    exit 1
fi

# 确定使用哪个命令
if docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose"
else
    COMPOSE_CMD="docker-compose"
fi

# 显示使用说明
function show_usage() {
    echo "用法: ./start.sh [选项]"
    echo ""
    echo "选项:"
    echo "  dev       - 启动开发环境（默认，带热重载）"
    echo "  debug     - 启动调试模式（支持远程调试）"
    echo "  build     - 重新构建镜像并启动"
    echo "  stop      - 停止所有服务"
    echo "  restart   - 重启服务"
    echo "  logs      - 查看日志"
    echo "  clean     - 清理容器和卷"
    echo "  help      - 显示此帮助信息"
    echo ""
}

# 启动开发环境
function start_dev() {
    echo -e "${YELLOW}正在启动开发环境...${NC}"
    echo -e "${GREEN}后端: http://localhost:2500${NC}"
    echo -e "${GREEN}前端: http://localhost:5777${NC}"
    echo ""
    $COMPOSE_CMD -f docker-compose.yaml up
}

# 启动调试模式
function start_debug() {
    echo -e "${YELLOW}正在启动调试模式...${NC}"
    echo -e "${GREEN}后端: http://localhost:2500${NC}"
    echo -e "${GREEN}前端: http://localhost:5777${NC}"
    echo -e "${GREEN}Delve 调试端口: localhost:2345${NC}"
    echo ""
    $COMPOSE_CMD -f docker-compose.yaml -f docker-compose.override.yaml up
}

# 重新构建并启动
function build_and_start() {
    echo -e "${YELLOW}正在重新构建镜像...${NC}"
    $COMPOSE_CMD -f docker-compose.yaml build --no-cache
    echo -e "${YELLOW}正在启动服务...${NC}"
    $COMPOSE_CMD -f docker-compose.yaml up
}

# 停止服务
function stop_services() {
    echo -e "${YELLOW}正在停止服务...${NC}"
    $COMPOSE_CMD -f docker-compose.yaml down
    echo -e "${GREEN}服务已停止${NC}"
}

# 重启服务
function restart_services() {
    echo -e "${YELLOW}正在重启服务...${NC}"
    $COMPOSE_CMD -f docker-compose.yaml restart
    echo -e "${GREEN}服务已重启${NC}"
}

# 查看日志
function show_logs() {
    $COMPOSE_CMD -f docker-compose.yaml logs -f
}

# 清理容器和卷
function clean_all() {
    echo -e "${RED}警告: 这将删除所有容器、卷和缓存数据${NC}"
    read -p "确定要继续吗? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        $COMPOSE_CMD -f docker-compose.yaml down -v --remove-orphans
        echo -e "${GREEN}清理完成${NC}"
    else
        echo -e "${YELLOW}已取消${NC}"
    fi
}

# 主逻辑
case "${1:-dev}" in
    dev)
        start_dev
        ;;
    debug)
        start_debug
        ;;
    build)
        build_and_start
        ;;
    stop)
        stop_services
        ;;
    restart)
        restart_services
        ;;
    logs)
        show_logs
        ;;
    clean)
        clean_all
        ;;
    help|--help|-h)
        show_usage
        ;;
    *)
        echo -e "${RED}未知选项: $1${NC}"
        show_usage
        exit 1
        ;;
esac
