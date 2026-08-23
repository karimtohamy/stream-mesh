package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"stream-mesh/media-sync/internal/config"
)

type TranscodeHandler func(ctx context.Context, job TransCodeEvent) error

type Listener struct {
	client *RabbitClient
	cfg    *config.Config
}

func NewListener(client *RabbitClient, cfg *config.Config) *Listener {
	return &Listener{
		client: client,
		cfg:    cfg,
	}
}

func (l *Listener) StartTranscodeConsumer(ctx context.Context, handler TranscodeHandler) error {
	ch := l.client.Channel()

	if err := ch.Qos(l.cfg.App.TransmuxWorkers, 0, false); err != nil {
		return fmt.Errorf("failed to set prefetch QoS: %w", err)
	}

	var msgs <-chan amqp.Delivery
	var err error

	msgs, err = ch.Consume(
		l.cfg.RabbitMQ.TranscodeQueue,
		"media-sync-transcoder",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming from %s: %w", l.cfg.RabbitMQ.TranscodeQueue, err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-msgs:
				if !ok {
					return
				}

				var job TransCodeEvent
				if err := json.Unmarshal(d.Body, &job); err != nil {
					log.Printf("[CONSUMER ERROR] Invalid JSON payload: %v", err)
					_ = d.Nack(false, false)
					continue
				}

				if err := handler(ctx, job); err != nil {
					log.Printf("[TRANSCODE ERROR] Failed job for Media %s: %v", job.MediaId, err)
					_ = d.Nack(false, false)
					continue
				}

				_ = d.Ack(false)
			}
		}
	}()

	return nil
}