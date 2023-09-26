package models

const (
	UserRole  = "user"
	AdminRole = "admin"
)

type User struct {
	Base

	Products []Product `json:"-" gorm:"foreignKey:OwnerId;constraint:OnDelete:CASCADE;"`
}
