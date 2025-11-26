package models

import (
	"time"
)

type Payment struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	ReservationID   uint      `json:"reservation_id" gorm:"not null"`
	Amount          float64   `json:"amount" gorm:"not null"`
	Status          string    `json:"status" gorm:"default:pending"` // pending, paid, failed, expired
	PaymentMethod   string    `json:"payment_method"`                // credit_card, gopay, shopeepay, etc.
	MidtransOrderID string    `json:"midtrans_order_id"`             // Order ID from Midtrans
	VaNumber        string    `json:"va_number"`
	VaBank          string    `json:"va_bank"`
	PaymentTime     time.Time `json:"payment_time"`
	CreatedAt       time.Time `json:"created_at"`

	// Relationship
	Reservation Reservation `json:"reservation" gorm:"foreignKey:ReservationID"`
}

type CreatePaymentRequest struct {
	ReservationID uint   `json:"reservation_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"`
}

type PaymentResponse struct {
	ID              uint      `json:"id"`
	ReservationID   uint      `json:"reservation_id"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status"`
	PaymentMethod   string    `json:"payment_method"`
	MidtransOrderID string    `json:"midtrans_order_id"`
	VaNumber        string    `json:"va_number,omitempty"`    // Opsional untuk bank transfer
	VaBank          string    `json:"va_bank,omitempty"`      // Opsional untuk bank transfer
	SnapToken       string    `json:"snap_token,omitempty"`   // Tambahkan ini untuk Snap
	RedirectURL     string    `json:"redirect_url,omitempty"` // Alternatif: redirect URL
	PaymentTime     time.Time `json:"payment_time"`
	CreatedAt       time.Time `json:"created_at"`
}

type MidtransNotification struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
