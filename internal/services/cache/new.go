package cache

import (
	"context"
	"logur.dev/logur"

	"github.com/go-redis/redis/v8"
	pkgerrors "github.com/pkg/errors"
)

// Client is struct containing redis
type Client struct {
	logger logur.LoggerFacade
	redis  *redis.Client
}

// New creates new cache client
func New(ctx context.Context, address string, logger logur.LoggerFacade) (*Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr: address,
	})
	if _, err := r.Ping(ctx).Result(); err != nil {
		return &Client{}, pkgerrors.WithStack(err)
	}
	return &Client{redis: r, logger: logger}, nil
}

// TestCache initializes and return redis connection
// Should not use in production
func TestClient(ctx context.Context, logger logur.LoggerFacade) *Client {
	c, err := New(ctx, "127.0.0.1:6379", logger)
	if err != nil {
		panic(err) // ideally should not use panic in production, but for testing is ok
	}
	return c
}
