package model

var (
	OrderCurrencyUSD   = "USD"
	OrderIntentCapture = "CAPTURE"

	OrderStatusCreated             = "CREATED"
	OrderStatusCompleted           = "COMPLETED"
	OrderStatusApproved            = "APPROVED"
	OrderStatusPayerActionRequired = "PAYER_ACTION_REQUIRED"

	OrderRelApprove     = "approve"
	OrderRelUpdate      = "update"
	OrderRelSelf        = "self"
	OrderRefPayerAction = "payer-action"
)

type UnitAmount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type PurchaseUnitItem struct {
	Name       string     `json:"name"`
	Quantity   string     `json:"quantity"`
	UnitAmount UnitAmount `json:"unit_amount"`
}

type PurchaseUnitAmountBreakdown struct {
	ItemTotal UnitAmount `json:"item_total"`
}

type PurchaseUnitAmount struct {
	CurrencyCode string                      `json:"currency_code"`
	Value        string                      `json:"value"`
	Breakdown    PurchaseUnitAmountBreakdown `json:"breakdown"`
}

type PurchaseUnit struct {
	Items  []PurchaseUnitItem `json:"items"`
	Amount PurchaseUnitAmount `json:"amount"`
}

type CreateOrderLinkResponse struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type CreateOrderRequest struct {
	Intent        string         `json:"intent"`
	PurchaseUnits []PurchaseUnit `json:"purchase_units"`
	PaymentSource struct {
		Paypal struct {
			ExperienceContext struct {
				ReturnUrl string `json:"return_url"`
				CancelUrl string `json:"cancel_url"`
			} `json:"experience_context"`
		} `json:"paypal"`
	} `json:"payment_source"`
}

type CreateOrderResponse struct {
	ID     string                    `json:"id"`
	Status string                    `json:"status"`
	Links  []CreateOrderLinkResponse `json:"links"`
}

type CaptureOrderPayerName struct {
	GivenName string `json:"given_name"`
	Surname   string `json:"surname"`
}

type CaptureOrderPayer struct {
	Name         CaptureOrderPayerName `json:"name"`
	EmailAddress string                `json:"email_address"`
	PayerID      string                `json:"payer_id"`
}

type CaptureOrderResponse struct {
	ID     string            `json:"id"`
	Status string            `json:"status"`
	Payer  CaptureOrderPayer `json:"payer"`
	// to know more fields - ref: https://developer.paypal.com/docs/api/orders/v2/#orders_capture
}

type GetOrderDetailResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
