package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisLimiterRepository struct {
	RClient *redis.Client
}

func NewRedisRepository(rdb *redis.Client) IRateLimiterRepository {
	return &redisLimiterRepository{
		RClient: rdb,
	}
}

func (r *redisLimiterRepository) AllowToken(ctx context.Context, key string, tokenMaxRequest int, tokenBlockTime string) (bool, error) {
	pipe := r.RClient.TxPipeline()

	incr := pipe.Incr(ctx, key)

	duration, err := time.ParseDuration(tokenBlockTime)
	if err != nil {
		return false, err
	}

	pipe.Expire(ctx, key, time.Duration(duration))
	if _, err := pipe.Exec(ctx); err != nil {
		return false, err
	}

	return incr.Val() <= int64(tokenMaxRequest), nil
}

func (r *redisLimiterRepository) AllowIP(ctx context.Context, key string, ipMaxRequest int, ipBlockTime string) (bool, error) {
	pipe := r.RClient.TxPipeline()

	incr := pipe.Incr(ctx, key)

	duration, err := time.ParseDuration(ipBlockTime)
	if err != nil {
		return false, err
	}

	pipe.Expire(ctx, key, time.Duration(duration))
	if _, err := pipe.Exec(ctx); err != nil {
		return false, err
	}

	return incr.Val() <= int64(ipMaxRequest), nil
}
