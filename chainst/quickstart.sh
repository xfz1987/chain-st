#!/bin/bash

# ChainST 区块链快速启动脚本

echo "======================================"
echo "  ChainST 区块链快速启动"
echo "======================================"
echo ""

# 检查是否在正确的目录
if [ ! -f "config.yml" ]; then
    echo "错误: 请在 chainst 项目根目录运行此脚本"
    exit 1
fi

# 检查依赖
echo "1. 检查依赖..."
if ! command -v go &> /dev/null; then
    echo "错误: 未找到 Go,请先安装 Go"
    exit 1
fi

if ! command -v ignite &> /dev/null; then
    echo "错误: 未找到 Ignite CLI,请先安装 Ignite"
    exit 1
fi

echo "   ✓ Go $(go version | awk '{print $3}')"
echo "   ✓ Ignite CLI $(ignite version 2>&1 | head -1 | awk '{print $4}')"
echo ""

# 清理旧数据(可选)
read -p "是否清理旧的区块链数据? (y/N): " clean
if [ "$clean" = "y" ] || [ "$clean" = "Y" ]; then
    echo "2. 清理旧数据..."
    ignite chain clean
    echo "   ✓ 数据已清理"
else
    echo "2. 跳过数据清理"
fi
echo ""

# 安装依赖
echo "3. 安装 Go 依赖..."
go mod tidy
echo "   ✓ 依赖安装完成"
echo ""

# 启动区块链
echo "4. 启动区块链..."
echo ""
echo "======================================"
echo "  区块链将在以下端口运行:"
echo "  - RPC:   http://localhost:26657"
echo "  - API:   http://localhost:1317"
echo "  - gRPC:  localhost:9090"
echo ""
echo "  区块链浏览器:"
echo "  - 浏览器: file://$(pwd)/explorer/index.html"
echo ""
echo "  按 Ctrl+C 停止区块链"
echo "======================================"
echo ""

# 启动
ignite chain serve
