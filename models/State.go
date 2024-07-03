package models

import "gorm.io/gorm"

type State struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"unique;not null"`
	Cities  []City
	Country string `gorm:"unique;not null;default:Brasil"`
}
