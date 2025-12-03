package model

import "time"

// Category 表示分类模型，支持父子层级。
type Category struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null;unique"`
	ParentId  *uint     `json:"parent_id,omitempty" gorm:"index"`
	Slug      string    `json:"slug" gorm:"type:varchar(100);not null;unique"`
	Sort      uint      `json:"sort" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
