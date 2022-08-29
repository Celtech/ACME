package queue

import (
	"context"
	"fmt"
	"github.com/Celtech/ACME/config"
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
	EVENT_ISSUE  = "issue"
	EVENT_RENEW  = "renew"
	EVENT_REVOKE = "revoke"
)

type QueueEvent struct {
	RequestId     int       `json:"RequestId"`
	Domain        string    `json:"Domain"`
	ChallengeType string    `json:"ChallengeType"`
	Type          string    `json:"Type"`
	Attempt       int       `json:"Attempt"`
	CreatedAt     time.Time `json:"CreatedAt"`
}

func NewQueue(queue string) *QueueManager {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.GetConfig().GetString("redis.host"), config.GetConfig().GetString("redis.port")),
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

func (q *QueueManager) Close() error {
	return q.client.Close()
}
