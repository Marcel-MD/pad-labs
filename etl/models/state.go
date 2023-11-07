package models

import "time"

const (
	UserDB    = "user-db"
	ProductDB = "product-db"
	OrderDB   = "order-db"
)

type State struct {
	Base

	DatabaseName string    `json:"database_name" gorm:"uniqueIndex"`
	LastSync     time.Time `json:"last_sync"`
}
