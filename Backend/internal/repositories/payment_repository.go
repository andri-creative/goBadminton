package repositories

import (
	"backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	UpdatePaymentStatus(ctx context.Context, orderID string, status string) error
	GetPaymentByOrderID(ctx context.Context, orderID string) (*models.Payment, error)

	GetUserPayments(ctx context.Context, userID uint) ([]models.Payment, error)
	GetPaymentByID(ctx context.Context, paymentID uint, userID uint) (*models.Payment, error)
	Update(ctx context.Context, payment *models.Payment) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// -----------------------------------------------------
// CREATE PAYMENT
// -----------------------------------------------------
func (r *paymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// -----------------------------------------------------
// UPDATE PAYMENT STATUS
// -----------------------------------------------------
func (r *paymentRepository) UpdatePaymentStatus(ctx context.Context, orderID string, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.Payment{}).
		Where("midtrans_order_id = ?", orderID).
		Update("status", status).Error
}

// -----------------------------------------------------
// GET PAYMENT BY MIDTRANS ORDER ID
// -----------------------------------------------------
func (r *paymentRepository) GetPaymentByOrderID(ctx context.Context, orderID string) (*models.Payment, error) {
	var payment models.Payment

	err := r.db.WithContext(ctx).
		Where("midtrans_order_id = ?", orderID).
		First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

// -----------------------------------------------------
// GET ALL USER PAYMENTS
// -----------------------------------------------------
func (r *paymentRepository) GetUserPayments(ctx context.Context, userID uint) ([]models.Payment, error) {
	var payments []models.Payment

	err := r.db.WithContext(ctx).
		Joins("JOIN reservations ON reservations.id = payments.reservation_id").
		Where("reservations.user_id = ?", userID).
		Preload("Reservation").
		Find(&payments).Error

	return payments, err
}

// -----------------------------------------------------
// GET PAYMENT BY ID (AND BELONGS TO USER)
// -----------------------------------------------------
func (r *paymentRepository) GetPaymentByID(ctx context.Context, paymentID uint, userID uint) (*models.Payment, error) {
	var payment models.Payment

	err := r.db.WithContext(ctx).
		Joins("JOIN reservations ON reservations.id = payments.reservation_id").
		Where("payments.id = ? AND reservations.user_id = ?", paymentID, userID).
		Preload("Reservation").
		First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

// -----------------------------------------------------
// UPDATE FULL PAYMENT
// -----------------------------------------------------
func (r *paymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}
