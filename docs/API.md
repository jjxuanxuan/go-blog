# API 文档与示例请求

轻量级博客后端，基于 Gin + GORM + MySQL，支持用户注册/登录（JWT）、文章与评论的增删改查。

Base URL: `http://127.0.0.1:8080`

- 业务前缀：`/api`
- 认证前缀：`/api/auth`

## 快速开始
- 配置环境变量：复制 `.env` 并按需修改
- 启动：`go run cmd/server/main.go`（默认 `http://127.0.0.1:8080`）

## 环境变量
- `APP_ENV`：development/production
- `DB_HOST`/`DB_PORT`/`DB_USER`/`DB_PASS`/`DB_NAME`：MySQL 连接
- `JWT_SECRET`：JWT 密钥（必须足够随机）
- `ACCESS_TOKEN_TTL`：访问令牌有效期（分钟，默认 120）
- `REFRESH_TOKEN_TTL`：刷新令牌有效期（分钟，默认 10080=7 天）

示例 DSN：`app:123456@tcp(127.0.0.1:3306)/go_blog?charset=utf8mb4&parseTime=true&loc=Local`

## 认证
- 登录成功后返回 `access_token` 与 `refresh_token`
- 受保护接口需设置：`Authorization: Bearer <access_token>`
- 中间件：`AuthMiddleware` 校验并解析 JWT，`RequireUser` 确保上下文存在有效用户 ID

401 可能返回：`缺少或非法Token`、`无效Token`、`解析Token失败`、`未登录`、`Token中缺少sub`、`无效用户ID`

## 通用返回规范
- 成功：`{"code":0,"message":"ok|...","data":{...}}`（`/api/me` 返回无 `code` 字段，见示例）
- 400：`{"code":400,"message":"参数错误"}`（部分接口附带 `detail`）
- 401：`{"code":401,"message":"未登录|无效Token|..."}`
- 403：`{"code":403,"message":"无权操作该文章|无权操作该评论"}`
- 404：`{"code":404,"message":"文章不存在|评论不存在|用户不存在"}`
- 409：`{"code":409,"message":"用户名或邮箱已存在"}`
- 500：`{"code":500,"message":"..."}`（部分接口附带 `detail`）

## 接口

### 1) 注册 `POST /api/auth/register`
- 请求体：
```json
{ "username": "alice", "email": "alice@example.com", "password": "secret123" }
```
- 成功响应：
```json
{
  "code": 0,
  "message": "注册成功",
  "data": {"id": 1, "username": "alice", "email": "alice@example.com"}
}
```
- 可能错误：409（用户名或邮箱已存在）
- 示例：
```bash
curl -X POST http://127.0.0.1:8080/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","email":"alice@example.com","password":"secret123"}'
```

### 2) 登录 `POST /api/auth/login`
- 请求体：
```json
{ "username": "alice", "password": "secret123" }
```
- 成功响应：
```json
{
  "code": 0,
  "message": "登录成功",
  "access_token": "<JWT>",
  "refresh_token": "<JWT>"
}
```
- 示例：
```bash
curl -X POST http://127.0.0.1:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"secret123"}'
```

### 3) 刷新访问令牌 `POST /api/auth/refresh`
- 请求体：
```json
{ "refresh_token": "<JWT>" }
```
- 成功响应：
```json
{ "code": 0, "message": "ok", "access_token": "<JWT>" }
```
- 示例：
```bash
curl -X POST http://127.0.0.1:8080/api/auth/refresh \
  -H 'Content-Type: application/json' \
  -d '{"refresh_token":"<JWT>"}'
```

### 4) 我的信息 `GET /api/me`（鉴权）
- 示例：
```bash
curl http://127.0.0.1:8080/api/me \
  -H 'Authorization: Bearer <ACCESS_JWT>'
```
- 成功响应：
```json
{ "message": "ok", "data": {"id":1, "username":"alice", "email":"alice@example.com"} }
```

### 5) 创建文章 `POST /api/posts`（鉴权）
- 请求体（不需要 user_id）：
```json
{ "title": "Hello", "content": "World" }
```
- 示例：
```bash
curl -X POST http://127.0.0.1:8080/api/posts \
  -H 'Authorization: Bearer <ACCESS_JWT>' \
  -H 'Content-Type: application/json' \
  -d '{"title":"Hello","content":"World"}'
```
- 成功响应：
```json
{
  "code": 0,
  "message": "创建文章成功",
  "data": {"id":1, "title":"Hello", "content":"World", "user_id":1, "created_at":"2024-01-01T00:00:00Z"}
}
```

### 6) 文章列表 `GET /api/posts`（鉴权）
- 示例：
```bash
curl http://127.0.0.1:8080/api/posts \
  -H 'Authorization: Bearer <ACCESS_JWT>'
```
- 成功响应（示例结构，省略字段）：
```json
{
  "code": 0,
  "message": "查询成功",
  "data": [
    {"id":1,"title":"Hello","user_id":1,"user":{"id":1,"username":"alice"}}
  ]
}
```

### 7) 文章详情 `GET /api/posts/:id`（鉴权）
- 示例：
```bash
curl http://127.0.0.1:8080/api/posts/1 \
  -H 'Authorization: Bearer <ACCESS_JWT>'
```
- 成功响应：
```json
{ "code":0, "message":"查询成功", "data": {"id":1,"title":"Hello","user":{"id":1,"username":"alice"}} }
```

### 8) 更新文章 `PUT /api/posts/:id`（鉴权，作者本人）
- 请求体（任意字段可选）：
```json
{ "title": "New Title", "content": "New Content" }
```
- 示例：
```bash
curl -X PUT http://127.0.0.1:8080/api/posts/1 \
  -H 'Authorization: Bearer <ACCESS_JWT>' \
  -H 'Content-Type: application/json' \
  -d '{"title":"New Title","content":"New Content"}'
```
- 成功响应：
```json
{ "code":0, "message":"更新成功", "data": {"id":1,"title":"New Title"} }
```

### 9) 删除文章 `DELETE /api/posts/:id`（鉴权，作者本人）
- 示例：
```bash
curl -X DELETE http://127.0.0.1:8080/api/posts/1 \
  -H 'Authorization: Bearer <ACCESS_JWT>'
```
- 成功响应：
```json
{ "code":0, "message":"删除成功" }
```

### 10) 创建评论 `POST /api/comments`（鉴权）
- 请求体：
```json
{ "post_id": 1, "content": "Nice post!" }
```
- 成功响应：
```json
{
  "code": 0,
  "message": "创建评论成功",
  "data": {"id": 1, "content": "Nice post!", "user_id": 1, "post_id": 1, "created_at": "2024-01-01T00:00:00Z"}
}
```
- 示例：
```bash
curl -X POST http://127.0.0.1:8080/api/comments \
  -H 'Authorization: Bearer <ACCESS_JWT>' \
  -H 'Content-Type: application/json' \
  -d '{"post_id":1,"content":"Nice post!"}'
```

### 11) 删除评论 `DELETE /api/comments/:id`（鉴权，作者本人）
- 示例：
```bash
curl -X DELETE http://127.0.0.1:8080/api/comments/1 \
  -H 'Authorization: Bearer <ACCESS_JWT>'
```
- 成功响应：
```json
{ "code":0, "message":"删除评论成功" }
```

### 12) 某篇文章的评论列表 `GET /api/posts/:post_id/comments`（鉴权）
- 查询参数：`page`（默认 1）、`page_size`（默认 10，最大 100）
- 示例：
```bash
curl 'http://127.0.0.1:8080/api/posts/1/comments?page=1&page_size=10' \
  -H 'Authorization: Bearer <ACCESS_JWT>'
```
- 成功响应（结构）：
```json
{
  "code": 0,
  "message": "查询评论成功",
  "data": {
    "page": 1,
    "page_size": 10,
    "total": 2,
    "list": [
      {
        "id": 1,
        "content": "Nice post!",
        "user": {"id": 2, "username": "bob"},
        "post_id": 1,
        "replies": [
          { "id": 2, "content": "+1", "user": {"id": 3, "username": "tom"}, "post_id": 1, "parent_id": 1 }
        ]
      }
    ]
  }
}
```

## 管理端（预留）
- 前缀：`/api/admin`（需 `admin` 角色，`AuthMiddleware` + `RequireUser` + `RequireRole("admin")`）
- 暂无开放接口，可在此扩展用户管理等能力。

## 其他说明
- 受保护路由统一经过 `AuthMiddleware` 与 `RequireUser`，未携带或非法 Token 将返回 401。
- 首次启动自动迁移数据表（`users`, `posts`, `comments`）。
- 删除用户会级联删除其文章与评论（`OnDelete:CASCADE`）。如需保留，可改为允许外键为空并使用 `SET NULL`。
