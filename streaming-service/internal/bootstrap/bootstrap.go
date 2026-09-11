package bootstrap

import (
	"context"
	"log"
	routes "stream-mesh/streaming/internal"
	"stream-mesh/streaming/internal/app"
	"stream-mesh/streaming/internal/controller"
	"stream-mesh/streaming/internal/listeners"
	"stream-mesh/streaming/internal/repository"
	videoService "stream-mesh/streaming/internal/service"
)

func Init(ctx context.Context, app *app.App) {
	controllers := &routes.Controllers{Video: initVideo(ctx, app)}
	routes.Register(app.Router, controllers)
}

func initVideo(ctx context.Context, app *app.App) *controller.VideoController {
	repo := repository.NewVideoRepository(app.Db)
	service := videoService.NewVideoService(repo, app.Redis)
	videoListener := listeners.NewListener(service)

	if err := app.Listener.Start(ctx, videoListener.OnTranscodeCompletes); err != nil {
		log.Print("video domain bootstrapping failed", err)
		panic(err)
	}
	return controller.NewVideoController(service)

}
