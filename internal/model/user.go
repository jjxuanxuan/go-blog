package model

import "time"

// User 表示用户模型（一个用户可以发表多篇文章）
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"size:64;uniqueIndex;not null"`
	Email     string    `json:"email"    gorm:"size:128;uniqueIndex;not null"`
	Password  string    `json:"-"        gorm:"size:255;not null"`        // 生产环境请存加密哈希
	Posts     []Post    `json:"posts,omitempty" gorm:"foreignKey:UserID"` // 一对多关联
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
