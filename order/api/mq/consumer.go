package mq

import (
	"encoding/json"
	"order/config"
	"order/models"
	"order/services"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/wagslane/go-rabbitmq"
)

type Consumer interface {
	Close() error
}

type consumer struct {
	conn           *rabbitmq.Conn
	consumer       *rabbitmq.Consumer
	userService    services.UserService
	productService services.ProductService
	orderService   services.OrderService
}

func NewConsumer(cfg config.Config, userService services.UserService, productService services.ProductService, orderService services.OrderService) (Consumer, error) {
	log.Info().Msg("Creating new RabbitMQ consumer")

	const retries = 5
	var conn *rabbitmq.Conn
	var err error

	for i := 0; i < retries; i++ {
		conn, err = rabbitmq.NewConn(
			cfg.RabbitMQUrl,
			rabbitmq.WithConnectionOptionsLogging,
		)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to connect to RabbitMQ. Retrying...")
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	c := &consumer{
		conn:           conn,
		userService:    userService,
		productService: productService,
		orderService:   orderService,
	}

	cons, err := rabbitmq.NewConsumer(
		conn,
		c.handleDelivery,
		models.OrderQueue,
		rabbitmq.WithConsumerOptionsQueueDurable,
	)
	if err != nil {
		return nil, err
	}

	c.consumer = cons

	return c, nil
}

func (c *consumer) Close() error {
	c.consumer.Close()
	return c.conn.Close()
}

func (c *consumer) handleDelivery(d rabbitmq.Delivery) rabbitmq.Action {
	log.Debug().Msg("Processing message delivery")

	var msg models.Message

	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal message")
		return rabbitmq.NackDiscard
	}

	switch msg.Type {
	case models.CreateUserMsgType:
		var user models.User
		err := json.Unmarshal(msg.Data, &user)
		if err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal user")
			return rabbitmq.NackDiscard
		}

		_, err = c.userService.Create(user)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create user")
			return rabbitmq.Ack
		}

	case models.CreateProductMsgType:
		var product models.Product
		err := json.Unmarshal(msg.Data, &product)
		if err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal product")
			return rabbitmq.NackDiscard
		}

		_, err = c.productService.Create(product)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create product")
			return rabbitmq.Ack
		}

	case models.UpdateOrderMsgType:
		var order models.UpdateOrder
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal order")
			return rabbitmq.NackDiscard
		}

		_, err = c.orderService.Update(order)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update order")
			return rabbitmq.Ack
		}
	}

	return rabbitmq.Ack
}
