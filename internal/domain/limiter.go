package domain

import (
	"context"
)

type RateLimiter struct {
	RateLimiterByIP        int
	RateLimiterTimeByIP    string
	RateLimiterByToken     int
	RateLimiterTimeByToken string
	IpFakeToTester         string
	Context                context.Context
}

func NewRateLimiter(rateLimiterByIP int,
	rateLimiterTimeByIP string,
	rateLimiterByToken int,
	rateLimiterTimeByToken string,
	ipFakeToTester string) *RateLimiter {
	return &RateLimiter{
		RateLimiterByIP:        rateLimiterByIP,
		RateLimiterTimeByIP:    rateLimiterTimeByIP,
		RateLimiterByToken:     rateLimiterByToken,
		RateLimiterTimeByToken: rateLimiterTimeByToken,
		IpFakeToTester:         ipFakeToTester,
		Context:                context.Background(),
	}
}
