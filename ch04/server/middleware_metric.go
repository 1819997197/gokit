package server

import (
	"context"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"time"
)

var fieldKeys = []string{"method"}
var RequestCount = kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
	Namespace: "will",
	Subsystem: "arithmetic_service",
	Name:      "request_count",
	Help:      "Number of requests received.",
}, fieldKeys)

var RequestLatency = kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
	Namespace: "will",
	Subsystem: "arithemetic_service",
	Name:      "request_latency",
	Help:      "Total duration of requests in microseconds.",
}, fieldKeys)

type metricMiddleware struct {
	next           OrderServer
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func MetricsMiddleware(requestCount metrics.Counter, requestLatency metrics.Histogram) Middleware {
	return func(next OrderServer) OrderServer {
		return &metricMiddleware{
			next,
			requestCount,
			requestLatency,
		}
	}
}

func (m *metricMiddleware) Uppercase(ctx context.Context, s string) (res string, err error) {
	defer func(beign time.Time) {
		lvs := []string{"method", "Uppercase"}
		m.requestCount.With(lvs...).Add(1)
		m.requestLatency.With(lvs...).Observe(time.Since(beign).Seconds())
	}(time.Now())
	return m.next.Uppercase(ctx, s)
}

func (m *metricMiddleware) Count(ctx context.Context, id string) int {
	defer func(beign time.Time) {
		lvs := []string{"method", "Count"}
		m.requestCount.With(lvs...).Add(1)
		m.requestLatency.With(lvs...).Observe(time.Since(beign).Seconds())
	}(time.Now())
	return m.next.Count(ctx, id)
}
