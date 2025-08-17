package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	ports "pinstack-user-service/internal/domain/ports/output"
	"pinstack-user-service/internal/infrastructure/config"

	"github.com/redis/go-redis/v9"
	"github.com/soloda1/pinstack-proto-definitions/custom_errors"
)

type Client struct {
	client *redis.Client
	log    ports.Logger
}

func NewClient(cfg config.Redis, log ports.Logger) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Error("Failed to connect to Redis", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info("Successfully connected to Redis",
		slog.String("address", cfg.Address),
		slog.Int("port", cfg.Port),
		slog.Int("db", cfg.DB))

	return &Client{
		client: rdb,
		log:    log,
	}, nil
}

func (c *Client) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.log.Debug("Cache miss", slog.String("key", key))
			return custom_errors.ErrCacheMiss
		}
		c.log.Error("Failed to get from cache",
			slog.String("key", key),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		c.log.Error("Failed to unmarshal cache value",
			slog.String("key", key),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to unmarshal cache value: %w", err)
	}

	c.log.Debug("Cache hit", slog.String("key", key))
	return nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		c.log.Error("Failed to marshal value for cache",
			slog.String("key", key),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		c.log.Error("Failed to set cache",
			slog.String("key", key),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to set cache: %w", err)
	}

	c.log.Debug("Successfully set cache",
		slog.String("key", key),
		slog.Duration("ttl", ttl))
	return nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	result, err := c.client.Del(ctx, key).Result()
	if err != nil {
		c.log.Error("Failed to delete from cache",
			slog.String("key", key),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to delete from cache: %w", err)
	}

	if result == 0 {
		c.log.Debug("Key not found for deletion", slog.String("key", key))
	} else {
		c.log.Debug("Successfully deleted from cache", slog.String("key", key))
	}

	return nil
}

func (c *Client) Close() error {
	if err := c.client.Close(); err != nil {
		c.log.Error("Failed to close Redis connection", slog.String("error", err.Error()))
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	c.log.Info("Redis connection closed")
	return nil
}

func (c *Client) Ping(ctx context.Context) error {
	if err := c.client.Ping(ctx).Err(); err != nil {
		c.log.Error("Redis ping failed", slog.String("error", err.Error()))
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}
