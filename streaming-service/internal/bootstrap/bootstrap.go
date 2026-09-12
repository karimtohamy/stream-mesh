package bootstrap

import (
	"context"
	"log"
	routes "stream-mesh/streaming/internal"
	"stream-mesh/streaming/internal/app"
	"stream-mesh/streaming/internal/controller"
	"stream-mesh/streaming/internal/listeners"
	"stream-mesh/streaming/internal/repository"
	"stream-mesh/streaming/internal/service"
	"stream-mesh/streaming/internal/ws"
)

func Init(ctx context.Context, app *app.App) {
	controllers := &routes.Controllers{Video: initVideo(ctx, app), Room: initRoom(ctx, app)}

	routes.Register(app.Cfg.App.Secret, app.Router, controllers)
}

func initVideo(ctx context.Context, app *app.App) *controller.VideoController {
	repo := repository.NewVideoRepository(app.Db)
	videoService := service.NewVideoService(repo, app.Redis)
	videoListener := listeners.NewListener(videoService)

	if err := app.Listener.Start(ctx, videoListener.OnTranscodeCompletes); err != nil {
		log.Print("video domain bootstrapping failed", err)
		panic(err)
	}
	return controller.NewVideoController(videoService)

}
func initRoom(ctx context.Context, app *app.App) *controller.RoomController {
	roomService := service.NewRoomService(app.Redis)
	hub := ws.NewHub(roomService)
	go hub.Run(ctx)
	return controller.NewRoomController(roomService, hub)
}
