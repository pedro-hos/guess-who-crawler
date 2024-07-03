package models

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey"`
	Answer       string `gorm:"not null"`
	WikipediaURL string `gorm:"not null"`
	Tips         []Tip
	CityId       uint
}
