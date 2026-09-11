# Tasks

## High Priority

### Fix Redis connection log order
Move `log.Println("connected to Redis")` in `app.go` to after the `Ping` check so it only prints when the connection is confirmed.
- File: `internal/app/app.go`

---

## Auth

### JWT middleware (done)

### User model and event-driven sync
Create a `User` model and sync users from the portal via RabbitMQ events (`user.created`, `user.updated`).
- Create `internal/models/User.go` with UUID primary key
- Create a `user.created` / `user.updated` consumer in `internal/listeners/`
- Wire into bootstrap

### Token revocation / logout blocklist
When tokens are long-lived, a stolen token stays valid until expiry. Implement a Redis blocklist for logout.
- On logout, store the token JTI in Redis with TTL = token expiry
- In auth middleware, check if JTI is in the blocklist before accepting the token
- Alternative: switch to short-lived access tokens (15min) + long-lived refresh tokens in DB

---

## Rooms

### Implement missing RoomService methods
- [ ] `GetRoom(ctx, code string) (*models.Room, error)` — fetch and deserialize from Redis
- [ ] `JoinRoom(ctx, code, userID string) error` — get room, append userID to Members, save back
- [ ] `LeaveRoom(ctx, code, userID string) error` — get room, remove userID from Members, save back
- File: `internal/service/RoomService.go`

### Sync room state to Redis on events
On `play`, `pause`, `seek` WS events, update the room's `PlaybackState` in Redis so it's the source of truth.
- Hook into the hub's broadcast path or handle in a dedicated event processor
- Only sync meaningful state changes, not every message

### Late join state sync
When a new user joins a room mid-session, catch them up to the current playback state without resetting other members.
- [ ] On client register in hub, broadcast `state_request` to the room with the new joiner's userId
- [ ] Frontend host handles `state_request` and responds with current player position and playing state
- [ ] Hub handles `state_sync` event — route response to only the target userId, not broadcast to room
- [ ] After `state_sync` is delivered, new client enters normal broadcast flow
- Files: `internal/ws/hub.go`, `internal/ws/client.go`

### WebRTC signaling
Add signaling support for peer-to-peer connections within a room.
- Define signaling message types: `offer`, `answer`, `ice_candidate`
- Route signaling messages to specific peers (not broadcast to whole room) via the hub
- File: new `internal/ws/signaling.go` or extend `message.go`

---

## Models

### Switch to UUID primary keys
Replace auto-increment IDs with UUIDs across all models.
- [ ] Create a `Base` struct in `internal/models/base.go` with UUID `Id` and `BeforeCreate` GORM hook using `github.com/google/uuid`
- [ ] Update `Video` model to embed `Base` instead of `gorm.Model`
- [ ] Update `User` model to embed `Base`

---

## Infrastructure

### Remove MySQL driver dependency
`go.mod` still has `github.com/go-sql-driver/mysql` as an indirect dependency. Clean it up.
- Run `go mod tidy` after verifying nothing imports it

### gRPC for transcode events
Currently using RabbitMQ for transcode events between portal and media-sync. Future option to use gRPC instead.
- Define a `.proto` file for the transcode event contract
- Implement gRPC server in media-sync and client in portal

### Dead-letter queue for failed transcode jobs
If a transcode job fails repeatedly, it currently gets dropped.
- Add a dead-letter exchange in RabbitMQ topology
- Route failed messages to a `media.transcode.failed` queue for inspection and retry
- File: `media-sync-service/internal/broker/rabbitmq.go`

### streaming.transcode.results queue cleanup
The queue name and topology in media-sync need to match what streaming-service expects.
- Verify exchange, routing key, and queue name are consistent across both services
