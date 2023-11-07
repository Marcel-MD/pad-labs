package models

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
