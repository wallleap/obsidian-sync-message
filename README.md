# OB Sync

一个跨平台的笔记同步系统，支持将文本、URL 和附件从 Web 端同步到 Obsidian 笔记库。

## 项目简介

OB Sync 由三个核心组件构成：

| 组件 | 技术栈 | 说明 |
|------|--------|------|
| **Server** | Go + Gin + SQLite | 后端服务，提供 REST API |
| **Frontend** | React + Vite + Tailwind | Web 前端，用于发送消息和文件 |
| **Obsidian Plugin** | TypeScript | Obsidian 插件，同步笔记到本地 |

## 功能特性

- **文本同步**：从 Web 端发送文本笔记，自动同步到 Obsidian
- **URL 抓取**：支持微信公众号、掘金等平台的 URL 自动抓取并转换为 Markdown
- **附件上传**：支持上传各类文件附件
- **增量同步**：基于时间戳的增量同步机制
- **自定义模板**：支持自定义笔记标题和 Frontmatter 模板

## 快速开始

### 1. 启动服务端

```bash
cd server
go mod download
go run cmd/main.go
```

服务将在 `http://localhost:8080` 启动。

### 2. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端开发服务器将在 `http://localhost:5173` 启动。

### 3. 安装 Obsidian 插件

```bash
cd obsidian-plugin
npm install
npm run build
```

将生成的 `main.js`、`manifest.json`、`styles.css` 复制到 Obsidian 笔记库的插件目录：

```
<Vault>/.obsidian/plugins/ob-sync/
```

在 Obsidian 设置中启用插件。

## Docker 部署

### 前置条件

- 安装 [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- 确保 Docker Compose v2 已安装

### 启动服务

```bash
# 进入项目目录
cd ob-sync

# 构建并启动所有服务
docker compose up --build

# 或后台运行
docker compose up --build -d
```

### 访问地址

| 服务 | 地址 |
|------|------|
| 前端 | http://localhost:12081 |
| 后端 API | http://localhost:12080 |

### 服务配置说明

- **server**：Go 后端服务，使用 SQLite 数据库
- **frontend**：React 前端，通过 Nginx 反向代理连接后端

### 数据持久化

- 数据库文件存储在 `server_data` 卷
- 日志文件存储在 `server_logs` 卷
- 上传文件存储在 `server_data/uploads`

### 停止服务

```bash
# 停止服务
docker compose down

# 停止服务并删除数据卷（谨慎使用）
docker compose down -v
```

## 使用流程

1. **获取用户 ID**：访问前端页面，点击生成新的用户 ID
2. **配置插件**：在 Obsidian 插件设置中填入用户 ID 和服务器地址
3. **发送消息**：通过 Web 前端发送文本、URL 或上传附件
4. **同步笔记**：在 Obsidian 中点击同步按钮或使用命令同步

## 项目结构

```
ob-sync/
├── server/                    # Go 后端服务
│   ├── cmd/main.go           # 入口文件
│   ├── config/               # 配置管理
│   └── internal/             # 内部模块
│       ├── handler/          # HTTP 处理器
│       ├── model/            # 数据模型
│       ├── repository/       # 数据访问层
│       └── util/             # 工具函数
├── frontend/                  # React 前端
│   └── src/
│       ├── api/              # API 客户端
│       ├── components/       # React 组件
│       └── hooks/            # 自定义 Hooks
└── obsidian-plugin/           # Obsidian 插件
    └── src/
        ├── main.ts           # 插件入口
        └── settings.ts       # 设置管理
```

## 详细文档

- [服务端文档](./server/README.md)
- [前端文档](./frontend/README.md)
- [Obsidian 插件文档](./obsidian-plugin/README.md)

## API 概览

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/user/generate` | POST | 生成新用户 ID |
| `/api/user/validate` | POST | 验证用户 ID |
| `/api/message/send` | POST | 发送文本/URL 消息 |
| `/api/message/upload` | POST | 上传附件 |
| `/api/message/sync` | POST | 同步消息 |
| `/api/message/file/:id` | GET | 下载附件 |

## 环境要求

- **服务端**：Go 1.22+
- **前端**：Node.js 18+
- **插件**：Obsidian 1.4.0+

## 许可证

MIT License
