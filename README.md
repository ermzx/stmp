# SMTP Mail

**Git**: https://github.com/ermzx/stmp.git

一个现代化的SMTP客户端Web应用。

## 快速开始

```bash
# 一键启动前后端
./scripts/start.sh        # Linux/Mac
scripts\start.bat         # Windows

# 只启动后端
./scripts/start.sh backend

# 只启动前端
./scripts/start.sh frontend
```

## 环境要求

- Go (版本 1.18 或更高)
- Node.js (版本 16 或更高)
- npm

## 端口配置

- **后端端口**: `config/config.yaml` 或环境变量 `SERVER_PORT`
- **前端端口**: `frontend/.env` 中的 `VITE_PORT`

## 文档

详细说明请参阅 `docs/` 目录下的文档。

## 许可证

BSD 3-Clause License
