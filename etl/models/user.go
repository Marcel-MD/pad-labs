package models

import "gorm.io/datatypes"

type User struct {
	Base

	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`

	Roles datatypes.JSONSlice[string] `json:"roles"`
}
