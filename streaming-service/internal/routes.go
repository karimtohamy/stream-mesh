package router

import (
	"stream-mesh/streaming/internal/controller"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	Video *controller.VideoController
}

func Register(engine *gin.Engine, c *Controllers) {
	v1 := engine.Group("/api/v1")
	registerVideoRoutes(v1, c.Video)
}
func registerVideoRoutes(g *gin.RouterGroup, v *controller.VideoController) {
	g.GET("/videos/:id", v.GetVideoById)
	g.GET("/videos/slug/:slug", v.GetVideoBySlug)

}
