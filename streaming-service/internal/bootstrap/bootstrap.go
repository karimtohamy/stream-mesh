package bootstrap

import (
	"context"
	"log"
	routes "stream-mesh/streaming/internal"
	"stream-mesh/streaming/internal/broker"
	videoController "stream-mesh/streaming/internal/controller"
	"stream-mesh/streaming/internal/listeners"
	"stream-mesh/streaming/internal/repository"
	videoService "stream-mesh/streaming/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(ctx context.Context, db *gorm.DB, engine *gin.Engine, listener *broker.Listener) {
	initVideo(ctx, db, engine, listener)
}

func initVideo(ctx context.Context, db *gorm.DB, engine *gin.Engine, listener *broker.Listener) {
	repo := repository.NewVideoRepository(db)
	service := videoService.NewVideoService(repo)
	videoListener := listeners.NewListener(service)

	controller := videoController.NewVideoController(service)
	routes.Register(
		engine,
		&routes.Controllers{
			Video: controller,
		},
	)
	if err := listener.Start(ctx, videoListener.OnTranscodeCompletes); err != nil {
		log.Print("video domain bootstrapping failed", err)
		panic(err)
	}

}
