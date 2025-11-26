package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"errors"
	"time"
)

type ReservationService interface {
	CreateReservation(ctx context.Context, userID uint, req *models.CreateReservationRequest) (*models.ReservationResponse, error)
	GetUserReservations(ctx context.Context, userID uint) ([]models.ReservationResponse, error)
	GetReservationByID(ctx context.Context, reservationID uint, userID uint) (*models.ReservationResponse, error)
	CancelReservation(ctx context.Context, reservationID uint, userID uint) error
}

type reservationService struct {
	reservationRepo repositories.ReservationRepository
	courtRepo       repositories.CourtRepository
}

func NewReservationService(
	reservationRepo repositories.ReservationRepository,
	courtRepo repositories.CourtRepository,
) ReservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		courtRepo:       courtRepo,
	}
}

func calculateDuration(timeSlot string) int {
	return 1 // Default 1 hour for now
}

func (s *reservationService) CreateReservation(ctx context.Context, userID uint, req *models.CreateReservationRequest) (*models.ReservationResponse, error) {
	// Parse date
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format. Use YYYY-MM-DD")
	}

	// Check if court exists
	_, err = s.courtRepo.GetCourtByID(ctx, req.CourtID)
	if err != nil {
		return nil, errors.New("court not found")
	}

	// Check court availability
	isAvailable, err := s.reservationRepo.CheckExistingReservation(ctx, parsedDate, req.TimeSlot, req.CourtID)
	if err != nil {
		return nil, errors.New("failed to check availability")
	}
	if isAvailable {
		return nil, errors.New("selected timeslot is already booked")
	}

	// Validate timeslot format (basic validation)
	if len(req.TimeSlot) != 11 || req.TimeSlot[5] != '-' {
		return nil, errors.New("invalid timeslot format. Use HH:MM-HH:MM")
	}

	// Get court details including PRICE
	court, err := s.courtRepo.GetCourtByID(ctx, req.CourtID)
	if err != nil {
		return nil, errors.New("court not found")
	}

	// CALCULATE DURATION from time slot
	duration := calculateDuration(req.TimeSlot) // Function baru

	// CALCULATE TOTAL AMOUNT
	totalAmount := court.PricePerHour * float64(duration)

	// Create reservation
	reservation := &models.Reservation{
		UserID:          userID,
		CourtID:         req.CourtID,
		ReservationDate: parsedDate,
		TimeSlot:        req.TimeSlot,
		DurationHours:   duration,    // NEW
		TotalAmount:     totalAmount, // NEW
		Status:          "pending",   // Will be confirmed after payment
	}

	err = s.reservationRepo.CreateReservation(ctx, reservation)
	if err != nil {
		return nil, errors.New("failed to create reservation")
	}

	// Get the created reservation with relationships
	createdReservation, err := s.reservationRepo.GetReservationByID(ctx, reservation.ID)
	if err != nil {
		return nil, errors.New("failed to fetch created reservation")
	}

	// Convert to response - ✅ TAMBAHKAN DURATION_HOURS & TOTAL_AMOUNT
	reservationResponse := &models.ReservationResponse{
		ID:              createdReservation.ID,
		UserID:          createdReservation.UserID,
		CourtID:         createdReservation.CourtID,
		CourtName:       createdReservation.Court.Name,
		ReservationDate: createdReservation.ReservationDate.Format("2006-01-02"),
		TimeSlot:        createdReservation.TimeSlot,
		DurationHours:   createdReservation.DurationHours, // ✅ NEW
		TotalAmount:     createdReservation.TotalAmount,   // ✅ NEW
		Status:          createdReservation.Status,
		CreatedAt:       createdReservation.CreatedAt,
	}

	return reservationResponse, nil
}

func (s *reservationService) GetUserReservations(ctx context.Context, userID uint) ([]models.ReservationResponse, error) {
	reservations, err := s.reservationRepo.GetUserReservations(ctx, userID)
	if err != nil {
		return nil, errors.New("failed to get user reservations")
	}

	var reservationResponses []models.ReservationResponse
	for _, reservation := range reservations {
		reservationResponse := models.ReservationResponse{
			ID:              reservation.ID,
			UserID:          reservation.UserID,
			CourtID:         reservation.CourtID,
			CourtName:       reservation.Court.Name,
			ReservationDate: reservation.ReservationDate.Format("2006-01-02"),
			TimeSlot:        reservation.TimeSlot,
			DurationHours:   reservation.DurationHours, // ✅ NEW
			TotalAmount:     reservation.TotalAmount,   // ✅ NEW
			Status:          reservation.Status,
			CreatedAt:       reservation.CreatedAt,
		}
		reservationResponses = append(reservationResponses, reservationResponse)
	}

	return reservationResponses, nil
}

func (s *reservationService) GetReservationByID(ctx context.Context, reservationID uint, userID uint) (*models.ReservationResponse, error) {
	reservation, err := s.reservationRepo.GetReservationByID(ctx, reservationID)
	if err != nil {
		return nil, errors.New("reservation not found")
	}

	// Check if reservation belongs to user
	if reservation.UserID != userID {
		return nil, errors.New("unauthorized to access this reservation")
	}

	reservationResponse := &models.ReservationResponse{
		ID:              reservation.ID,
		UserID:          reservation.UserID,
		CourtID:         reservation.CourtID,
		CourtName:       reservation.Court.Name,
		ReservationDate: reservation.ReservationDate.Format("2006-01-02"),
		TimeSlot:        reservation.TimeSlot,
		DurationHours:   reservation.DurationHours, // ✅ NEW
		TotalAmount:     reservation.TotalAmount,   // ✅ NEW
		Status:          reservation.Status,
		CreatedAt:       reservation.CreatedAt,
	}

	return reservationResponse, nil
}

func (s *reservationService) CancelReservation(ctx context.Context, reservationID uint, userID uint) error {
	// First get the reservation to check ownership
	reservation, err := s.reservationRepo.GetReservationByID(ctx, reservationID)
	if err != nil {
		return errors.New("reservation not found")
	}

	// Check if reservation belongs to user
	if reservation.UserID != userID {
		return errors.New("unauthorized to cancel this reservation")
	}

	// Check if reservation can be cancelled (only pending reservations)
	if reservation.Status != "pending" {
		return errors.New("only pending reservations can be cancelled")
	}

	// Update reservation status to cancelled
	err = s.reservationRepo.UpdateReservationStatus(ctx, reservationID, "cancelled")
	if err != nil {
		return errors.New("failed to cancel reservation")
	}

	return nil
}
