package router

import (
	"stream-mesh/streaming/internal/controller"
	"stream-mesh/streaming/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	Video *controller.VideoController
	Room  *controller.RoomController
}

func Register(secret string, engine *gin.Engine, c *Controllers) {
	v1 := engine.Group("/api/v1")
	v1.Use(middleware.Auth(secret))
	registerVideoRoutes(v1, c.Video)
	registerRoomRoutes(v1, c.Room)
}
func registerVideoRoutes(g *gin.RouterGroup, v *controller.VideoController) {
	g.GET("/videos/:id", v.GetVideoById)
	g.GET("/videos/slug/:slug", v.GetVideoBySlug)

}

func registerRoomRoutes(g *gin.RouterGroup, r *controller.RoomController) {
	g.GET("/room/create", r.CreateRoom)
}
