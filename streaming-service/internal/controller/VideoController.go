package controller

import (
	"fmt"
	"net/http"
	"stream-mesh/streaming/internal/payload"
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
	resp, err := v.service.GetVideoByID(c, id)
	if err != nil {
		msg := fmt.Sprintf("Video for id %s not found", id)
		payload.Fail(c, http.StatusNotFound, "NOT_FOUND", msg)
		fmt.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (v *VideoController) GetVideoBySlug(c *gin.Context) {
	slug := c.Param("slug")
	resp, err := v.service.GetBySlug(c, slug)
	if err != nil {
		msg := fmt.Sprintf("Video for slug %s not found", slug)
		payload.Fail(c, http.StatusNotFound, "NOT_FOUND", msg)
		fmt.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
