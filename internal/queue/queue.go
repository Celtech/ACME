package queue

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var QueueMgr *QueueManager

type QueueManager struct {
	client *redis.Client
	ctx    context.Context
	queue  string
}

const (
	EVENT_REQUEST = "request"
	EVENT_RENEW   = "renew"
	EVENT_REVOKE  = "revoke"
)

type QueueEvent struct {
	Domain        string
	ChallengeType string
	Type          string
	Attempt       int
	CreatedAt     time.Time
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
