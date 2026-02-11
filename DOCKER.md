# Docker 开发环境指南

本项目使用 Docker Compose 一键启动前后端开发环境，**无需在本地安装 Go 或 Node.js**。

## 快速开始

### 1. 一键启动开发环境（推荐）

```bash
./start.sh
```

或指定模式：

```bash
./start.sh dev      # 开发模式（热重载）
./start.sh debug    # 调试模式（支持 Go 远程调试）
```

启动后访问：
- **前端**: http://localhost:5777
- **后端 API**: http://localhost:2500
- **健康检查**: http://localhost:2500/health

### 2. start.sh 脚本命令

| 命令 | 说明 |
|------|------|
| `./start.sh` | 启动开发环境（默认，带热重载） |
| `./start.sh dev` | 同默认，启动开发环境 |
| `./start.sh debug` | 启动调试模式（支持 Go Delve 远程调试） |
| `./start.sh build` | 重新构建镜像并启动 |
| `./start.sh stop` | 停止所有服务 |
| `./start.sh restart` | 重启服务 |
| `./start.sh logs` | 查看实时日志 |
| `./start.sh clean` | 清理容器和卷（数据会丢失） |
| `./start.sh help` | 显示帮助信息 |

### 3. 手动使用 Docker Compose

如果不使用脚本，可以直接使用 docker-compose 命令：

```bash
# 启动开发环境
docker-compose up

# 后台运行
docker-compose up -d

# 启动调试模式
docker-compose -f docker-compose.yaml -f docker-compose.override.yaml up

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 重新构建
docker-compose build --no-cache
```

## 服务说明

### 后端服务 (backend)

- **镜像**: 基于 `golang:1.24-alpine` 构建
- **端口**: 2500
- **热重载**: 使用 [Air](https://github.com/air-verse/air) 实现代码变更自动重启
- **工作目录**: `/app` (挂载本地 `./backend` 目录)
- **缓存**: Go mod 缓存卷加速依赖下载

### 前端服务 (frontend)

- **镜像**: `node:20-slim`
- **端口**: 5777
- **包管理器**: pnpm
- **工作目录**: `/app` (挂载本地 `./frontend` 目录)
- **缓存**: node_modules 和 pnpm 缓存卷加速安装

## 调试模式

### Go 远程调试 (Delve)

使用 `./start.sh debug` 启动时，Delve 调试器会在端口 2345 监听。

**VS Code 配置** (`.vscode/launch.json`):

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Connect to Docker",
      "type": "go",
      "request": "attach",
      "mode": "remote",
      "remotePath": "/app",
      "port": 2345,
      "host": "localhost"
    }
  ]
}
```

## 目录挂载

本地代码变更会实时同步到容器中：

- `./backend` → `/app` (后端容器)
- `./frontend` → `/app` (前端容器)

## 注意事项

1. **首次启动较慢**: 首次启动需要下载依赖和构建镜像，可能需要 3-5 分钟，请耐心等待
2. **pnpm 安装**: 前端使用 pnpm 作为包管理器，会自动启用 corepack
3. **网络**: 前后端使用 `thinking-models-network` 桥接网络通信
4. **健康检查**: 后端启动后会进行健康检查，通过后才启动前端

## 故障排除

### 前端启动失败或退出

可能是依赖安装问题，尝试清理重启：

```bash
./start.sh clean
./start.sh
```

### 端口被占用

如果 2500 或 5777 端口被占用，修改 `docker-compose.yaml` 中的端口映射：

```yaml
ports:
  - "2501:2500"  # 将主机端口改为 2501
  - "5778:5777"  # 将主机端口改为 5778
```

### 热重载不工作

检查本地文件系统变更是否同步到容器：

```bash
# 查看容器内文件
docker-compose exec backend ls -la /app
docker-compose exec frontend ls -la /app
```

### 进入容器调试

```bash
# 进入后端容器
docker-compose exec backend sh

# 进入前端容器
docker-compose exec frontend bash

# 查看后端日志
docker-compose logs -f backend

# 查看前端日志
docker-compose logs -f frontend
```

## 生产部署

生产环境建议使用 `backend/Dockerfile`（非 dev 版本）进行构建，该配置使用多阶段构建生成精简镜像。

```bash
# 生产构建
docker build -f backend/Dockerfile -t thinking-models-backend:latest backend/
```

## 技术栈版本

- **Go**: 1.24
- **Node.js**: 20
- **pnpm**: 10.28.2
- **Air**: latest（热重载工具）
- **Delve**: latest（Go 调试器）
