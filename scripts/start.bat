@echo off
chcp 65001 >nul

REM SMTP Mail 启动脚本
REM 用法: start.bat [backend|frontend|all]

set SCRIPT_DIR=%~dp0
set PROJECT_DIR=%~dp0..

REM 默认启动全部
set MODE=%1
if "%MODE%"=="" set MODE=all

if "%MODE%"=="backend" (
    echo 启动后端服务...
    cd /d "%PROJECT_DIR%\backend"
    go run main.go
    goto :eof
)

if "%MODE%"=="frontend" (
    echo 启动前端服务...
    cd /d "%PROJECT_DIR%\frontend"
    npm run dev
    goto :eof
)

if "%MODE%"=="all" (
    echo 启动前后端服务...
    echo.
    echo 服务已启动:
    echo   - 后端: http://localhost:8800
    echo   - 前端: http://localhost:50011
    echo.
    echo 按 Ctrl+C 停止所有服务
    echo.

    REM 启动后端
    cd /d "%PROJECT_DIR%\backend"
    start /B go run main.go > backend.log 2>&1

    REM 启动前端
    cd /d "%PROJECT_DIR%\frontend"
    start /B npm run dev > frontend.log 2>&1

    REM 等待用户中断
    pause >nul

    REM 停止服务
    taskkill /F /IM node.exe >nul 2>&1
    taskkill /F /IM go.exe >nul 2>&1
    echo 服务已停止
    goto :eof
)

echo 用法: start.bat [backend^|frontend^|all]
echo.
echo 参数:
echo   backend  - 只启动后端服务
echo   frontend - 只启动前端服务
echo   all      - 启动全部服务（默认）