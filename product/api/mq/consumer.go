package mq

import (
	"encoding/json"
	"product/config"
	"product/models"
	"product/services"
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
	producer       Producer
}

func NewConsumer(cfg config.Config, userService services.UserService, productService services.ProductService, producer Producer) (Consumer, error) {
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
		producer:       producer,
	}

	cons, err := rabbitmq.NewConsumer(
		conn,
		c.handleDelivery,
		models.ProductQueue,
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

	case models.CreateOrderMsgType:
		var order models.OrderMessage
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal order")
			return rabbitmq.NackDiscard
		}

		product, err := c.productService.FindById(order.ProductId)
		if err != nil {
			log.Error().Err(err).Msg("Failed to find product")
			return rabbitmq.Ack
		}

		if product.Stock < order.Quantity {
			log.Warn().Msg("Insufficient stock")
			c.producer.SendMsg(models.UpdateOrderMsgType, models.UpdateOrderMessage{
				ID:             order.ID,
				ProductOwnerId: product.OwnerId,
				Status:         models.CanceledStatus,
				Cost:           0,
			}, []string{models.OrderQueue})
			return rabbitmq.Ack
		}

		product.Stock -= order.Quantity
		product, err = c.productService.Update(product)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update product")
			return rabbitmq.Ack
		}

		c.producer.SendMsg(models.UpdateOrderMsgType, models.UpdateOrderMessage{
			ID:             order.ID,
			ProductOwnerId: product.OwnerId,
			Status:         models.ShippedStatus,
			Cost:           product.Price * order.Quantity,
		}, []string{models.OrderQueue})
	}

	return rabbitmq.Ack
}
