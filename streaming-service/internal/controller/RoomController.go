package controller

import (
	"net/http"
	"stream-mesh/streaming/internal/payload"
	"stream-mesh/streaming/internal/service"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	service *service.RoomService
}

func NewRoomController(roomService *service.RoomService) *RoomController {
	return &RoomController{
		service: roomService,
	}
}
func (r *RoomController) CreateRoom(c *gin.Context) {
	userId := c.GetString("user_id")
	room, err := r.service.CreateRoom(c, userId)
	if err != nil {
		payload.Fail(c, http.StatusBadRequest, "ROOM_CREATION_FAILED", "something went wrong")
		return
	}
	payload.OK(c, http.StatusCreated, room)
}
