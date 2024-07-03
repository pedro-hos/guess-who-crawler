package models

import "gorm.io/gorm"

type Tip struct {
	gorm.Model
	ID     uint   `gorm:"primaryKey"`
	Text   string `gorm:"not null"`
	CardId uint
}
