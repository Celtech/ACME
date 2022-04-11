package queue

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type QueueManager struct {
	client *redis.Client
	ctx    context.Context
	queue  string
}

func NewQueue(queue string) *QueueManager {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.TODO()

	return &QueueManager{
		client: redisClient,
		ctx:    ctx,
		queue:  queue,
	}
}
