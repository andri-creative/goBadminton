package repositories

import (
	"backend/internal/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type ReservationRepository interface {
	CreateReservation(ctx context.Context, reservation *models.Reservation) error
	GetReservationByID(ctx context.Context, id uint) (*models.Reservation, error)
	GetUserReservations(ctx context.Context, userID uint) ([]models.Reservation, error)
	GetReservationsByDateAndCourt(ctx context.Context, date time.Time, courtID uint) ([]models.Reservation, error)
	UpdateReservationStatus(ctx context.Context, id uint, status string) error
	CheckExistingReservation(ctx context.Context, date time.Time, timeSlot string, courtID uint) (bool, error)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) CreateReservation(ctx context.Context, reservation *models.Reservation) error {
	return r.db.WithContext(ctx).Create(reservation).Error
}

func (r *reservationRepository) GetReservationByID(ctx context.Context, id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Court").
		First(&reservation, id).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) GetUserReservations(ctx context.Context, userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.WithContext(ctx).
		Preload("Court").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepository) GetReservationsByDateAndCourt(ctx context.Context, date time.Time, courtID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.WithContext(ctx).
		Where("reservation_date = ? AND court_id = ? AND status IN (?, ?)",
			date, courtID, "pending", "confirmed").
		Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepository) UpdateReservationStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.Reservation{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *reservationRepository) CheckExistingReservation(ctx context.Context, date time.Time, timeSlot string, courtID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Reservation{}).
		Where("court_id = ? AND reservation_date = ? AND time_slot = ? AND status IN (?, ?)",
			courtID, date, timeSlot, "pending", "confirmed").
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
