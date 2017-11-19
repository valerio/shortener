package data

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

const counterKey = "url_counter_key"

type Storage interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Next() (int64, error)
}

type redisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr, password string, db int) *redisStorage {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	err := c.Ping().Err()
	if err != nil {
		log.Fatalf("Could not reach redis instance: %s", err.Error())
	}

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
