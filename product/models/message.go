package models

const (
	CreateUserMsgType = "CreateUser"

	CreateProductMsgType = "CreateProduct"
	UpdateProductMsgType = "UpdateProduct"
	DeleteProductMsgType = "DeleteProduct"

	CreateOrderMsgType = "CreateOrder"
	UpdateOrderMsgType = "UpdateOrder"
)

const (
	ProductQueue = "product"
	OrderQueue   = "order"
)

type Message struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}
