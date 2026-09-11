# Streaming Service - TODO / Future Work

## Edge Cases & Business Logic

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

## Rooms (Redis-only, ephemeral)

- [ ] Room model: `id`, `members`, `currentVideo`, `playbackState`, `updatedAt`
- [ ] Room code: random 16-char string used as Redis key (`room:<code>`)
- [ ] WebSocket hub for real-time playback sync within a room
- [ ] WebRTC signaling for peer connections

## Models

- [ ] Switch primary keys to UUID (`github.com/google/uuid`) — custom `Base` struct with `BeforeCreate` hook
- [ ] Update `Video` model to embed `Base` instead of `gorm.Model`

## Infrastructure

- [ ] gRPC for transcode events (future — currently RabbitMQ)
- [ ] `streaming.transcode.results` queue topology cleanup in media-sync

## Broker

- [ ] Add dead-letter queue handling for failed transcode events
