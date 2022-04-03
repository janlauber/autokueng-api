package models

import "gorm.io/gorm"

type Link struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Image       string `json:"image"`
}
