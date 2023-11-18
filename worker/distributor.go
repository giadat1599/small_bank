package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, ops ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}


// Send tasks to redis queue
func NewRedisDistributor(redisOps asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOps)
	return &RedisTaskDistributor{
		client: client,		
	}
}