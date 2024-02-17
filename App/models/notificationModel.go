package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	Name   string `gorm:"type:varchar(255)"`
	UserID int    `gorm:"default:null"`
}
