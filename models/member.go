package models

import "github.com/jinzhu/gorm"

type Member struct {
	gorm.Model
	Name  string `json:"name"`
	Role  string `json:"role"`
	Image string `json:"image"`
	Quote string `json:"quote"`
}
