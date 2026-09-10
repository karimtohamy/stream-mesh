package controller

import (
	"fmt"
	"net/http"
	"stream-mesh/streaming/internal/service"

	"github.com/gin-gonic/gin"
)

type VideoController struct {
	service *service.VideoService
}

func NewVideoController(svc *service.VideoService) *VideoController {
	return &VideoController{
		service: svc,
	}
}

func (v *VideoController) GetVideoById(c *gin.Context) {
	id := c.Param("id")
	resp, err := v.service.GetVideoByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (v *VideoController) GetVideoBySlug(c *gin.Context) {
	slug := c.Param("slug")
	resp, err := v.service.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
