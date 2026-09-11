package controller

import (
	"net/http"
	"stream-mesh/streaming/internal/payload"
	"stream-mesh/streaming/internal/service"
	"stream-mesh/streaming/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type RoomController struct {
	service *service.RoomService
	hub     *ws.Hub
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewRoomController(roomService *service.RoomService, hub *ws.Hub) *RoomController {
	return &RoomController{
		service: roomService,
		hub:     hub,
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
func (r *RoomController) Connect(c *gin.Context) {
	roomCode := c.Param("code")
	userId := c.GetString("user_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		payload.Fail(c, http.StatusInternalServerError, "WS_UPGRADE_FAILED", "failed to upgrade connection")
		return
	}

	client := ws.NewClient(r.hub, conn, roomCode, userId)
	r.hub.Register(client)

	go client.ReadPump()
	go client.WritePump()
}
