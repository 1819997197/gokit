package main

import (
	"flag"
	"fmt"
	"gokit/profilesvc/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	profilesvc "gokit/profilesvc/server"
)

func main() {
	var (
		consulHost  = flag.String("consul.host", "127.0.0.1", "consul ip address")
		consulPort  = flag.String("consul.port", "8500", "consul port")
		serviceHost = flag.String("service.host", "127.0.0.1", "service ip address")
		servicePort = flag.String("service.port", "8080", "service port")
	)
	flag.Parse()

	logger.InitLog()

	var s profilesvc.Service
	{
		s = profilesvc.NewInmemService()
		s = profilesvc.LoggingMiddleware(logger.Log)(s)
	}

	var h http.Handler
	{
		h = profilesvc.MakeHTTPHandler(s, log.With(logger.Log, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	//创建注册对象
	registar := profilesvc.Register(*consulHost, *consulPort, *serviceHost, *servicePort, logger.Log)
	go func() {
		registar.Register()
		logger.Log.Log("Http Server start at port:" + *servicePort)
		errs <- http.ListenAndServe(":"+*servicePort, h)
	}()

	logger.Log.Log("exit", <-errs)
	//服务退出取消注册
	registar.Deregister()
}
