package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     int
	Database int
}

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg RedisConfig) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: "",
		DB:       cfg.Database,
	})

	return &RedisClient{client: rdb}
}

func (c *RedisClient) RPush(ctx context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return c.client.RPush(ctx, key, data).Err()
}

func (c *RedisClient) LRange(ctx context.Context, key string, dest interface{}) error {
	if len(key) == 0 {
		return errors.New("key is empty")
	}

	vals, err := c.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	slicePtr := reflect.ValueOf(dest)
	if slicePtr.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer to a slice")
	}

	sliceValue := slicePtr.Elem()
	elemType := sliceValue.Type().Elem()

	for _, v := range vals {
		elemPtr := reflect.New(elemType)
		if err := json.Unmarshal([]byte(v), elemPtr.Interface()); err != nil {
			return err
		}
		sliceValue.Set(reflect.Append(sliceValue, elemPtr.Elem()))
	}

	return nil
}

func (c *RedisClient) GetJson(ctx context.Context, key string, src interface{}) error {
	if len(key) <= 0 {
		return errors.New("key is empty")
	}
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), &src)
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisClient) Close() error {
	return c.client.Close()
}
