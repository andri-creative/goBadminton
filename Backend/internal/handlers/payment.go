package handlers

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService  services.MidtransService
	reservationRepo repositories.ReservationRepository
	userRepo        repositories.UserRepository
	paymentRepo     repositories.PaymentRepository
}

func NewPaymentHandler(
	paymentService services.MidtransService,
	reservationRepo repositories.ReservationRepository,
	userRepo repositories.UserRepository,
	paymentRepo repositories.PaymentRepository,
) *PaymentHandler {
	return &PaymentHandler{
		paymentService:  paymentService,
		reservationRepo: reservationRepo,
		userRepo:        userRepo,
		paymentRepo:     paymentRepo,
	}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get Reservation
	reservation, err := h.reservationRepo.GetReservationByID(c, req.ReservationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reservation not found"})
		return
	}

	// Get User
	user, err := h.userRepo.GetUserByID(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Create Payment - SEKARANG menggunakan PaymentResponse
	paymentResp, err := h.paymentService.CreatePayment(c, reservation, user, req.PaymentMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response berdasarkan payment method
	response := gin.H{
		"message":        "Payment created successfully",
		"order_id":       paymentResp.OrderID,
		"amount":         paymentResp.Amount,
		"status":         paymentResp.Status,
		"payment_method": req.PaymentMethod,
	}

	// Tambahkan fields berdasarkan jenis payment
	if paymentResp.SnapToken != "" {
		response["token"] = paymentResp.SnapToken
		response["redirect_url"] = paymentResp.RedirectURL
		fmt.Printf("SNAP Payment - Token: %s\n", paymentResp.SnapToken)
	}

	if paymentResp.VaNumber != "" {
		response["va_number"] = paymentResp.VaNumber
		response["va_bank"] = paymentResp.VaBank
		fmt.Printf("Bank Transfer - VA: %s, Bank: %s\n", paymentResp.VaNumber, paymentResp.VaBank)
	}

	c.JSON(http.StatusCreated, response)
}

// Handler lainnya tetap sama...
func (h *PaymentHandler) HandlePaymentNotification(c *gin.Context) {
	// Log everything about the request
	fmt.Printf("=== RAW REQUEST INSPECTION ===\n")
	fmt.Printf("Content-Type: %s\n", c.GetHeader("Content-Type"))
	fmt.Printf("Content-Length: %s\n", c.GetHeader("Content-Length"))

	// Read raw body
	body, err := c.GetRawData()
	if err != nil {
		fmt.Printf("ERROR reading raw body: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read request body: " + err.Error()})
		return
	}

	fmt.Printf("Raw body as string: %s\n", string(body))
	fmt.Printf("Raw body length: %d\n", len(body))

	// ✅ FIX: Handle empty body (Midtrans test notification)
	if len(body) == 0 {
		fmt.Printf("INFO: Empty request body received - Midtrans test notification\n")
		// Return 200 OK untuk test notification
		c.JSON(http.StatusOK, gin.H{"message": "Test notification received successfully"})
		return
	}

	// Try to parse as JSON
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Printf("ERROR unmarshaling JSON: %v\n", err)
		fmt.Printf("Error details: %s\n", err.Error())

		// Try to see what's wrong with the JSON
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			fmt.Printf("Syntax error at offset %d\n", syntaxErr.Offset)
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload: " + err.Error()})
		return
	}

	fmt.Printf("Successfully parsed payload: %+v\n", payload)

	// ✅ FIX: Handle test notifications from Midtrans
	if orderID, exists := payload["order_id"]; exists {
		orderIDStr := fmt.Sprintf("%v", orderID)
		if strings.Contains(orderIDStr, "payment_notif_test") || strings.Contains(orderIDStr, "test") {
			fmt.Printf("INFO: Midtrans test notification - OrderID: %s\n", orderIDStr)
			c.JSON(http.StatusOK, gin.H{"message": "Test notification processed successfully"})
			return
		}
	}

	// Process the notification
	if err := h.paymentService.HandleNotification(c.Request.Context(), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification processed successfully"})
}

func (h *PaymentHandler) GetUserPayments(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	payments, err := h.paymentRepo.GetUserPayments(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
		"count":    len(payments),
	})
}

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	userID, _ := c.Get("userID")
	idStr := c.Param("id")

	paymentID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	payment, err := h.paymentRepo.GetPaymentByID(c, uint(paymentID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment": payment})
}
