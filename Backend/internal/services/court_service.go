package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"time"
)

type CourtService interface {
	GetAllCourts(ctx context.Context) ([]models.CourtResponse, error)
	GetAvailableCourts(ctx context.Context, date string) ([]models.AvailableSlotResponse, error)
	CheckTimeSlotAvailability(ctx context.Context, req models.CheckAvailabilityRequest) (bool, error)
	GetCourtByID(ctx context.Context, id uint) (*models.CourtResponse, error)
}

type courtService struct {
	courtRepo repositories.CourtRepository
}

func NewCourtService(courtRepo repositories.CourtRepository) CourtService {
	return &courtService{courtRepo: courtRepo}
}

func (s *courtService) GetAllCourts(ctx context.Context) ([]models.CourtResponse, error) {
	courts, err := s.courtRepo.GetAllCourts(ctx)
	if err != nil {
		return nil, err
	}

	var courtResponses []models.CourtResponse
	for _, court := range courts {
		courtResponses = append(courtResponses, models.CourtResponse{
			ID:           court.ID,
			Name:         court.Name,
			Location:     court.Location,     // NEW
			PricePerHour: court.PricePerHour, // NEW
			Status:       court.Status,
		})
	}

	return courtResponses, nil
}

func (s *courtService) GetAvailableCourts(ctx context.Context, date string) ([]models.AvailableSlotResponse, error) {
	// Parse date
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	// Get all courts
	courts, err := s.courtRepo.GetAllCourts(ctx)
	if err != nil {
		return nil, err
	}

	var availableCourts []models.AvailableSlotResponse

	for _, court := range courts {
		// Get available time slots for this court
		availableSlots, err := s.courtRepo.GetAvailableTimeSlots(ctx, parsedDate, court.ID)
		if err != nil {
			return nil, err
		}

		// Convert to TimeSlot models
		var timeSlots []models.TimeSlot
		for _, slot := range availableSlots {
			timeSlots = append(timeSlots, models.TimeSlot{
				Time:     slot,
				IsBooked: false,
			})
		}

		// If no available slots, skip this court
		if len(timeSlots) == 0 {
			continue
		}

		availableCourt := models.AvailableSlotResponse{
			CourtID:   court.ID,
			CourtName: court.Name,
			Date:      date,
			TimeSlots: timeSlots,
		}

		availableCourts = append(availableCourts, availableCourt)
	}

	return availableCourts, nil
}

func (s *courtService) CheckTimeSlotAvailability(ctx context.Context, req models.CheckAvailabilityRequest) (bool, error) {
	// Parse date
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return false, err
	}

	// Check if the time slot is available for the court
	isAvailable, err := s.courtRepo.CheckCourtAvailability(ctx, parsedDate, req.TimeSlot, req.CourtID)
	if err != nil {
		return false, err
	}

	return isAvailable, nil
}

func (s *courtService) GetCourtByID(ctx context.Context, id uint) (*models.CourtResponse, error) {
	court, err := s.courtRepo.GetCourtByID(ctx, id)
	if err != nil {
		return nil, err
	}

	courtResponse := &models.CourtResponse{
		ID:     court.ID,
		Name:   court.Name,
		Status: court.Status,
	}

	return courtResponse, nil
}
