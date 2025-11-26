package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransService interface {
	CreatePayment(ctx context.Context, reservation *models.Reservation, user *models.User, paymentMethod string) (*PaymentResponse, error)
	HandleNotification(ctx context.Context, payload map[string]interface{}) error
}

type PaymentResponse struct {
	SnapToken   string                  `json:"snap_token,omitempty"`
	VaNumber    string                  `json:"va_number,omitempty"`
	VaBank      string                  `json:"va_bank,omitempty"`
	RedirectURL string                  `json:"redirect_url,omitempty"`
	CoreAPIResp *coreapi.ChargeResponse `json:"core_api_response,omitempty"`
	SnapResp    *snap.Response          `json:"snap_response,omitempty"`
	OrderID     string                  `json:"order_id"`
	Amount      int64                   `json:"amount"`
	Status      string                  `json:"status"`
}

type midtransService struct {
	coreClient      coreapi.Client
	snapClient      snap.Client
	paymentRepo     repositories.PaymentRepository
	reservationRepo repositories.ReservationRepository
}

func NewMidtransService(paymentRepo repositories.PaymentRepository, reservationRepo repositories.ReservationRepository) MidtransService {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		panic("MIDTRANS_SERVER_KEY is not set")
	}

	// Initialize both clients
	coreClient := coreapi.Client{}
	coreClient.New(serverKey, midtrans.Sandbox)

	snapClient := snap.Client{}
	snapClient.New(serverKey, midtrans.Sandbox)

	return &midtransService{
		coreClient:      coreClient,
		snapClient:      snapClient,
		paymentRepo:     paymentRepo,
		reservationRepo: reservationRepo,
	}
}

func (s *midtransService) CreatePayment(ctx context.Context, reservation *models.Reservation, user *models.User, paymentMethod string) (*PaymentResponse, error) {
	amount := int64(reservation.TotalAmount)
	orderID := fmt.Sprintf("ORDER-%d-%d", reservation.ID, reservation.CreatedAt.Unix())

	fmt.Printf("🎯 Creating payment for reservation %d, method: %s, amount: %d\n",
		reservation.ID, paymentMethod, amount)

	// Gunakan SNAP untuk semua payment method kecuali bank transfer
	if paymentMethod != "bank_transfer" {
		fmt.Printf("🔄 Using SNAP for payment method: %s\n", paymentMethod)
		return s.createSnapPayment(ctx, reservation, user, paymentMethod, amount, orderID)
	}

	// Untuk bank transfer, gunakan Core API
	fmt.Printf("🔄 Using CoreAPI for bank transfer\n")
	return s.createCoreAPIPayment(ctx, reservation, user, paymentMethod, amount, orderID)
}

// Untuk Gopay, QRIS, Credit Card, etc. (Snap Popup)
func (s *midtransService) createSnapPayment(ctx context.Context, reservation *models.Reservation, user *models.User, paymentMethod string, amount int64, orderID string) (*PaymentResponse, error) {

	// Map payment method to Snap payment type
	var enabledPayments []snap.SnapPaymentType
	switch paymentMethod {
	case "gopay":
		enabledPayments = []snap.SnapPaymentType{snap.PaymentTypeGopay}
	case "qris":
		enabledPayments = []snap.SnapPaymentType{"qris"} // ✅ FIXED
	case "credit_card":
		enabledPayments = []snap.SnapPaymentType{snap.PaymentTypeCreditCard}
	case "shopeepay":
		enabledPayments = []snap.SnapPaymentType{snap.PaymentTypeShopeepay}
	default:
		// Default: enable semua payment methods
		enabledPayments = []snap.SnapPaymentType{
			snap.PaymentTypeGopay,
			"qris", // ✅ FIXED
			snap.PaymentTypeCreditCard,
			snap.PaymentTypeShopeepay,
		}
	}

	fmt.Printf("✅ Enabled payments for Snap: %v\n", enabledPayments)

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
		EnabledPayments: enabledPayments,
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
			Phone: user.Phone,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    fmt.Sprintf("COURT-%d", reservation.CourtID),
				Price: amount,
				Qty:   1,
				Name:  fmt.Sprintf("Court Booking - %s", reservation.TimeSlot),
			},
		},
	}

	fmt.Printf("🔄 Sending Snap request for order: %s\n", orderID)
	snapResp, err := s.snapClient.CreateTransaction(snapReq)
	if err != nil {
		fmt.Printf("❌ ERROR creating Snap transaction: %v\n", err)
		return nil, fmt.Errorf("failed to create Snap transaction: %v", err)
	}

	fmt.Printf("✅ Snap response received. Token: %s, RedirectURL: %s\n",
		snapResp.Token, snapResp.RedirectURL)

	// Save to database
	payment := &models.Payment{
		ReservationID:   reservation.ID,
		Amount:          float64(amount),
		Status:          "pending",
		PaymentMethod:   paymentMethod,
		MidtransOrderID: orderID,
		VaNumber:        "", // Tidak ada VA untuk non-bank transfer
		VaBank:          "",
	}

	if err := s.paymentRepo.CreatePayment(ctx, payment); err != nil {
		fmt.Printf("❌ ERROR saving payment to database: %v\n", err)
		return nil, fmt.Errorf("failed to save payment: %v", err)
	}

	fmt.Printf("✅ Payment saved to database with ID: %d\n", payment.ID)

	return &PaymentResponse{
		SnapToken:   snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
		SnapResp:    snapResp,
		OrderID:     orderID,
		Amount:      amount,
		Status:      "pending",
	}, nil
}

// Untuk Bank Transfer (Core API)
func (s *midtransService) createCoreAPIPayment(ctx context.Context, reservation *models.Reservation, user *models.User, paymentMethod string, amount int64, orderID string) (*PaymentResponse, error) {

	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeBankTransfer,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
		BankTransfer: &coreapi.BankTransferDetails{
			Bank: midtrans.BankBca,
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
			Phone: user.Phone,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    fmt.Sprintf("COURT-%d", reservation.CourtID),
				Price: amount,
				Qty:   1,
				Name:  fmt.Sprintf("Court Booking - %s", reservation.TimeSlot),
			},
		},
	}

	fmt.Printf("🔄 Sending CoreAPI request for bank transfer, order: %s\n", orderID)
	coreResp, err := s.coreClient.ChargeTransaction(chargeReq)
	if err != nil {
		fmt.Printf("❌ ERROR creating CoreAPI transaction: %v\n", err)
		return nil, fmt.Errorf("failed to create Core API transaction: %v", err)
	}

	fmt.Printf("✅ CoreAPI response received. Status: %s\n", coreResp.StatusMessage)

	// Save VA Info
	vaNumber := ""
	vaBank := ""

	if len(coreResp.VaNumbers) > 0 {
		vaNumber = coreResp.VaNumbers[0].VANumber
		vaBank = coreResp.VaNumbers[0].Bank
		fmt.Printf("✅ VA Number: %s, Bank: %s\n", vaNumber, vaBank)
	}

	payment := &models.Payment{
		ReservationID:   reservation.ID,
		Amount:          float64(amount),
		Status:          "pending",
		PaymentMethod:   paymentMethod,
		MidtransOrderID: orderID,
		VaNumber:        vaNumber,
		VaBank:          vaBank,
	}

	if err := s.paymentRepo.CreatePayment(ctx, payment); err != nil {
		fmt.Printf("❌ ERROR saving payment to database: %v\n", err)
		return nil, fmt.Errorf("failed to save payment: %v", err)
	}

	fmt.Printf("✅ Bank transfer payment saved to database with ID: %d\n", payment.ID)

	return &PaymentResponse{
		VaNumber:    vaNumber,
		VaBank:      vaBank,
		CoreAPIResp: coreResp,
		OrderID:     orderID,
		Amount:      amount,
		Status:      "pending",
	}, nil
}

// FIXED PARSING NOTIFICATION
func (s *midtransService) HandleNotification(ctx context.Context, payload map[string]interface{}) error {
	fmt.Printf("=== MIDTRANS NOTIFICATION RECEIVED ===\n")
	fmt.Printf("Raw payload: %+v\n", payload)

	jsonBytes, _ := json.Marshal(payload)
	var notif coreapi.TransactionStatusResponse
	if err := json.Unmarshal(jsonBytes, &notif); err != nil {
		fmt.Printf("ERROR parsing notification: %v\n", err)
		return fmt.Errorf("failed to parse Midtrans notification: %v", err)
	}

	fmt.Printf("Parsed notification - OrderID: %s, Status: %s\n", notif.OrderID, notif.TransactionStatus)

	var newStatus string
	switch notif.TransactionStatus {
	case "settlement":
		newStatus = "paid"
		fmt.Printf("Setting payment status to: %s\n", newStatus)
	case "pending":
		newStatus = "pending"
	case "expire":
		newStatus = "expired"
	case "cancel", "deny":
		newStatus = "failed"
	default:
		newStatus = "failed"
	}

	// Update payment
	fmt.Printf("Updating payment for OrderID: %s to status: %s\n", notif.OrderID, newStatus)
	if err := s.paymentRepo.UpdatePaymentStatus(ctx, notif.OrderID, newStatus); err != nil {
		fmt.Printf("ERROR updating payment status: %v\n", err)
		return fmt.Errorf("failed to update payment status: %v", err)
	}
	fmt.Printf("Payment updated successfully\n")

	// Update reservation if paid
	if newStatus == "paid" {
		fmt.Printf("Payment is PAID, updating reservation...\n")
		pay, err := s.paymentRepo.GetPaymentByOrderID(ctx, notif.OrderID)
		if err != nil {
			fmt.Printf("ERROR getting payment: %v\n", err)
			return fmt.Errorf("failed to get payment by order ID: %v", err)
		}
		fmt.Printf("Found payment with ReservationID: %d\n", pay.ReservationID)

		if err := s.reservationRepo.UpdateReservationStatus(ctx, pay.ReservationID, "confirmed"); err != nil {
			fmt.Printf("ERROR updating reservation: %v\n", err)
			return fmt.Errorf("failed to update reservation status: %v", err)
		}
		fmt.Printf("Reservation %d updated to 'confirmed'\n", pay.ReservationID)
	}

	fmt.Printf("=== NOTIFICATION PROCESSING COMPLETE ===\n")
	return nil
}
