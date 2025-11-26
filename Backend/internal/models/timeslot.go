package models

type TimeSlot struct {
	Time     string `json:"time"`      // Format: "10:00-11:00"
	IsBooked bool   `json:"is_booked"` // Available or not
}

type AvailableSlotResponse struct {
	CourtID   uint       `json:"court_id"`
	CourtName string     `json:"court_name"`
	Date      string     `json:"date"` // Format: "2006-01-02"
	TimeSlots []TimeSlot `json:"time_slots"`
}

type AvailableCourtRequest struct {
	Date string `json:"date" binding:"required"` // Format: "2006-01-02"
}
