package data

import (
	"time"

	"github.com/go-redis/redis"
)

const counterKey = "url_counter_key"

var Store Storage = newRedisStorage()

type Storage interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Next() (int64, error)
}

type redisStorage struct {
	client *redis.Client
}

func newRedisStorage() *redisStorage {
	c := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &redisStorage{
		client: c,
	}
}

func (rs redisStorage) Set(key string, value string) error {
	_, err := rs.client.Set(key, value, time.Hour*24).Result()
	return err
}

func (rs redisStorage) Get(key string) (string, error) {
	value, err := rs.client.Get(key).Result()
	return value, err
}

func (rs redisStorage) Next() (int64, error) {
	value, err := rs.client.Incr(counterKey).Result()
	return value, err
}
