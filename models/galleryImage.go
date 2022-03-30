package models

import "gorm.io/gorm"

type GalleryImage struct {
	gorm.Model
	Title string `json:"title"`
	Image string `json:"image"`
}
