package model

import "time"

type Tag struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null;unique"`
	Slug      string    `json:"slug" gorm:"type:varchar(100);not null;unique"`
	Weight    int       `json:"weight" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
