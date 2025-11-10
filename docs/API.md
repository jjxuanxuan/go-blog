# API 文档与示例请求

轻量级博客后端，基于 Gin + GORM + MySQL，支持用户注册/登录（JWT）与文章的增删改查。

## 快速开始
- 配置环境变量：复制 `.env` 并按需修改
- 启动：`go run cmd/server/main.go`（默认 `http://127.0.0.1:8080`）

## 环境变量
- `APP_ENV`：development/production
- `DB_HOST`/`DB_PORT`/`DB_USER`/`DB_PASS`/`DB_NAME`：MySQL 连接
- `JWT_SECRET`：JWT 密钥（必须足够随机）
- `ACCESS_TOKEN_TTL`：令牌有效期（分钟）

示例 DSN：`app:123456@tcp(127.0.0.1:3306)/go_blog?charset=utf8mb4&parseTime=true&loc=Local`

## 认证
- 登录成功后获取 `token`
- 受保护接口需设置请求头：`Authorization: Bearer <token>`
- 中间件：`AuthMiddleware` 校验并解析 JWT，`RequireUser` 确保上下文存在有效用户 ID

401 可能返回：`缺少或非法Token`、`无效Token`、`解析Token失败`、`未登录`、`Token中缺少sub`

## 通用返回规范
- 成功：`{"code":0,"message":"ok|...","data":{...}}`（登录接口除外）
- 400：`{"code":400,"message":"参数错误"}`
- 401：`{"code":401,"message":"未登录|无效Token|..."}`
- 403：`{"code":403,"message":"无权操作该文章"}`
- 404：`{"code":404,"message":"文章不存在|用户不存在"}`
- 500：`{"code":500,"message":"..."}`

## 接口

### 1) 注册 `POST /register`
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
- 示例：
```bash
curl -X POST http://127.0.0.1:8080/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","email":"alice@example.com","password":"secret123"}'
```

### 2) 登录 `POST /login`
- 请求体：
```json
{ "username": "alice", "password": "secret123" }
```
- 成功响应：
```json
{ "message": "登录成功", "token": "<JWT>" }
```
- 示例：
```bash
curl -X POST http://127.0.0.1:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"secret123"}'
```

### 3) 我的信息 `GET /me`（鉴权）
- 示例：
```bash
curl http://127.0.0.1:8080/me \
  -H 'Authorization: Bearer <JWT>'
```
- 成功响应：
```json
{ "message": "ok", "data": {"id":1, "username":"alice", "email":"alice@example.com"} }
```

### 4) 创建文章 `POST /posts`（鉴权）
- 请求体（不需要 user_id）：
```json
{ "title": "Hello", "content": "World" }
```
- 示例：
```bash
curl -X POST http://127.0.0.1:8080/posts \
  -H 'Authorization: Bearer <JWT>' \
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

### 5) 文章列表 `GET /posts`（鉴权）
- 示例：
```bash
curl http://127.0.0.1:8080/posts \
  -H 'Authorization: Bearer <JWT>'
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

### 6) 文章详情 `GET /posts/:id`（鉴权）
- 示例：
```bash
curl http://127.0.0.1:8080/posts/1 \
  -H 'Authorization: Bearer <JWT>'
```
- 成功响应：
```json
{ "code":0, "message":"查询成功", "data": {"id":1,"title":"Hello","user":{"id":1,"username":"alice"}} }
```

### 7) 更新文章 `PUT /posts/:id`（鉴权，作者本人）
- 请求体（任意字段可选）：
```json
{ "title": "New Title", "content": "New Content" }
```
- 示例：
```bash
curl -X PUT http://127.0.0.1:8080/posts/1 \
  -H 'Authorization: Bearer <JWT>' \
  -H 'Content-Type: application/json' \
  -d '{"title":"New Title","content":"New Content"}'
```
- 成功响应：
```json
{ "code":0, "message":"更新成功", "data": {"id":1,"title":"New Title"} }
```

### 8) 删除文章 `DELETE /posts/:id`（鉴权，作者本人）
- 示例：
```bash
curl -X DELETE http://127.0.0.1:8080/posts/1 \
  -H 'Authorization: Bearer <JWT>'
```
- 成功响应：
```json
{ "code":0, "message":"删除成功" }
```

## 其他说明
- 受保护路由统一经过 `AuthMiddleware` 与 `RequireUser`，未携带或非法 Token 将返回 401。
- 首次启动自动迁移数据表（`users`, `posts`）。
- 删除用户会级联删除其文章（`OnDelete:CASCADE`）。如需保留文章，可改为允许 `UserID` 为空并使用 `SET NULL`。
