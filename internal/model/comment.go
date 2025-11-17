package model

import "time"

type Comment struct {
	Id       uint   `json:"id" gorm:"primaryKey"`
	Content  string `json:"content" gorm:"type:longtext;not null"`
	UserId   uint   `json:"user_id" gorm:"index;not null"`
	PostId   uint   `json:"post_id" gorm:"index;not null"`
	ParentId *uint  `json:"parent_id" gorm:"index"`

	// 关联用户和文章
	User    User      `json:"-" gorm:"foreignKey:UserId"`
	Post    Post      `json:"-" gorm:"foreignKey:PostId"`
	Replies []Comment `json:"-" gorm:"foreignKey:ParentId"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
