package models

type TimeSlot struct {
	Time     string `json:"time"`
	IsBooked bool   `json:"is_booked"`
}

type AvailableSlotResponse struct {
	CourtID   uint       `json:"court_id"`
	CourtName string     `json:"court_name"`
	Date      string     `json:"date"`
	TimeSlots []TimeSlot `json:"time_slots"`
}

type AvailableCourtRequest struct {
	Date string `json:"date" binding:"required"`
}
