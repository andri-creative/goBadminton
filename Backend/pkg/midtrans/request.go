package midtrans

type ChargeRequest struct {
	PaymentType        string             `json:"payment_type"`
	TransactionDetails TransactionDetails `json:"transaction_details"`
	CustomerDetails    *CustomerDetails   `json:"customer_details,omitempty"`
	ItemDetails        []ItemDetail       `json:"item_details,omitempty"`
	BankTransfer       *BankTransfer      `json:"bank_transfer,omitempty"`
	EWallet            *EWallet           `json:"ewallet,omitempty"`
}

type TransactionDetails struct {
	OrderID  string `json:"order_id"`
	GrossAmt int64  `json:"gross_amount"`
}

type CustomerDetails struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
}

type ItemDetail struct {
	ID    string `json:"id"`
	Price int64  `json:"price"`
	Qty   int    `json:"quantity"`
	Name  string `json:"name"`
}

type BankTransfer struct {
	Bank string `json:"bank,omitempty"` // bca, bni, bri, etc.
}

type EWallet struct {
	Channel string `json:"channel,omitempty"` // gopay, shopeepay, etc.
}
