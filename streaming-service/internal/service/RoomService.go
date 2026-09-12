package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"slices"
	"stream-mesh/streaming/internal/cache"
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
	//generate 16 char string for the room code
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 16)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func (r *RoomService) UpdatePlaybackState(ctx context.Context, event, roomCode string, data any) error {
	room, err := r.GetRoom(ctx, roomCode)
	log.Printf("[ROOM] GetRoom for %s: room=%v err=%v", roomCode, room, err)
	if err != nil {
		return err
	}
	if room == nil {
		return nil
	}
	//programmatic ping to refresh ttl of room in redis
	if event == "ping" {
		//reset the current state but with a fresh ttl
		if err := cache.Set[*models.Room](ctx, r.redis, "room:"+room.Code, room, time.Hour); err != nil {
			log.Printf("[ROOM] ping TTL refresh failed: %v", err)
			return err
		}
		log.Printf("[ROOM] ping TTL refreshed for %s", roomCode)
		return nil
	}
	//update room state in case of any event such as play pause
	bytes, _ := json.Marshal(data)
	var ps models.PlaybackState
	if err := json.Unmarshal(bytes, &ps); err != nil {
		return err
	}
	room.PlaybackState = &ps
	room.UpdatedAt = time.Now()

	if err := cache.Set[*models.Room](ctx, r.redis, "room:"+room.Code, room, time.Hour); err != nil {
		return err
	}
	return nil
}

func (r *RoomService) GetRoom(ctx context.Context, roomCode string) (*models.Room, error) {
	room, err := cache.Get[models.Room](ctx, r.redis, "room:"+roomCode)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (r *RoomService) JoinRoom(ctx context.Context, roomCode, userID string) (*models.Room, error) {
	room, err := r.GetRoom(ctx, roomCode)
	log.Printf("[ROOM] GetRoom for %s: room=%v err=%v", roomCode, room, err)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, fmt.Errorf("room not found")
	}
	if slices.Contains(room.Members, userID) {
		return room, nil
	}
	room.Members = append(room.Members, userID)
	room.UpdatedAt = time.Now()
	if err := cache.Set[*models.Room](ctx, r.redis, "room:"+room.Code, room, time.Hour); err != nil {
		return nil, err
	}
	return room, nil
}
func (r *RoomService) LeaveRoom(ctx context.Context, roomCode, userID string) error {
	room, err := r.GetRoom(ctx, roomCode)
	if err != nil {
		return err
	}
	if room == nil {
		return fmt.Errorf("room not found")
	}
	for i, member := range room.Members {
		if member == userID {
			room.Members = append(room.Members[:i], room.Members[i+1:]...)
			break
		}
	}
	room.UpdatedAt = time.Now()
	if err := cache.Set[*models.Room](ctx, r.redis, "room:"+room.Code, room, time.Hour); err != nil {
		return err
	}
	return nil
}
