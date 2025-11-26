package models

import (
	"time"
)

type Court struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"not null"`
	Status       string    `json:"status" gorm:"default:active"`
	Location     string    `json:"location" gorm:"not null"`
	PricePerHour float64   `json:"price_per_hour" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
}

type CourtResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Status       string  `json:"status"`
	Location     string  `json:"location"`
	PricePerHour float64 `json:"price_per_hour"`
}
