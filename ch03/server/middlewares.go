package server

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(OrderServer) OrderServer

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next OrderServer) OrderServer {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   OrderServer
	logger log.Logger
}

func (mw *loggingMiddleware) Uppercase(ctx context.Context, s string) (res string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Uppercase", "id", s, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Uppercase(ctx, s)
}

func (mw *loggingMiddleware) Count(ctx context.Context, id string) int {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetProfile", "id", id, "took", time.Since(begin))
	}(time.Now())
	return mw.next.Count(ctx, id)
}
