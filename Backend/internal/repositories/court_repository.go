package repositories

import (
	"backend/internal/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type CourtRepository interface {
	GetAllCourts(ctx context.Context) ([]models.Court, error)
	GetCourtByID(ctx context.Context, id uint) (*models.Court, error)
	GetAvailableTimeSlots(ctx context.Context, date time.Time, courtID uint) ([]string, error)
	CheckCourtAvailability(ctx context.Context, date time.Time, timeSlot string, courtID uint) (bool, error)
	GetReservedSlots(ctx context.Context, date time.Time) ([]models.Reservation, error)
}

type courtRepository struct {
	db *gorm.DB
}

func NewCourtRepository(db *gorm.DB) CourtRepository {
	return &courtRepository{db: db}
}

func (r *courtRepository) GetAllCourts(ctx context.Context) ([]models.Court, error) {
	var courts []models.Court
	err := r.db.WithContext(ctx).Where("status = ?", "active").Find(&courts).Error
	if err != nil {
		return nil, err
	}
	return courts, nil
}

func (r *courtRepository) GetCourtByID(ctx context.Context, id uint) (*models.Court, error) {
	var court models.Court
	err := r.db.WithContext(ctx).First(&court, id).Error
	if err != nil {
		return nil, err
	}
	return &court, nil
}

func (r *courtRepository) GetAvailableTimeSlots(ctx context.Context, date time.Time, courtID uint) ([]string, error) {
	// Define all possible time slots
	allTimeSlots := []string{
		"07:00-08:00", "08:00-09:00", "09:00-10:00", "10:00-11:00",
		"11:00-12:00", "12:00-13:00", "13:00-14:00", "14:00-15:00",
		"15:00-16:00", "16:00-17:00", "17:00-18:00", "18:00-19:00",
		"19:00-20:00", "20:00-21:00",
	}

	// Get reserved slots for this court and date
	var reservedSlots []string
	err := r.db.WithContext(ctx).
		Model(&models.Reservation{}).
		Where("court_id = ? AND reservation_date = ? AND status IN (?, ?)",
			courtID, date, "pending", "confirmed").
		Pluck("time_slot", &reservedSlots).Error

	if err != nil {
		return nil, err
	}

	// Create a map of reserved slots for quick lookup
	reservedMap := make(map[string]bool)
	for _, slot := range reservedSlots {
		reservedMap[slot] = true
	}

	// Filter available slots
	var availableSlots []string
	for _, slot := range allTimeSlots {
		if !reservedMap[slot] {
			availableSlots = append(availableSlots, slot)
		}
	}

	return availableSlots, nil
}

func (r *courtRepository) CheckCourtAvailability(ctx context.Context, date time.Time, timeSlot string, courtID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Reservation{}).
		Where("court_id = ? AND reservation_date = ? AND time_slot = ? AND status IN (?, ?)",
			courtID, date, timeSlot, "pending", "confirmed").
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *courtRepository) GetReservedSlots(ctx context.Context, date time.Time) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.WithContext(ctx).
		Where("reservation_date = ? AND status IN (?, ?)", date, "pending", "confirmed").
		Find(&reservations).Error

	if err != nil {
		return nil, err
	}
	return reservations, nil
}
