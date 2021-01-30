package server

import (
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

func MakeHTTPHandler(svc OrderServer) http.Handler {
	endpoint := makeUppercaseEndpoint(svc)
	//ratebucket := ratelimit.NewBucket(time.Second*1, 3)
	//endpoint = NewTokenBucketLimitterWithJuju(ratebucket)(endpoint)
	ratebucket := rate.NewLimiter(rate.Every(time.Second*1), 3)
	endpoint = NewTokenBucketLimitterWithBuildIn(ratebucket)(endpoint) //限流中间件
	uppercaseHandler := httptransport.NewServer(
		endpoint,
		decodeUppercaseRequest,
		encodeResponse,
	)
	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Methods("POST").Path("/uppercase").Handler(uppercaseHandler)
	r.Methods("GET").Path("/count/{name}").Handler(countHandler)
	r.Path("/metrics").Handler(promhttp.Handler()) //用于Prometheus轮循拉取监控指标的接口
	return r
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, ErrBadRouting
	}
	var request countRequest
	request.S = name
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
