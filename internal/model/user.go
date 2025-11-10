// Package model 定义数据库模型（GORM）。
package model

import "time"

// User 表示用户模型（一个用户可以发表多篇文章）。
// Role 用于 RBAC，默认 "user"；可设置为 "admin" 以访问管理路由。
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username" gorm:"size:64;uniqueIndex;not null"`
    Email     string    `json:"email"    gorm:"size:128;uniqueIndex;not null"`
    Password  string    `json:"-"        gorm:"size:255;not null"`        // 存加密哈希
    Role      string    `json:"role"     gorm:"size:16;not null;default:user"`
    Posts     []Post    `json:"posts,omitempty" gorm:"foreignKey:UserID"` // 一对多关联
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
