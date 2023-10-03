package models

const (
	PendingStatus   = "pending"
	ShippedStatus   = "shipped"
	CanceledStatus  = "canceled"
	CompletedStatus = "completed"
)

type OrderMessage struct {
	Base
	ProductId string `json:"product_id"`
	UserId    string `json:"user_id"`
	Quantity  int64  `json:"quantity"`
}

type UpdateOrderMessage struct {
	ID             string `json:"id"`
	ProductOwnerId string `json:"product_owner_id"`
	Status         string `json:"status"`
	Cost           int64  `json:"cost"`
}
