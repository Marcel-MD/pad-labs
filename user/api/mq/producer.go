package mq

import (
	"encoding/json"
	"time"
	"user/config"
	"user/models"

	"github.com/rs/zerolog/log"
	"github.com/wagslane/go-rabbitmq"
)

type Producer interface {
	SendMsg(msgType string, msg interface{}, queues []string)
	Close() error
}

func NewProducer(cfg config.Config) (Producer, error) {
	log.Info().Msg("Creating new RabbitMQ producer")

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

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
	)
	if err != nil {
		return nil, err
	}

	prod := &producer{
		conn:      conn,
		publisher: publisher,
	}

	return prod, nil
}

type producer struct {
	conn      *rabbitmq.Conn
	publisher *rabbitmq.Publisher
}

func (s *producer) Close() error {
	s.publisher.Close()
	return s.conn.Close()
}

func (p *producer) SendMsg(msgType string, msg interface{}, queues []string) {
	go func() {
		data, err := json.Marshal(msg)
		if err != nil {
			log.Print("Failed to marshal message")
			return
		}

		message := models.Message{
			Type: msgType,
			Data: data,
		}

		body, err := json.Marshal(message)
		if err != nil {
			log.Print("Failed to marshal message")
			return
		}

		err = p.publisher.Publish(
			body,
			queues,
			rabbitmq.WithPublishOptionsContentType("application/json"),
		)
		if err != nil {
			log.Print("Failed to send message")
		}
	}()
}
