package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CourtHandler struct {
	courtService services.CourtService
}

func NewCourtHandler(courtService services.CourtService) *CourtHandler {
	return &CourtHandler{courtService: courtService}
}

// GetAllCourts godoc
// @Summary Get all courts
// @Description Get list of all badminton courts
// @Tags courts
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /courts [get]
func (h *CourtHandler) GetAllCourts(c *gin.Context) {
	courts, err := h.courtService.GetAllCourts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get courts",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courts": courts,
	})
}

// GetAvailableCourts godoc
// @Summary Get available courts and timeslots
// @Description Get available courts with their available timeslots for a specific date
// @Tags courts
// @Accept json
// @Produce json
// @Param date query string true "Date in YYYY-MM-DD format"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /courts/available [get]
func (h *CourtHandler) GetAvailableCourts(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Date parameter is required",
		})
		return
	}

	availableCourts, err := h.courtService.GetAvailableCourts(c.Request.Context(), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get available courts",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"date":             date,
		"available_courts": availableCourts,
	})
}

// CheckAvailability godoc
// @Summary Check specific timeslot availability
// @Description Check if a specific timeslot is available for a court
// @Tags courts
// @Accept json
// @Produce json
// @Param request body models.CheckAvailabilityRequest true "Availability check data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /courts/check-availability [post]
func (h *CourtHandler) CheckAvailability(c *gin.Context) {
	var req models.CheckAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	isAvailable, err := h.courtService.CheckTimeSlotAvailability(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check availability",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_available": isAvailable,
		"date":         req.Date,
		"time_slot":    req.TimeSlot,
		"court_id":     req.CourtID,
	})
}

// GetCourtByID godoc
// @Summary Get court by ID
// @Description Get specific court details
// @Tags courts
// @Accept json
// @Produce json
// @Param id path int true "Court ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /courts/{id} [get]
func (h *CourtHandler) GetCourtByID(c *gin.Context) {
	var req struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid court ID",
		})
		return
	}

	court, err := h.courtService.GetCourtByID(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Court not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"court": court,
	})
}
