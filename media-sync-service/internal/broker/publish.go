package broker

import (
	"context"
	"fmt"
	"stream-mesh/media-sync/internal/config"
	"time"
)

type Publisher struct {
	template *RabbitTemplate
	cfg      *config.Config
}

func NewPublisher(template *RabbitTemplate, cfg *config.Config) *Publisher {
	return &Publisher{
		template: template,
		cfg:      cfg,
	}
}

func (p *Publisher) PublishTranscodeCompleted(ctx context.Context, MediaId string, manifestURL string, targetBucket string) error {
	payload := TransmuxEventResponse{
		MediaId:      MediaId,
		ManifestURL:  manifestURL,
		TargetBucket: targetBucket,
		CompletedAt:  time.Now().UTC(),
	}

	if err := p.template.ConvertAndSend(ctx, p.cfg.RabbitMQ.Exchange, p.cfg.RabbitMQ.TranscodeCmpKey, payload); err != nil {
		return fmt.Errorf("failed to publish transcode completed event: %w", err)
	}

	return nil
}
func (p *Publisher) PublishUsageTick(ctx context.Context, eventID, streamID, userID string, durationSec int, costPerMin float64) error {
	payload := UsageTickPayload{
		EventID:         eventID,
		StreamID:        streamID,
		UserID:          userID,
		DurationSeconds: durationSec,
		CostPerMinute:   costPerMin,
		Timestamp:       time.Now().UTC(),
	}

	if err := p.template.ConvertAndSend(ctx, p.cfg.RabbitMQ.Exchange, p.cfg.RabbitMQ.UsageTickKey, payload); err != nil {
		return fmt.Errorf("failed to publish usage tick event: %w", err)
	}

	return nil
}
