package models

type Product struct {
	Base
	OwnerId string `json:"owner_id"`
	Name    string `json:"name"`
	Price   int64  `json:"price"`
	Stock   int64  `json:"stock"`
	Owner   User   `json:"-" gorm:"foreignKey:OwnerId"`
}

type ProductMessage struct {
	Base
	OwnerId string `json:"owner_id"`
}
