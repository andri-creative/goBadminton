package midtrans

type ChargeResponse struct {
	Token       string     `json:"token"`
	RedirectURL string     `json:"redirect_url"`
	OrderID     string     `json:"order_id"`
	StatusCode  string     `json:"status_code"`
	VaNumbers   []VaNumber `json:"va_numbers"`
}

type Notification struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
	PaymentType       string `json:"payment_type"`
}

type TransactionStatusResponse struct {
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
}

type VaNumber struct {
	Bank     string `json:"bank"`
	VaNumber string `json:"va_number"`
}
