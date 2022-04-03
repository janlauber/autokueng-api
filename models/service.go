package models

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
	Image   string `json:"image"`
}
