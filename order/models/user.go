package models

type User struct {
	Base
	Products []Product `json:"-" gorm:"foreignKey:OwnerId;constraint:OnDelete:CASCADE;"`
	Orders   []Order   `json:"-" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}
