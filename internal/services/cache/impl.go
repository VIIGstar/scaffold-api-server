package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	pkgerrors "github.com/pkg/errors"
)

// Ping checks for connection
func (c Client) Ping(ctx context.Context) error {
	_, err := c.redis.Ping(ctx).Result()
	return pkgerrors.WithStack(err)
}

// Set takes in interface
func (c Client) Set(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		if err == redis.Nil {
			return redis.Nil
		}
		return pkgerrors.WithStack(err)
	}

	return pkgerrors.WithStack(c.redis.Set(ctx, key, b, duration).Err())
}

// Get attempts to retrieve byte and unmarshal into result
func (c Client) Get(ctx context.Context, key string, value interface{}) error {
	b, err := c.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return redis.Nil
		}
		return pkgerrors.WithStack(err)
	}

	return pkgerrors.WithStack(json.Unmarshal(b, value))
}

// Delete attempts to delete entry by key
func (c Client) Delete(ctx context.Context, key string) error {
	return pkgerrors.WithStack(c.redis.Del(ctx, key).Err())
}

// Exists checks if a key exists in cache
func (c Client) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.redis.Exists(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, redis.Nil
		}
		return false, pkgerrors.WithStack(err)
	}

	return count >= 1, nil
}

// HMSet set hashmap interface with hash key
func (c Client) HMSet(ctx context.Context, hashKey string, data interface{}) error {
	mapInterface, ok := data.(map[string]interface{})
	if !ok {
		return errors.New("can not get map interface")
	}
	return c.redis.HMSet(ctx, hashKey, mapInterface).Err()
}

// HMGet get array interface with hash key
func (c Client) HMGet(ctx context.Context, hashKey string, arrValue interface{}, emptyJson string, keys ...string) error {
	res := c.redis.HMGet(ctx, hashKey, keys...)
	if res.Err() != nil {
		return res.Err()
	}
	jsonData := "["
	for _, v := range res.Val() {
		tmp, ok := v.(string)
		if !ok {
			tmp = emptyJson
		}
		jsonData += tmp
		jsonData += ","
	}
	jsonData = jsonData[:len(jsonData)-1]
	jsonData += "]"
	b := []byte(jsonData)
	return json.Unmarshal(b, &arrValue)
}
