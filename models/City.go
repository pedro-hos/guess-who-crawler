package models

import "gorm.io/gorm"

type City struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"not null"`
	StateId uint
	Cards   []Card
}
