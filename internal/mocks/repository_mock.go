package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRateLimiterRepo struct {
	mock.Mock
}

func (m *MockRateLimiterRepo) AllowToken(ctx context.Context, token string, maxReq int, blockTime string) (bool, error) {
	args := m.Called(ctx, token, maxReq, blockTime)
	return args.Bool(0), args.Error(1)
}

func (m *MockRateLimiterRepo) AllowIP(ctx context.Context, ip string, maxReq int, blockTime string) (bool, error) {
	args := m.Called(ctx, ip, maxReq, blockTime)
	return args.Bool(0), args.Error(1)
}
