#!/bin/bash

# SMTP Mail 启动脚本
# 用法: ./start.sh [backend|frontend|all]

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# 默认启动全部
MODE=${1:-all}

case "$MODE" in
    backend)
        echo "启动后端服务..."
        cd "$PROJECT_DIR/backend"
        go run main.go
        ;;
    frontend)
        echo "启动前端服务..."
        cd "$PROJECT_DIR/frontend"
        npm run dev
        ;;
    all|*)
        echo "启动前后端服务..."
        # 启动后端（后台）
        cd "$PROJECT_DIR/backend"
        go run main.go &
        BACKEND_PID=$!
        echo "后端 PID: $BACKEND_PID"
        
        # 启动前端
        cd "$PROJECT_DIR/frontend"
        npm run dev &
        FRONTEND_PID=$!
        echo "前端 PID: $FRONTEND_PID"
        
        echo ""
        echo "服务已启动:"
        echo "  - 后端: http://localhost:8800"
        echo "  - 前端: http://localhost:50011"
        echo ""
        echo "按 Ctrl+C 停止所有服务"
        
        # 等待中断信号
        trap "kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit" SIGINT SIGTERM
        wait
        ;;
esac