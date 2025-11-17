package model

import "time"

type PostTag struct {
	PostId    uint      `json:"post_id" gorm:"primaryKey,index"`
	TagId     uint      `json:"tag_id" gorm:"primaryKey,index"`
	CreatedAt time.Time `json:"created_at"`
}
