package model

import "time"

// Post 表示文章模型（每篇文章属于一个用户）
type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"   gorm:"size:200;not null"`
	Content   string    `json:"content" gorm:"type:longtext"`
	UserID    uint      `json:"user_id" gorm:"index;not null"` // 外键
	User      *User     `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
