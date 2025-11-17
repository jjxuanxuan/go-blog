# go-blog

轻量级博客后端，基于 Gin + GORM + MySQL。支持用户注册/登录（JWT）、文章与评论的增删改查。

## 快速开始
- 准备：Go 1.20+、可用的 MySQL 实例
- 设置环境变量（示例）：
  - `DB_HOST=127.0.0.1` `DB_PORT=3306` `DB_USER=app` `DB_PASS=123456` `DB_NAME=go_blog`
  - `JWT_SECRET=change_me`（随机且足够复杂）
  - 可选：`ACCESS_TOKEN_TTL`（默认 120 分钟）、`REFRESH_TOKEN_TTL`（默认 7 天）
- 启动：`go run cmd/server/main.go`
- Base URL：`http://127.0.0.1:8080`

## 更多文档
- 详见 `docs/API.md`（包含完整示例与错误说明）。
