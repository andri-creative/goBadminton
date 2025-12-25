package models

import (
	"time"
)

type Reservation struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"user_id" gorm:"not null"`
	CourtID         uint      `json:"court_id" gorm:"not null"`
	ReservationDate time.Time `json:"reservation_date" gorm:"type:date;not null"`
	TimeSlot        string    `json:"time_slot" gorm:"not null"`
	DurationHours   int       `json:"duration_hours" gorm:"default:1"`
	TotalAmount     float64   `json:"total_amount" gorm:"not null"`
	Status          string    `json:"status" gorm:"default:pending"`

	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User  User  `json:"user" gorm:"foreignKey:UserID"`
	Court Court `json:"court" gorm:"foreignKey:CourtID"`
}

type CreateReservationRequest struct {
	CourtID  uint   `json:"court_id" binding:"required"`
	Date     string `json:"date" binding:"required"`
	TimeSlot string `json:"time_slot" binding:"required"`
}

type ReservationResponse struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	CourtID         uint      `json:"court_id"`
	CourtName       string    `json:"court_name"`
	ReservationDate string    `json:"reservation_date"`
	TimeSlot        string    `json:"time_slot"`
	DurationHours   int       `json:"duration_hours"`
	TotalAmount     float64   `json:"total_amount"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type CheckAvailabilityRequest struct {
	Date     string `json:"date" binding:"required"`
	TimeSlot string `json:"time_slot" binding:"required"`
	CourtID  uint   `json:"court_id" binding:"required"`
}
