package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	transporthttp "github.com/go-kit/kit/transport/http"
	"github.com/hashicorp/consul/api"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	//DirectClient()
	ServiceDiscoveryClient()
}

func ServiceDiscoveryClient() {
	res, err := ServiceDiscovery("GET", "http://localhost:8500", encodeRequest, decodeResponse, nil, "profilesvc", true, "arithmetic", "raysonxin")
	if err != nil {
		fmt.Println("err1: ", err)
		return
	}
	fmt.Println("res: ", res)
}

// ServiceDiscovery: 通过服务发现的形式调用服务
// registryAddress: 注册中心的地址
// servicesName: 注册的服务名称
// tags: 可用标签
// passingOnly: true 只返回通过健康监测的实例
// method:方法
// enc: http.EncodeRequestFunc dec: http.DecodeResponseFunc 这两个函数具体等一下会在Transport中进行详细解释
// requestStruct: 根据EndPoint定义的request结构体传参
func ServiceDiscovery(method string, registryAddress string, enc transporthttp.EncodeRequestFunc, dec transporthttp.DecodeResponseFunc, requestStruct interface{}, servicesName string, passingOnly bool, tags ...string) (interface{}, error) {
	consulCfg := api.DefaultConfig()
	consulCfg.Address = registryAddress
	consulClient, err := api.NewClient(consulCfg)
	if err != nil {
		return nil, err
	}

	// 获取服务ip + port 自己做负载均衡
	//result, _, err := consulClient.Catalog().Service("profilesvc", "", nil)
	//if err != nil {
	//	fmt.Println("ServiceDiscovery err: ", err)
	//	return nil, err
	//}
	//for _, v := range result {
	//	fmt.Println("id:", v.ServiceID, ", address:", v.ServiceAddress, ", port:", v.ServicePort)
	//}

	client := consul.NewClient(consulClient)
	logger := log.NewLogfmtLogger(os.Stdout)
	instances := consul.NewInstancer(client, logger, servicesName, tags, passingOnly)
	fmt.Println("instances: ", instances)
	f := func(servicesUrl string) (endpoint.Endpoint, io.Closer, error) {
		// 解析url
		target, err := url.Parse("http://" + servicesUrl)
		if err != nil {
			return nil, nil, err
		}
		target.Path = "/profiles/1234"
		return transporthttp.NewClient(strings.ToUpper(method), target, enc, dec).Endpoint(), nil, nil
	}
	endpointer := sd.NewEndpointer(instances, f, logger)
	var (
		maxAttempts = 3                      // per request, before giving up
		maxTime     = 250 * time.Millisecond // wallclock time, before giving up
	)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	return retry(context.Background(), requestStruct)
}

func DirectClient() {
	res, err := Direct("GET", "http://localhost:8080/profiles/1234", encodeRequest, decodeResponse, nil)
	if err != nil {
		fmt.Println("err1: ", err)
		return
	}
	fmt.Println("res: ", res)
}

// 直接调用服务
func Direct(method, fullUrl string, enc transporthttp.EncodeRequestFunc, dec transporthttp.DecodeResponseFunc, requestStruct interface{}) (interface{}, error) {
	target, err := url.Parse(fullUrl)
	if err != nil {
		fmt.Println("direct err:", err)
		return nil, err
	}
	client := transporthttp.NewClient(strings.ToUpper(method), target, enc, dec)
	return client.Endpoint()(context.Background(), requestStruct)
}

func encodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	if request == nil {
		return nil
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response getProfileResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

type getProfileResponse struct {
	Profile Profile `json:"profile,omitempty"`
	Err     error   `json:"err,omitempty"`
}

type Profile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Addresses []Address `json:"addresses,omitempty"`
}
type Address struct {
	ID       string `json:"id"`
	Location string `json:"location,omitempty"`
}

var (
	ErrNotFound = errors.New("not found")
)
