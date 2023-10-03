package models

type Product struct {
	Base
	OwnerId string  `json:"owner_id"`
	Owner   User    `json:"-" gorm:"foreignKey:OwnerId"`
	Orders  []Order `json:"-" gorm:"foreignKey:ProductId"`
}
