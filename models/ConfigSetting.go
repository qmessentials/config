package models

type ConfigSetting struct {
	Name   string   `gorm:"unique"`
	Values []string `gorm:"type:text[]"`
}
