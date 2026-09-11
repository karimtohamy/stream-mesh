package service

import (
	"context"
	"encoding/json"
	"math/rand"
	"stream-mesh/streaming/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type RoomService struct {
	redis *redis.Client
}

func NewRoomService(redis *redis.Client) *RoomService {
	return &RoomService{
		redis: redis,
	}
}
func (r *RoomService) CreateRoom(ctx context.Context, userID string) (*models.Room, error) {
	room := &models.Room{
		Code:      generateRoomCode(),
		Members:   []string{userID},
		UpdatedAt: time.Now(),
	}
	data, err := json.Marshal(room)
	if err != nil {
		return nil, err
	}
	if err := r.redis.Set(ctx, "room:"+room.Code, data, time.Hour).Err(); err != nil {
		return nil, err
	}
	return room, nil
}

func generateRoomCode() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 16)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
