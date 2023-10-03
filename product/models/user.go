package models

type User struct {
	Base
	Products []Product `json:"-" gorm:"foreignKey:OwnerId"`
}
