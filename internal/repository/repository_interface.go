package repository

import "context"

type IRateLimiterRepository interface {
	AllowToken(ctx context.Context, key string, tokenMaxRequest int, tokenBlockTime string) (bool, error)
	AllowIP(ctx context.Context, key string, ipMaxRequest int, ipBlockTime string)  (bool, error)
}
