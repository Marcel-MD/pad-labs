package models

const (
	PendingStatus   = "pending"
	ShippedStatus   = "shipped"
	CanceledStatus  = "canceled"
	CompletedStatus = "completed"
)

type Order struct {
	Base
	ProductId       string  `json:"product_id"`
	UserId          string  `json:"user_id"`
	Quantity        int64   `json:"quantity"`
	Cost            int64   `json:"cost"`
	Status          string  `json:"status"`
	ShippingAddress string  `json:"shipping_address"`
	Product         Product `json:"-" gorm:"foreignKey:ProductId"`
	User            User    `json:"-" gorm:"foreignKey:UserId"`
}

type OrderMessage struct {
	Base
	ProductId string `json:"product_id"`
	UserId    string `json:"user_id"`
	Quantity  int64  `json:"quantity"`
}

type UpdateOrder struct {
	ID             string `json:"id"`
	ProductOwnerId string `json:"product_owner_id"`
	Status         string `json:"status"`
	Cost           int64  `json:"cost"`
}
