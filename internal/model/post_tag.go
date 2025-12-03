package model

import "time"

// PostTag 文章与标签关联表，使用联合主键表示一条绑定关系。
type PostTag struct {
	PostId    uint      `json:"post_id" gorm:"primaryKey,index"`
	TagId     uint      `json:"tag_id" gorm:"primaryKey,index"`
	CreatedAt time.Time `json:"created_at"`
}
