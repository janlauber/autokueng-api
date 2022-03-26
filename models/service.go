package models

import "github.com/jinzhu/gorm"

type Service struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
	Picture string `json:"picture"`
}
