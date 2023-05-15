package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ProductCode string `gorm:"unique" json:"productCode"`
	Description string `json:"description"`
}
