# OB Sync Frontend

OB Sync 的 Web 前端应用，基于 React + Vite + Tailwind CSS 构建，用于发送消息和上传附件。

## 环境要求

- Node.js 18 或更高版本
- npm 或 pnpm

## 快速开始

### 安装依赖

```bash
npm install
# 或
pnpm install
```

### 开发模式

```bash
npm run dev
```

开发服务器将在 `http://localhost:5173` 启动。

### 构建生产版本

```bash
npm run build
```

构建产物将输出到 `dist/` 目录。

### 预览生产构建

```bash
npm run preview
```

## 功能介绍

### 用户认证

- **生成用户 ID**：首次使用时生成唯一的用户标识
- **验证用户 ID**：已有用户 ID 可直接登录
- **本地存储**：用户 ID 安全存储在浏览器本地
- **导出 ID**：支持下载用户 ID 文件备份

### 发送消息

支持三种消息类型：

#### 1. 文本消息

直接输入文本内容发送，适合快速记录笔记。

#### 2. URL 消息

输入网页 URL，服务端会自动：
- 抓取网页内容
- 提取标题和正文
- 转换为 Markdown 格式

支持的平台：
- 微信公众号文章
- 掘金文章
- 其他通用网页

#### 3. 附件上传

上传任意文件，附件会存储在服务端，同步时下载到 Obsidian 笔记库。

### 消息同步

- 点击同步按钮获取最新消息
- 自动记录上次同步时间
- 增量同步，避免重复

## 项目结构

```
frontend/
├── src/
│   ├── api/
│   │   └── index.ts         # API 客户端
│   ├── components/
│   │   ├── AuthPage.tsx     # 认证页面
│   │   └── MainPage.tsx     # 主页面
│   ├── hooks/
│   │   └── useUserID.ts     # 用户 ID 管理 Hook
│   ├── App.tsx              # 根组件
│   ├── main.tsx             # 入口文件
│   └── index.css            # 全局样式
├── index.html               # HTML 模板
├── vite.config.ts           # Vite 配置
├── tailwind.config.js       # Tailwind 配置
├── postcss.config.js        # PostCSS 配置
├── tsconfig.json            # TypeScript 配置
└── package.json
```

## API 客户端

`src/api/index.ts` 封装了与服务端的所有通信：

```typescript
// 生成用户 ID
generateUserID(): Promise<string>

// 验证用户 ID
validateUserID(userID: string): Promise<boolean>

// 发送消息
sendMessage(userID: string, type: string, content: string, url?: string): Promise<MessageResponse>

// 上传附件
uploadAttachment(userID: string, file: File): Promise<MessageResponse>

// 同步消息
syncMessages(userID: string, lastSyncTime: string): Promise<Message[]>
```

## 组件说明

### AuthPage

认证页面，提供：
- 生成新用户 ID
- 输入已有用户 ID 登录
- 用户 ID 本地存储

### MainPage

主页面，包含：
- 消息列表展示
- 发送消息表单
- URL 输入框
- 文件上传
- 同步按钮
- 处理状态显示

## 样式

使用 Tailwind CSS 进行样式开发：

- 响应式设计
- 深色模式支持（可扩展）
- 自定义颜色主题

## 开发指南

### 修改 API 地址

在 `src/api/index.ts` 中修改 `API_BASE_URL`：

```typescript
const API_BASE_URL = 'http://your-server:8080';
```

### 添加新功能

1. 在 `src/api/index.ts` 添加 API 方法
2. 在 `src/components/` 创建或修改组件
3. 更新样式和交互

### 环境变量

创建 `.env` 文件：

```
VITE_API_URL=http://localhost:8080
```

## 技术栈

- **框架**：React 18
- **构建工具**：Vite 4
- **样式**：Tailwind CSS 3
- **语言**：TypeScript 5
- **本地存储**：localForage

## 部署

### 静态部署

构建后可直接部署到任何静态文件服务器：

```bash
npm run build
# 将 dist/ 目录部署到服务器
```

### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Vercel/Netlify 部署

直接连接 GitHub 仓库，自动构建部署。
