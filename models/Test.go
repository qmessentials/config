package models

import "gorm.io/gorm"

type Test struct {
	gorm.Model
	TestName           string    `gorm:"unique" json:"testName"`
	UnitType           string    `json:"unitType"`
	References         *[]string `gorm:"type:text[]" json:"references"`
	Standards          *[]string `gorm:"type:text[]" json:"standards"`
	AvailableModifiers *[]string `gorm:"type:text[]" json:"availableModifiers"`
}
