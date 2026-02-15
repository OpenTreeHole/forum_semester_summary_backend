# Forum Semester Summary Backend

论坛学期总结后端服务 - 为每个用户提供个性化的 HTML 页面访问服务

本 readme 由 ai 生成，仅供参考。

## 项目简介

这是一个基于 Go 语言开发的 HTTP 静态文件服务器，用于为论坛用户提供学期总结页面。服务通过用户身份验证（从 HTTP Header 中获取 `X-Consumer-Username`），为每个用户提供对应的个性化 HTML 页面及相关资源文件。

## 核心功能

### 1. 用户身份验证
- 从请求头 `X-Consumer-Username` 获取用户 ID
- 验证用户 ID 的有效性（必须在 1 到 TOTAL_USER 范围内）
- 未授权用户自动重定向到登录页面

### 2. 静态文件服务
- 为每个用户提供独立的资源目录 `resource/{user_id}/`
- 支持访问用户目录下的任意文件
- 默认返回 `{user_id}.html` 作为主页
- 自动处理目录访问和文件访问

### 3. 并行文件下载（可选）
`download.go` 提供了高并发下载功能（当前被注释），可用于批量下载用户数据：
- 支持最多 1000 个并发下载任务
- 从指定 BASE_URL 下载每个用户的 HTML 文件
- 使用 Go 协程和 channel 实现高效并行下载

## 环境变量配置

| 变量名 | 必需 | 说明 | 示例 |
|--------|------|------|------|
| `TOTAL_USER` | 是 | 总用户数量 | `1000` |
| `AUTH_URL` | 是 | 未授权用户重定向的登录页面 URL | `https://example.com/login` |
| `BASE_URL` | 否 | 下载文件时的基础 URL（仅在使用下载功能时需要） | `https://api.example.com?user_id=` |

## 目录结构

```
forum_semester_summary_backend/
├── main.go           # 主服务程序
├── download.go       # 文件下载工具（可选）
├── resource/         # 用户资源目录
│   ├── 1/
│   │   ├── 1.html    # 用户 1 的主页
│   │   └── ...       # 其他资源文件
│   ├── 2/
│   │   ├── 2.html
│   │   └── ...
│   └── ...
└── README.md
```

## 快速开始

### 1. 安装依赖

```bash
go mod init forum_semester_summary_backend
go mod tidy
```

### 2. 配置环境变量

```bash
export TOTAL_USER=1000
export AUTH_URL=https://example.com/login
```

### 3. 准备资源文件

将用户的 HTML 文件和相关资源放入 `resource/{user_id}/` 目录。

### 4. 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## 使用示例

### 访问用户页面

```bash
# 需要在请求头中包含 X-Consumer-Username
curl -H "X-Consumer-Username: 123" http://localhost:8080/
```

### 访问用户资源文件

```bash
curl -H "X-Consumer-Username: 123" http://localhost:8080/assets/style.css
```

## API 路由

### `GET /*`

**功能**: 根据用户 ID 返回对应的静态文件

**请求头**:
- `X-Consumer-Username`: 用户 ID（必需）

**响应**:
- `200 OK`: 成功返回文件内容
- `302 Found`: 未授权，重定向到登录页面
- `400 Bad Request`: 用户 ID 无效或超出范围
- `404 Not Found`: 用户资源目录或文件不存在

## 相关项目

- 数据仓库: https://github.com/OpenTreeHole/DanXi-newyear