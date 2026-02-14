#!/bin/bash

# Swagger 文档生成脚本

set -e

echo "========================================="
echo "  Swagger 文档生成"
echo "========================================="
echo ""

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# 切换到 backend 目录（脚本在 backend/scripts/ 下）
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"

echo "切换到目录: $BACKEND_DIR"
cd "$BACKEND_DIR"
echo ""

# 检查 swag 是否安装
if ! command -v swag &> /dev/null; then
    echo "错误: swag 未安装"
    echo "请运行: go install github.com/swaggo/swag/cmd/swag@latest"
    exit 1
fi

# 检查 go.mod 是否存在
if [ ! -f "go.mod" ]; then
    echo "错误: 未找到 go.mod 文件"
    echo "当前目录: $(pwd)"
    exit 1
fi

# 生成文档
echo "正在生成 Swagger 文档..."
swag init --parseDependency --parseInternal

# 检查生成结果
if [ -f "docs/swagger.json" ]; then
    FILE_SIZE=$(du -h docs/swagger.json | cut -f1)
    echo ""
    echo "✓ 生成成功！"
    echo ""
    echo "  文件: docs/swagger.json"
    echo "  大小: $FILE_SIZE"
    echo ""
    echo "Swagger UI: http://localhost:2500/swagger/index.html"
    echo ""
else
    echo ""
    echo "✗ 生成失败"
    exit 1
fi
