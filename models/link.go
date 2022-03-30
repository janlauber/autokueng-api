package models

import "github.com/jinzhu/gorm"

type Link struct {
	gorm.Model
	URL   string `json:"url"`
	Image string `json:"image"`
}
