package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReservationHandler struct {
	reservationService services.ReservationService
}

func NewReservationHandler(reservationService services.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

// CreateReservation godoc
// @Summary Create a new reservation
// @Description Create a new badminton court reservation
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateReservationRequest true "Reservation data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reservations [post]
func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User must be logged in to make reservation",
		})
		return
	}

	var req models.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	reservation, err := h.reservationService.CreateReservation(
		c.Request.Context(),
		userID.(uint),
		&req,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Reservation created successfully",
		"reservation": reservation,
	})
}

// GetUserReservations godoc
// @Summary Get user reservations
// @Description Get all reservations for the authenticated user
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reservations [get]
func (h *ReservationHandler) GetUserReservations(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	reservations, err := h.reservationService.GetUserReservations(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reservations": reservations,
		"count":        len(reservations),
	})
}

// GetReservationByID godoc
// @Summary Get reservation by ID
// @Description Get specific reservation details
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Reservation ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /reservations/{id} [get]
func (h *ReservationHandler) GetReservationByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	reservationIDStr := c.Param("id")
	reservationID, err := strconv.ParseUint(reservationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid reservation ID",
		})
		return
	}

	reservation, err := h.reservationService.GetReservationByID(
		c.Request.Context(),
		uint(reservationID),
		userID.(uint),
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reservation": reservation,
	})
}

// CancelReservation godoc
// @Summary Cancel reservation
// @Description Cancel a pending reservation
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Reservation ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /reservations/{id}/cancel [put]
func (h *ReservationHandler) CancelReservation(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	reservationIDStr := c.Param("id")
	reservationID, err := strconv.ParseUint(reservationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid reservation ID",
		})
		return
	}

	err = h.reservationService.CancelReservation(
		c.Request.Context(),
		uint(reservationID),
		userID.(uint),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reservation cancelled successfully",
	})
}
