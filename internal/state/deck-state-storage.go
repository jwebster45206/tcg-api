package state

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jwebster45206/tcg-api/internal/deckstate"
)

type DeckStateStorage interface {
	SaveDeckState(ctx context.Context, gameID string, state *deckstate.DeckState) error
	GetDeckState(ctx context.Context, gameID string) (*deckstate.DeckState, error)
	DeleteDeckState(ctx context.Context, gameID string) error
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
