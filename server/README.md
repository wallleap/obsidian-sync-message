# OB Sync Server

OB Sync 的后端服务，基于 Go + Gin + SQLite 构建，提供 RESTful API 用于消息同步。

## 环境要求

- Go 1.22 或更高版本
- SQLite3

## 快速开始

### 安装依赖

```bash
go mod download
```

### 运行服务

```bash
go run cmd/main.go
```

服务默认在 `http://localhost:8080` 启动。

### 构建生产版本

```bash
go build -o ob-sync-server cmd/main.go
./ob-sync-server
```

## 配置

服务支持通过环境变量进行配置：

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `SERVER_PORT` | `8080` | 服务监听端口 |

### 数据存储

服务启动时会自动创建以下目录：

- `data/` - 数据库和上传文件存储
  - `obsync.db` - SQLite 数据库
  - `uploads/` - 上传的附件文件
- `logs/` - 日志文件

## API 文档

### 用户管理

#### 生成用户 ID

```http
POST /api/user/generate
```

**响应：**
```json
{
  "user_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

#### 验证用户 ID

```http
POST /api/user/validate
Content-Type: application/json

{
  "user_id": "your-user-id"
}
```

**响应：**
```json
{
  "valid": true
}
```

### 消息管理

#### 发送消息

```http
POST /api/message/send
Content-Type: application/json

{
  "user_id": "your-user-id",
  "type": "text|url",
  "content": "消息内容",
  "original_url": "原始URL（type为url时）"
}
```

**响应：**
```json
{
  "message": "Message sent successfully",
  "id": "message-id"
}
```

#### 上传附件

```http
POST /api/message/upload
Content-Type: multipart/form-data

user_id: your-user-id
file: [文件]
```

**响应：**
```json
{
  "message": "File uploaded successfully",
  "message_id": "message-id"
}
```

#### 同步消息

```http
POST /api/message/sync
Content-Type: application/json

{
  "user_id": "your-user-id",
  "last_sync_time": "2024-01-01T00:00:00Z"
}
```

**响应：**
```json
[
  {
    "id": "message-uuid",
    "type": "text|url|attachment",
    "title": "标题",
    "content": "内容",
    "original_url": "原始URL",
    "created_at": "2024-01-01T12:00:00Z",
    "attachment": {
      "filename": "文件名",
      "file_type": "文件类型"
    }
  }
]
```

#### 下载附件

```http
GET /api/message/file/:id
```

返回附件文件内容。

## URL 抓取功能

服务端支持自动抓取 URL 内容并转换为 Markdown：

### 支持的平台

| 平台 | 抓取方式 | 说明 |
|------|----------|------|
| 微信公众号 | 浏览器渲染 | 自动处理微信验证 |
| 掘金 | HTTP 抓取 | 提取文章正文 |
| 通用网页 | HTTP 抓取 | 使用通用规则提取 |

### 抓取流程

1. 检测 URL 类型
2. 选择合适的抓取策略
3. 提取标题和正文
4. 转换 HTML 为 Markdown
5. 返回处理后的内容

## 项目结构

```
server/
├── cmd/
│   └── main.go              # 入口文件
├── config/
│   └── config.go            # 配置管理
├── internal/
│   ├── handler/
│   │   ├── message_handler.go  # 消息处理器
│   │   └── user_handler.go     # 用户处理器
│   ├── model/
│   │   └── models.go        # 数据模型
│   ├── repository/
│   │   ├── attachment.go    # 附件仓库
│   │   ├── message.go       # 消息仓库
│   │   └── user.go          # 用户仓库
│   └── util/
│       ├── logger.go        # 日志工具
│       ├── snowflake.go     # ID 生成器
│       ├── url_fetcher.go   # URL 抓取
│       └── plugins/         # URL 处理插件
│           ├── interface.go
│           ├── manager.go
│           ├── wechat.go    # 微信公众号
│           ├── juejin.go    # 掘金
│           ├── default.go   # 通用处理
│           ├── converter.go # HTML 转 Markdown
│           └── browser_renderer.go # 浏览器渲染
├── go.mod
└── go.sum
```

## 技术栈

- **Web 框架**：[Gin](https://github.com/gin-gonic/gin)
- **ORM**：[GORM](https://github.com/go-gorm/gorm)
- **数据库**：SQLite3
- **ID 生成**：Snowflake 算法
- **HTML 解析**：goquery
- **Markdown 转换**：自定义转换器

## 开发

### 添加新的 URL 处理插件

1. 在 `internal/util/plugins/` 创建新文件
2. 实现 `Plugin` 接口：

```go
type Plugin interface {
    CanHandle(url string) bool
    ExtractContent(html string) (title string, content string)
}
```

3. 在 `manager.go` 中注册插件

### 日志

日志文件存储在 `logs/` 目录，包含：
- 请求日志
- 错误日志
- URL 抓取日志

## 部署建议

### Docker 部署

```dockerfile
FROM golang:1.22-alpine
WORKDIR /app
COPY . .
RUN go build -o server cmd/main.go
EXPOSE 8080
CMD ["./server"]
```

### 生产环境配置

- 使用反向代理（Nginx/Caddy）
- 配置 HTTPS
- 设置适当的 CORS 策略
- 定期备份数据库
