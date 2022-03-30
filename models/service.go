package models

import "github.com/jinzhu/gorm"

type Service struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
	Image   string `json:"image"`
}
