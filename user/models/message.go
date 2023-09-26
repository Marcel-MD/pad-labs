package models

const (
	CreateUserMsgType = "CreateUser"
)

const (
	ProductQueue = "product"
	OrderQueue   = "order"
)

type Message struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}
