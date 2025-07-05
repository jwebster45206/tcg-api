package state

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/deckstate"
)

type DeckStateLock interface {
	Release() error
	IsExpired() bool
	OwnerID() string
	Extend(duration time.Duration) error
}

type LockError struct {
	GameID  string
	Message string
}

func (e *LockError) Error() string {
	return e.Message
}

type DeckStateStorage interface {
	SaveDeckState(ctx context.Context, gameID string, state *deckstate.DeckState) error
	GetDeckState(ctx context.Context, gameID string) (*deckstate.DeckState, error)
	DeleteDeckState(ctx context.Context, gameID string) error

	LockDeckState(ctx context.Context, gameID string, timeout time.Duration) (DeckStateLock, error)
	IsLocked(ctx context.Context, gameID string) (bool, string, error) // locked, ownerID, error
}

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(redisClient *redis.Client) *RedisStorage {
	return &RedisStorage{
		client: redisClient,
	}
}

// Ensure that RedisStorage implements DeckStateStorage interface
var _ DeckStateStorage = (*RedisStorage)(nil)

const deckStateKeyPrefix = "deckstate:"

func (r *RedisStorage) SaveDeckState(ctx context.Context, gameID string, state *deckstate.DeckState) error {
	key := deckStateKeyPrefix + gameID
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, 4*time.Hour).Err()
}

func (r *RedisStorage) GetDeckState(ctx context.Context, gameID string) (*deckstate.DeckState, error) {
	key := deckStateKeyPrefix + gameID
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // No state found
		}
		return nil, err // Other error
	}
	var state deckstate.DeckState
	if err := json.Unmarshal([]byte(data), &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *RedisStorage) DeleteDeckState(ctx context.Context, gameID string) error {
	key := deckStateKeyPrefix + gameID
	return r.client.Del(ctx, key).Err()
}

type redisDeckStateLock struct {
	redisClient *redis.Client
	key         string
	ownerID     string
	expiration  time.Time
}

func (l *redisDeckStateLock) Release() error {
	_, err := l.redisClient.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
		pipe.Del(context.Background(), l.key)
		return nil
	})
	return err
}

func (l *redisDeckStateLock) IsExpired() bool {
	return time.Now().After(l.expiration)
}

func (l *redisDeckStateLock) OwnerID() string {
	return l.ownerID
}

func (l *redisDeckStateLock) Extend(duration time.Duration) error {
	l.expiration = time.Now().Add(duration)
	_, err := l.redisClient.Expire(context.Background(), l.key, duration).Result()
	return err
}

func (r *RedisStorage) LockDeckState(ctx context.Context, gameID string, timeout time.Duration) (DeckStateLock, error) {
	key := deckStateKeyPrefix + gameID + ":lock"
	ownerID := uuid.New().String()
	expiration := time.Now().Add(timeout)

	ok, err := r.client.SetNX(ctx, key, ownerID, timeout).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, &LockError{GameID: gameID, Message: "deck state is locked by another process"}
	}

	return &redisDeckStateLock{
		redisClient: r.client,
		key:         key,
		ownerID:     ownerID,
		expiration:  expiration,
	}, nil
}

func (r *RedisStorage) IsLocked(ctx context.Context, gameID string) (bool, string, error) {
	key := deckStateKeyPrefix + gameID + ":lock"
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, "", nil // Not locked
		}
		return false, "", err // Other error
	}
	return true, val, nil
}
