package server

import (
	"context"
	"strings"
)

type OrderServer interface {
	Uppercase(ctx context.Context, s string) (string, error)
	Count(ctx context.Context, s string) int
}

type OrderService struct{}

func (order *OrderService) Uppercase(ctx context.Context, s string) (string, error) {
	Log.Log("Uppercase", s)
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (order *OrderService) Count(ctx context.Context, s string) int {
	return len(s)
}
