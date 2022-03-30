package models

import "github.com/jinzhu/gorm"

type GalleryImage struct {
	gorm.Model
	Title string `json:"title"`
	Image string `json:"image"`
}
