package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"stream-mesh/media-sync/internal/config"
)

type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     *config.Config
}

type RabbitTemplate struct {
	client *RabbitClient
}


func NewRabbitClient(cfg *config.Config) (*RabbitClient, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ.URL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}
	client := &RabbitClient{
		conn:    conn,
		channel: ch,
		cfg:     cfg,
	}
	if err := client.declareTopology(); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

func (r *RabbitClient) declareTopology() error {
	err := r.channel.ExchangeDeclare(
		r.cfg.RabbitMQ.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	dlxName := r.cfg.RabbitMQ.Exchange + ".dlx"
	err = r.channel.ExchangeDeclare(
		dlxName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	dlqName := r.cfg.RabbitMQ.TranscodeQueue + ".dlq"
	dlqKey := "transcode.dead"
	_, err = r.channel.QueueDeclare(
		dlqName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = r.channel.QueueBind(
		dlqName,
		dlqKey,
		dlxName,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	queueArgs := amqp.Table{
		"x-dead-letter-exchange":    dlxName,
		"x-dead-letter-routing-key": dlqKey,
	}

	_, err = r.channel.QueueDeclare(
		r.cfg.RabbitMQ.TranscodeQueue,
		true,
		false,
		false,
		false,
		queueArgs,
	)
	if err != nil {
		return err
	}

	return r.channel.QueueBind(
		r.cfg.RabbitMQ.TranscodeQueue,
		r.cfg.RabbitMQ.TranscodeReqKey,
		r.cfg.RabbitMQ.Exchange,
		false,
		nil,
	)
}
func (r *RabbitClient) Channel() *amqp.Channel {
	return r.channel
}

func (r *RabbitClient) Config() *config.Config {
	return r.cfg
}

func (r *RabbitClient) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

func NewRabbitTemplate(client *RabbitClient) *RabbitTemplate {
	return &RabbitTemplate{client: client}
}

func (t *RabbitTemplate) ConvertAndSend(ctx context.Context, exchange, routingKey string, message any) error {
	var body []byte
	var err error

	switch v := message.(type) {
	case []byte:
		body = v
	case string:
		body = []byte(v)
	default:
		body, err = json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal payload to json: %w", err)
		}
	}

	return t.client.channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now().UTC(),
			Body:         body,
		},
	)
}