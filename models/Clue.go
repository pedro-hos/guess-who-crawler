package models

import "gorm.io/gorm"

type Clue struct {
	gorm.Model
	ID     uint   `gorm:"primaryKey"`
	Text   string `gorm:"not null"`
	CardId uint
}
