package listeners

import (
	"context"
	"stream-mesh/streaming/internal/broker"
	"stream-mesh/streaming/internal/models"
	"stream-mesh/streaming/internal/service"
)

type Listener struct {
	service *service.VideoService
}

func NewListener(service *service.VideoService) *Listener {
	return &Listener{
		service: service,
	}
}

func (l *Listener) OnTranscodeCompletes(ctx context.Context, event broker.NewVideoEvent) error {

	video := &models.Video{
		Slug:        event.Slug,
		ManifestURL: event.ManifestURL,
		Thumbnail:   event.ThumbNailURL,
	}
	return l.service.Save(ctx, video)
}
