package queue

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisQueue struct {
	client    *redis.Client
	queueName string
	dlqName   string
}

func NewRedisQueue(client *redis.Client, queueName, dlqName string) *RedisQueue {
	return &RedisQueue{
		client:    client,
		queueName: queueName,
		dlqName:   dlqName,
	}
}

func (q *RedisQueue) Enqueue(ctx context.Context, data []byte) error {
	return q.client.LPush(ctx, q.queueName, data).Err()
}

// DequeueBlocking blocks up to timeout and returns the message.
// Returns ("", nil) on timeout with no data.
func (q *RedisQueue) DequeueBlocking(ctx context.Context, timeout time.Duration) ([]byte, error) {
	res, err := q.client.BRPop(ctx, timeout, q.queueName).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	if len(res) != 2 {
		return nil, nil
	}
	return []byte(res[1]), nil
}

func (q *RedisQueue) EnqueueDLQ(ctx context.Context, data []byte) error {
	return q.client.LPush(ctx, q.dlqName, data).Err()
}

func (q *RedisQueue) GetQueueLength(ctx context.Context) (int64, error) {
	return q.client.LLen(ctx, q.queueName).Result()
}
