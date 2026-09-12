# Streaming Service - TODO / Future Work

## High Priority

- [ ] **Redis connected log message missing** — `app.go` prints the Redis connection log before the client is actually pinged and confirmed. Move `log.Println("connected to Redis")` to after the `Ping` check.

## Edge Cases & Business Logic

- [ ] **Room TTL heartbeat** — rooms expire after 1 hour of no state updates. Frontend sends a `ping` event every 55 minutes per WS connection. Hub handles `ping` by calling `roomService.UpdatePlaybackState` to refresh the TTL without changing state and without broadcasting to the room.

- [ ] **Video deleted from DB while cached in Redis**
  - If a video is removed from the `videos` table but still lives in Redis cache, `GetVideoBySlug` will return stale data until TTL expires.
  - Options to consider:
    1. On delete, explicitly evict the cache key (`video:slug:<slug>`)
    2. Use a short TTL so stale data expires quickly
    3. Publish a `video.deleted` event and have the service invalidate cache on receipt

## Auth

- [ ] JWT middleware for protected routes (streaming-service validates tokens issued by portal)
- [ ] Cache user lookups in Redis to avoid hitting DB on every request
- [ ] User model + event-driven sync from portal via RabbitMQ (`user.created`, `user.updated`)
- [ ] Token revocation / logout blocklist in Redis — needed if access tokens are long-lived, since a stolen token stays valid until expiry. Consider switching to short-lived access tokens (15min–1hr) + long-lived refresh tokens stored in DB that can be revoked.

## Rooms (Redis-only, ephemeral)

- [ ] Room model: `id`, `members`, `currentVideo`, `playbackState`, `updatedAt`
- [ ] Room code: random 16-char string used as Redis key (`room:<code>`)
- [x] WebSocket hub for real-time playback sync within a room
- [ ] WebRTC signaling for peer connections
- [ ] Implement missing RoomService methods: `GetRoom`, `JoinRoom`, `LeaveRoom` — sync member list to Redis on every join/leave
- [ ] Sync meaningful room state (current video, playback position) to Redis on `play`, `pause`, `seek` events — Redis is source of truth for room state, not every WS message
- [ ] **Late join state sync**
  - When a new user connects to a room, hub broadcasts a `state_request` event to the room with the new user's ID
  - An existing active client (e.g. the host) receives it and responds with `{ event: "state_sync", data: { position, playing } }`
  - Hub routes the `state_sync` response to only the joining user, not the whole room
  - After receiving `state_sync` the new user is a normal room member and receives all subsequent broadcasts
  - This ensures state comes from a live active player, not a potentially stale Redis snapshot
  - Subtasks:
    - [ ] On client register in hub, broadcast `state_request` to the room with the new joiner's userId
    - [ ] Add `state_request` event handling on the frontend — host responds with current player position and playing state
    - [ ] Add `state_sync` event handling in hub — route response to only the target userId, not the whole room
    - [ ] After `state_sync` is sent, new client enters normal broadcast flow

## Models

- [ ] Switch primary keys to UUID (`github.com/google/uuid`) — custom `Base` struct with `BeforeCreate` hook
- [ ] Update `Video` model to embed `Base` instead of `gorm.Model`

## Infrastructure

- [ ] gRPC for transcode events (future — currently RabbitMQ)
- [ ] `streaming.transcode.results` queue topology cleanup in media-sync

## Broker

- [ ] Add dead-letter queue handling for failed transcode events
